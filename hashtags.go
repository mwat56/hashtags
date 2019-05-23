/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
                  All rights reserved
              EMail : <support@mwat.de>
*/

package hashtags

import (
	"bufio"
	"hash/crc32"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync/atomic"
)

type (
	// `tSourceList` is storing the IDs of #hashtags/@mentions.
	tSourceList []string

	// `tHashMap` is indexed by #hashtags/@mentions holding a `tSourceList`.
	tHashMap map[string]*tSourceList

	// TCountItem holds a #hashtag/@mention and its number of occurances.
	//
	// @see CountedList()
	TCountItem = struct {
		Count int    // number of IDs for this #hashtag/@mention
		Tag   string // name of #hashtag/@mention
	}

	// A list of `TCountItem`s
	tCountList []TCountItem

	// Data cache for `CountedList()`
	tCountCache struct {
		µCRC    uint32
		µCounts tCountList
	}

	// THashList is a list of `#hashtags` and `@mentions`
	// pointing to sources (i.e. IDs).
	THashList struct {
		fn      string // the filename to use
		hl      tHashMap
		µChange uint32
		µCC     tCountCache
	}
)

// `add()` appends 'aID` to the list
//
// `aID` the source ID to add to the list.
func (sl *tSourceList) add(aID string) *tSourceList {
	for _, id := range *sl {
		if id == aID {
			// already in list
			return sl
		}
	}
	*sl = append(*sl, aID)

	return sl
} // add()

// `clear()` removes all entries in this list.
func (sl *tSourceList) clear() *tSourceList {
	(*sl) = (*sl)[:0]

	return sl
} // clear()

// `indexOf()` returns the list index of `aID`.
//
// `aID` is the string to look up.
func (sl *tSourceList) indexOf(aID string) int {
	for result, id := range *sl {
		if id == aID {
			return result
		}
	}

	return -1
} // indexOf()

// `removeID()` deletes the list entry of `aID`.
//
// `aID` is the string to look up.
func (sl *tSourceList) removeID(aID string) *tSourceList {
	idx := sl.indexOf(aID)
	if 0 > idx {
		return sl
	}

	slen := len(*sl) - 1
	if 0 > slen {
		// can't remove from empty list …
		return sl
	}

	if 0 == idx {
		if 0 == slen {
			*sl = (*sl)[:0]
		} else {
			*sl = (*sl)[1:]
		}
	} else if slen == idx {
		*sl = (*sl)[:slen]
	} else {
		*sl = append((*sl)[:idx], (*sl)[idx+1:]...)
	}

	return sl
} // removeID()

// `renameID()` replaces all occurances of `aOldID` by `aNewID`.
//
// This method is intended for rare cases when the ID of a document
// gets changed.
//
// `aOldID` is to be replaced in this list.
//
// `aNewID` is the replacement in this list.
func (sl *tSourceList) renameID(aOldID, aNewID string) *tSourceList {
	for idx, id := range *sl {
		if id == aOldID {
			(*sl)[idx] = aNewID
			return sl.sort()
		}
	}

	return sl
} // renameID()

// `sort()` returns the sorted list.
func (sl *tSourceList) sort() *tSourceList {
	sort.Slice(*sl, func(i, j int) bool {
		return ((*sl)[i] < (*sl)[j]) // ascending
	})

	return sl
} // sort()

// String returns the list as a linefeed seperated string.
//
// (Implements `Stringer` interface)
func (sl *tSourceList) String() string {
	sl.sort()

	return strings.Join(*sl, "\n")
} // String()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

// `add()` appends `aID` to the list associated with `aMapIdx`.
//
// If either `aMapIdx` or `aID` are empty strings they are silently
// ignored (i.e. this method does nothing).
//
// `aDelim` is the start character of words to use (i.e. either '@' or '#').
//
// `aMapIdx` is the list index to lookup.
//
// `aID` is to be added to the hash list.
func (hl *THashList) add(aDelim byte, aMapIdx, aID string) *THashList {
	if (0 == len(aMapIdx)) || (0 == len(aID)) {
		return hl
	}
	aMapIdx = strings.ToLower(aMapIdx) // prepare for case-insensitive search
	if aMapIdx[0] != aDelim {
		aMapIdx = string(aDelim) + aMapIdx
	}

	return hl.add0(aMapIdx, aID)
} // add()

// `add0()` appends `aID` to the list associated with `aMapIdx`.
//
// `aMapIdx` is the list index to lookup.
//
// `aID` is to be added to the hash list.
func (hl *THashList) add0(aMapIdx, aID string) *THashList {
	if sl, ok := hl.hl[aMapIdx]; ok {
		sl.add(aID).sort()
	} else {
		sl := make(tSourceList, 1, 32)
		sl[0] = aID
		hl.hl[aMapIdx] = &sl
	}
	atomic.StoreUint32(&hl.µChange, 0)

	return hl
} // add0()

// Checksum returns the list's CRC32 checksum.
//
// This method can be used to get a kind of 'footprint'.
func (hl *THashList) Checksum() uint32 {
	if 0 != atomic.LoadUint32(&hl.µChange) {
		return hl.µChange
	}
	// We use `String()` because it sorts internally thus
	// generating reproducible results:
	atomic.StoreUint32(&hl.µChange, crc32.Update(0, crc32.MakeTable(crc32.Castagnoli), []byte(hl.String())))

	return hl.µChange
} // Checksum()

// Clear empties the internal data structures:
// all `#hashtags` and `@mentions` are deleted.
func (hl *THashList) Clear() *THashList {
	for mapIdx, sl := range hl.hl {
		sl.clear()
		delete(hl.hl, mapIdx)
	}
	atomic.StoreUint32(&hl.µChange, 0)

	return hl
} // Clear()

// CountedList returns a list of #hashtags/@mentions with
// their respective count of associated IDs.
func (hl *THashList) CountedList() []TCountItem {
	if (hl.µCC.µCRC == hl.µChange) && (0 < len(hl.µCC.µCounts)) {
		return hl.µCC.µCounts
	}

	hl.µCC.µCounts = nil
	hl.µCC.µCRC = hl.µChange
	result := make(tCountList, 0, len(hl.hl))
	for mapIdx, sl := range hl.hl {
		result = append(result, TCountItem{len(*sl), mapIdx})
	}
	if 0 < len(result) {
		sort.Slice(result, func(i, j int) bool {
			// ignore [#@] for sorting
			return (result[i].Tag[1:] < result[j].Tag[1:])
		})
	}
	hl.µCC.µCounts = result

	return result
} // CountedList()

// Filename returns the configured filename for reading/storing this list.
func (hl *THashList) Filename() string {
	return hl.fn
} // Filename()

// HashAdd appends `aID` to the list of `aHash`.
//
// If either `aHash` or `aID` are empty strings they are silently
// ignored (i.e. this method does nothing).
//
// `aHash` is the list index to lookup.
//
// `aID` is to be added to the hash list.
func (hl *THashList) HashAdd(aHash, aID string) *THashList {
	return hl.add('#', aHash, aID)
} // HashAdd()

// HashLen returns the number of IDs stored for `aHash`.
//
// `aHash` identifies the ID list to lookup.
func (hl *THashList) HashLen(aHash string) int {
	return hl.idxLen('#', aHash)
} // HashLen()

// HashList returns a list of IDs associated with `aHash`.
//
// `aHash` identifies the ID list to lookup.
func (hl *THashList) HashList(aHash string) []string {
	return hl.list('#', aHash)
} // HashList()

// HashRemove deletes `aID` from the list of `aHash`.
//
// `aHash` identifies the ID list to lookup.
//
// `aID` is the source to remove from the list.
func (hl *THashList) HashRemove(aHash, aID string) *THashList {
	return hl.remove('#', aHash, aID)
} // HashRemove()

// IDlist returns a list of #hashtags and @mentions associated with `aID`.
func (hl *THashList) IDlist(aID string) (rList []string) {
	for mapIdx, sl := range hl.hl {
		if 0 <= sl.indexOf(aID) {
			rList = append(rList, mapIdx)
		}
	}
	if 0 < len(rList) {
		sort.Slice(rList, func(i, j int) bool {
			return (rList[i] < rList[j]) // ascending
		})
	}

	return
} // IDlist()

var (
	// match: #hashtag|@mention
	hashMentionRE = regexp.MustCompile(`(?i)\b?([@#][\wÄÖÜß-]+)`)
)

// IDparse checks whether `aText` contains strings starting with `[@|#]`
// and – if found – adds them to the respective list.
//
// `aID` is the ID to add to the list.
//
// `aText` is the text to search.
func (hl *THashList) IDparse(aID string, aText []byte) *THashList {
	oldCRC := hl.Checksum()
	defer func() {
		if oldCRC != atomic.LoadUint32(&hl.µChange) {
			hl.Store()
		}
	}()

	matches := hashMentionRE.FindAllSubmatch(aText, -1)
	if (nil == matches) || (0 >= len(matches)) {
		return hl
	}
	for _, sub := range matches {
		if 0 < len(sub[1]) {
			hl.add0(strings.ToLower(string(sub[1])), aID)
		}
	}

	return hl
} // IDparse()

// IDremove deletes all @hashtags/@mentions associated with `aID`.
//
// `aID` is to be deleted from all lists.
func (hl *THashList) IDremove(aID string) *THashList {
	oldCRC := hl.Checksum()
	defer func() {
		if oldCRC != atomic.LoadUint32(&hl.µChange) {
			hl.Store()
		}
	}()

	for mapIdx, sl := range hl.hl {
		sl.removeID(aID)
		if 0 == len(*hl.hl[mapIdx]) {
			delete(hl.hl, mapIdx)
		}
	}
	atomic.StoreUint32(&hl.µChange, 0)

	return hl
} // IDremove()

// IDrename replaces all occurances of `aOldID` by `aNewID`.
//
// This method is intended for rare cases when the ID of a document
// needs to get changed.
//
// `aOldID` is to be replaced in all lists.
//
// `aNewID` is the replacement in all lists.
func (hl *THashList) IDrename(aOldID, aNewID string) *THashList {
	for _, sl := range hl.hl {
		sl.renameID(aOldID, aNewID)
	}
	atomic.StoreUint32(&hl.µChange, 0)
	hl.Store()

	return hl
} // IDrename()

// IDupdate checks `aText` removing all #hashtags/@mentions no longer
// present and adds #hashtags/@mentions new in `aText`.
//
// `aID` is the ID to update.
//
// `aText` is the text to use.
func (hl *THashList) IDupdate(aID string, aText []byte) *THashList {
	hl.IDremove(aID)

	oldCRC := hl.Checksum()
	defer func() {
		if oldCRC != atomic.LoadUint32(&hl.µChange) {
			hl.Store()
		}
	}()

	hl.IDremove(aID)
	matches := hashMentionRE.FindAllSubmatch(aText, -1)
	if (nil == matches) || (0 >= len(matches)) {
		return hl
	}
	for _, sub := range matches {
		if 0 < len(sub[1]) {
			hl = hl.add(sub[1][0], string(sub[1]), aID)
		}
	}

	return hl
} // IDupdate()

// `idxLen()` returns the number of IDs stored for `aMapIdx`.
//
// `aDelim` is the first character of words to use (i.e. either '@' or '#').
//
// `aMapIdx` identifies the ID list to lookup.
func (hl *THashList) idxLen(aDelim byte, aMapIdx string) int {
	if 0 == len(aMapIdx) {
		return -1
	}
	aMapIdx = strings.ToLower(aMapIdx)
	if aMapIdx[0] != aDelim {
		aMapIdx = string(aDelim) + aMapIdx
	}
	if sl, ok := (hl.hl)[aMapIdx]; ok {
		return len(*sl)
	}

	return -1
} // idxLen()

// Len returns the current length of the list i.e. how many #hashtags
// and @mentions are currently stored in the list.
func (hl *THashList) Len() int {
	return len(hl.hl)
} // Len()

// LenTotal returns the length of all #hashtag/@mention lists together.
func (hl *THashList) LenTotal() int {
	result := len(hl.hl)
	for _, sl := range hl.hl {
		result += len(*sl)
	}

	return result
} // LenTotal()

// `list()` returns a list of IDs associated with `aMapIdx`.
//
// `aDelim` is the start of words to search (i.e. either '@' or '#').
//
// `aMapIdx` identifies the sources list to lookup.
func (hl *THashList) list(aDelim byte, aMapIdx string) (rList []string) {
	if 0 == len(aMapIdx) {
		return
	}
	aMapIdx = strings.ToLower(aMapIdx)
	if aMapIdx[0] != aDelim {
		aMapIdx = string(aDelim) + aMapIdx
	}
	if sl, ok := hl.hl[aMapIdx]; ok {
		sl.sort()
		rList = []string(*sl)
	}

	return
} // list()

// Load reads the configured filen returning the data structure
// read from the file and a possible error condition.
//
// If the hash file doesn't exist that is not considered an error.
// If there is an error, it will be of type `*PathError`.
func (hl *THashList) Load() (*THashList, error) {
	file, err := os.OpenFile(hl.fn, os.O_RDONLY, 0)
	if nil != err {
		if os.IsNotExist(err) {
			return hl, nil
		}
		return hl, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	_, err = hl.Clear().read(scanner)

	return hl, err
} // Load()

// MentionAdd appends `aID` to the list of `aMention`.
//
// If either `aMention` or `aID` are empty strings they are
// silently ignored (i.e. this method does nothing).
//
// `aMention` is the list index to lookup.
//
// `aID` is to be added to the hash list.
func (hl *THashList) MentionAdd(aMention, aID string) *THashList {
	return hl.add('@', aMention, aID)
} // MentionAdd()

// MentionLen returns the number of IDs stored for `aMention`.
//
// `aMention` identifies the ID list to lookup.
func (hl *THashList) MentionLen(aMention string) int {
	return hl.idxLen('@', aMention)
} // MentionLen()

// MentionList returns a list of IDs associated with `aHash`.
//
// `aMention` identifies the ID list to lookup.
func (hl *THashList) MentionList(aMention string) []string {
	return hl.list('@', aMention)
} // MentionList()

// MentionRemove deletes `aID` from the list of `aMention`.
//
// `aMention` identifies the ID list to lookup.
//
// `aID` is the source to remove from the list.
func (hl *THashList) MentionRemove(aMention, aID string) *THashList {
	return hl.remove('@', aMention, aID)
} // MentionRemove()

var (
	// match: [#Hashtag|@mention]
	hashHeadRE = regexp.MustCompile(`^\[\s*([#@][^\]]*?)\s*\]$`)
)

// `read()` parses a file written by `Store()` returning the
// number of bytes read and a possible error.
//
// This method reads one line of the file at a time.
func (hl *THashList) read(aScanner *bufio.Scanner) (rRead int, rErr error) {
	var mapIdx string

	for lineRead := aScanner.Scan(); lineRead; lineRead = aScanner.Scan() {
		line := aScanner.Text()
		rRead += len(line) + 1 // add trailing LF

		line = strings.TrimSpace(line)
		if 0 == len(line) {
			continue
		}

		if matches := hashHeadRE.FindStringSubmatch(line); nil != matches {
			mapIdx = strings.TrimSpace(matches[1])
		} else {
			hl.add(mapIdx[0], mapIdx, line)
		}
	}
	atomic.StoreUint32(&hl.µChange, 0)
	rErr = aScanner.Err()

	return
} // read()

// `remove()` deletes `aID` from the list of `aMapIdx`.
//
// `aDelim` is the start character of words to use (i.e. either '@' or '#').
//
// `aMapIdx` identifies the sources list to lookup.
//
// `aID` is the source to remove from the list.
func (hl *THashList) remove(aDelim byte, aMapIdx, aID string) *THashList {
	if (0 == len(aMapIdx)) || (0 == len(aID)) {
		return hl
	}
	aMapIdx = strings.ToLower(aMapIdx)
	if aMapIdx[0] != aDelim {
		aMapIdx = string(aDelim) + aMapIdx
	}
	if sl, ok := hl.hl[aMapIdx]; ok {
		sl.removeID(aID)
		if 0 == len(*hl.hl[aMapIdx]) {
			delete(hl.hl, aMapIdx)
		}
		atomic.StoreUint32(&hl.µChange, 0)
	}

	return hl
} // remove()

// SetFilename sets `aFilename` to use by this list.
func (hl *THashList) SetFilename(aFilename string) *THashList {
	hl.fn = aFilename

	return hl
} // SetFilename()

// Store writes the whole list to the configured filen
// returning the number of bytes written and a possible error.
//
// If there is an error, it will be of type `*PathError`.
func (hl *THashList) Store() (int, error) {
	s := hl.String()
	file, err := os.OpenFile(hl.fn, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0664)
	if nil != err {
		return 0, err
	}
	defer file.Close()

	return file.Write([]byte(s))
} // Store()

// String returns the whole list as a linefeed separated string.
func (hl *THashList) String() string {
	var (
		result string
		tmp    tSourceList
	)
	for hash := range hl.hl {
		tmp = append(tmp, hash)
	}
	// sort the order of hashtags to get a reproducible result
	sort.Slice(tmp, func(i, j int) bool {
		// ignore leading [@#] when sorting
		return (tmp[i][1:] < tmp[j][1:]) // ascending
	})
	for _, hash := range tmp {
		sl, _ := hl.hl[hash]
		result += "[" + hash + "]\n" + sl.String() + "\n"
	}

	return result
} // String()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

type (
	// TWalkFunc is used by `Walk()` when visiting an entry
	// in the #hashtag/@mention lists.
	//
	// see `Walk()`
	TWalkFunc func(aHash, aID string) (rValid bool)

	// THashWalker is used by `Walker()` when visiting an entry
	// in the #hashtag/@mentions lists.
	//
	// see `Walker()`
	THashWalker interface {
		Walk(aHash, aID string) bool
	}
)

// Walk traverses through all entries in the #hashtag/@mention lists
// calling `aFunc` for each entry.
//
// `aFunc` is the function called for each ID in all lists.
func (hl *THashList) Walk(aFunc TWalkFunc) {
	oldCRC := hl.Checksum()
	defer func() {
		if oldCRC != atomic.LoadUint32(&hl.µChange) {
			hl.Store()
		}
	}()

	for hash, sl := range hl.hl {
		for _, id := range *sl {
			if !aFunc(hash, id) {
				sl.removeID(id)
				atomic.StoreUint32(&hl.µChange, 0)
			}
		}
	}
	for hash, sl := range hl.hl {
		if 0 == len(*sl) {
			delete(hl.hl, hash)
			atomic.StoreUint32(&hl.µChange, 0)
		}
	}
} // Walk()

// Walker traverses through all entries in the INI list sections
// calling `aWalker` for each entry.
//
// `aWalker` is an object implementing the `TIniWalker` interface.
func (hl *THashList) Walker(aWalker THashWalker) {
	hl.Walk(aWalker.Walk)
} // Walker()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

// New returns a new `THashList` instance after reading
// the given file.
//
// If the hash file doesn't exist that is not considered an error.
// If there is an error, it will be of type *PathError.
//
// `aFilename` is the name of the file to use for reading and storing.
func New(aFilename string) (*THashList, error) {
	result := THashList{
		fn: aFilename,
		hl: make(tHashMap, 64),
	}

	return result.Load()
} // New()

/* EoF */
