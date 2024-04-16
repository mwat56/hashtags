/*
Copyright © 2019, 2024  M.Watermann, 10247 Berlin, Germany

	    All rights reserved
	EMail : <support@mwat.de>
*/
package hashtags

//lint:file-ignore ST1017 - I prefer Yoda conditions

import (
	"bufio"
	"encoding/gob"
	"errors"
	"hash/crc32"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
)

type (
	// `tSourceList` is storing the IDs of a single #hashtag/@mention.
	tSourceList []string

	// `tHashMap` is indexed by #hashtags/@mentions holding a `tSourceList`.
	tHashMap map[string]*tSourceList

	// TCountItem holds a #hashtag/@mention and its number of occurrences.
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
		fn      string        // the filename to use
		hl      tHashMap      // the actual map list of sources/IDs
		mtx     *sync.RWMutex // safeguard against concurrent accesses
		µChange uint32        // internal change flag
		µCC     tCountCache   // cache for `CountedList()`
	}
)

var (
	// UseBinaryStorage determines whether to use binary storage
	// or not (i.e. plain text).
	//
	// Loading/storing binary data is about three times as fast with
	// the `THashList` data than reading and parsing plain text data.
	UseBinaryStorage = true
)

// `add()` appends 'aID` to the list.
//
//	`aID` the source ID to add to the list.
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
//	`aID` is the string to look up.
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
//	`aID` is the string to look up.
func (sl *tSourceList) removeID(aID string) *tSourceList {
	idx := sl.indexOf(aID)
	if 0 > idx {
		return sl
	}

	sLen := len(*sl) - 1
	if 0 > sLen {
		// can't remove from empty list …
		return sl
	}

	switch idx {
	case 0:
		if 0 == sLen {
			*sl = (*sl)[:0]
		} else {
			*sl = (*sl)[1:]
		}
	case sLen:
		*sl = (*sl)[:sLen]
	default:
		*sl = append((*sl)[:idx], (*sl)[idx+1:]...)
	}

	return sl
} // removeID()

// `renameID()` replaces all occurrences of `aOldID` by `aNewID`.
//
// This method is intended for rare cases when the ID of a document
// gets changed.
//
//	`aOldID` is to be replaced in this list.
//	`aNewID` is the replacement in this list.
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

// String returns the list as a linefeed separated string.
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
// ignored (i.e. this method does nothing) returning the current list.
//
//	`aDelim` is the start character of words to use (i.e. either '@' or '#').
//	`aMapIdx` is the list index to lookup.
//	`aID` is to be added to the hash list.
func (hl *THashList) add(aDelim byte, aMapIdx, aID string) *THashList {
	// the mutex.Lock is done by the callers

	if (0 == len(aMapIdx)) || (0 == len(aID)) {
		return hl
	}
	aMapIdx = strings.ToLower(aMapIdx) // prepare for case-insensitive search
	if aMapIdx[0] != aDelim {
		aMapIdx = string(aDelim) + aMapIdx
	}

	return hl.add0(aMapIdx, aID)
} // add()

// `add0()` appends `aID` to the tag list associated with `aMapIdx`.
//
//	`aMapIdx` is the list index to lookup.
//	`aID` is to be added to the hash list.
func (hl *THashList) add0(aMapIdx, aID string) *THashList {
	// the mutex.Lock is done by the callers

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

// `checksum()` returns the list's CRC32 checksum.
func (hl *THashList) checksum() uint32 {
	// the mutex.Lock is done by the callers

	if 0 == atomic.LoadUint32(&hl.µChange) {
		// We use `string()` because it sorts internally
		// thus generating reproducible results:
		atomic.StoreUint32(&hl.µChange, crc32.Update(0, crc32.MakeTable(crc32.Castagnoli), []byte(hl.string())))
	}

	return hl.µChange
} // checksum()

// Checksum returns the list's CRC32 checksum.
//
// This method can be used to get a kind of 'footprint'.
func (hl *THashList) Checksum() uint32 {
	hl.mtx.RLock()
	defer hl.mtx.RUnlock()

	return hl.checksum()
} // Checksum()

// `clear()` empties the internal data structures:
// all `#hashtags` and `@mentions` are deleted.
func (hl *THashList) clear() *THashList {
	// the mutex.Lock is done by the callers
	for mapIdx, sl := range hl.hl {
		sl.clear()
		delete(hl.hl, mapIdx)
	}
	atomic.StoreUint32(&hl.µChange, 0)

	return hl
} // clear()

// Clear empties the internal data structures:
// all `#hashtags` and `@mentions` are deleted.
func (hl *THashList) Clear() *THashList {
	hl.mtx.Lock()
	defer hl.mtx.Unlock()

	return hl.clear()
} // Clear()

// `count()` returns the number of hashtags (if `aDelim == '#'`) or
// mentions (if `aDelim == '@'`).
//
//	`aDelim` is the start of words to search (i.e. either '@' or '#').
func (hl *THashList) count(aDelim byte) (rLen int) {
	hl.mtx.RLock()
	defer hl.mtx.RUnlock()

	for idx := range hl.hl {
		if idx[0] == aDelim {
			rLen++
		}
	}

	return
} // count()

// CountedList returns a list of #hashtags/@mentions with
// their respective count of associated IDs.
func (hl *THashList) CountedList() (rList []TCountItem) {
	hl.mtx.Lock()
	defer hl.mtx.Unlock()

	if (hl.checksum() == hl.µCC.µCRC) && (0 < len(hl.µCC.µCounts)) {
		return hl.µCC.µCounts
	}

	hl.µCC.µCounts = nil
	hl.µCC.µCRC = hl.checksum()
	rList = make(tCountList, 0, len(hl.hl))
	for mapIdx, sl := range hl.hl {
		rList = append(rList, TCountItem{len(*sl), mapIdx})
	}
	if 0 < len(rList) {
		sort.Slice(rList, func(i, j int) bool {
			// ignore [#@] for sorting
			return (rList[i].Tag[1:] < rList[j].Tag[1:])
		})
	}
	hl.µCC.µCounts = rList

	return
} // CountedList()

// Filename returns the configured filename for reading/storing this list.
func (hl *THashList) Filename() string {
	hl.mtx.RLock()
	defer hl.mtx.RUnlock()

	return hl.fn
} // Filename()

// HashAdd appends `aID` to the list of `aHash`.
//
// If either `aHash` or `aID` are empty strings they are silently
// ignored (i.e. this method does nothing).
//
//	`aHash` is the list index to lookup.
//	`aID` is to be added to the hash list.
func (hl *THashList) HashAdd(aHash, aID string) *THashList {
	if (0 == len(aHash)) || (0 == len(aID)) {
		return hl
	}
	hl.mtx.Lock()
	defer hl.mtx.Unlock()

	return hl.add('#', aHash, aID)
} // HashAdd()

// HashCount returns the number of hashtags in the list.
func (hl *THashList) HashCount() int {
	return hl.count('#')
} // HashCount()

// HashLen returns the number of IDs stored for `aHash`.
//
//	`aHash` identifies the ID list to lookup.
func (hl *THashList) HashLen(aHash string) int {
	return hl.idxLen('#', aHash)
} // HashLen()

// HashList returns a list of IDs associated with `aHash`.
//
//	`aHash` identifies the ID list to lookup.
func (hl *THashList) HashList(aHash string) []string {
	return hl.list('#', aHash)
} // HashList()

// HashRemove deletes `aID` from the list of `aHash`.
//
//	`aHash` identifies the ID list to lookup.
//	`aID` is the source to remove from the list.
func (hl *THashList) HashRemove(aHash, aID string) *THashList {
	return hl.remove('#', aHash, aID)
} // HashRemove()

// IDlist returns a list of #hashtags and @mentions associated with `aID`.
func (hl *THashList) IDlist(aID string) (rList []string) {
	if 0 == len(aID) {
		return
	}
	hl.mtx.RLock()
	defer hl.mtx.RUnlock()

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

// IDparse checks whether `aText` contains strings starting with `[@|#]`
// and – if found – adds them to the respective list.
//
//	`aID` is the ID to add to the list.
//	`aText` is the text to search.
func (hl *THashList) IDparse(aID string, aText []byte) *THashList {
	if (0 == len(aID)) || (0 == len(aText)) {
		return hl
	}
	hl.mtx.Lock()
	defer hl.mtx.Unlock()

	oldCRC := hl.checksum()
	defer func() {
		if oldCRC != atomic.LoadUint32(&hl.µChange) {
			_, _ = hl.store()
		}
	}()

	return hl.parseID(aID, aText)
} // IDparse()

// IDremove deletes all @hashtags/@mentions associated with `aID`.
//
//	`aID` is to be deleted from all lists.
func (hl *THashList) IDremove(aID string) *THashList {
	if 0 == len(aID) {
		return hl
	}
	hl.mtx.Lock()
	defer hl.mtx.Unlock()

	return hl.removeID(aID)
} // IDremove()

// IDrename replaces all occurrences of `aOldID` by `aNewID`.
//
// This method is intended for rare cases when the ID of a document
// needs to get changed.
//
//	`aOldID` is to be replaced in all lists.
//	`aNewID` is the replacement in all lists.
func (hl *THashList) IDrename(aOldID, aNewID string) *THashList {
	if (0 == len(aOldID)) || (0 == len(aNewID)) {
		return hl
	}
	hl.mtx.Lock()
	defer hl.mtx.Unlock()

	for _, sl := range hl.hl {
		sl.renameID(aOldID, aNewID)
	}
	atomic.StoreUint32(&hl.µChange, 0)
	_, _ = hl.store()

	return hl
} // IDrename()

// IDupdate checks `aText` removing all #hashtags/@mentions no longer
// present and adding #hashtags/@mentions new in `aText`.
//
//	`aID` is the ID to update.
//	`aText` is the text to use.
func (hl *THashList) IDupdate(aID string, aText []byte) *THashList {
	if (0 == len(aID)) || (0 == len(aText)) {
		return hl
	}
	hl.mtx.Lock()
	defer hl.mtx.Unlock()

	oldCRC := hl.checksum()
	defer func() {
		if oldCRC != atomic.LoadUint32(&hl.µChange) {
			_, _ = hl.store()
		}
	}()

	return hl.updateID(aID, aText)
} // IDupdate()

// `idxLen()` returns the number of IDs stored for `aMapIdx`.
//
//	`aDelim` is the first character of words to use (i.e. either '@' or '#').
//	`aMapIdx` identifies the ID list to lookup.
func (hl *THashList) idxLen(aDelim byte, aMapIdx string) int {
	if 0 == len(aMapIdx) {
		return -1
	}
	hl.mtx.RLock()
	defer hl.mtx.RUnlock()

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
	hl.mtx.RLock()
	defer hl.mtx.RUnlock()

	return len(hl.hl)
} // Len()

// LenTotal returns the length of all #hashtag/@mention lists together.
func (hl *THashList) LenTotal() (rLen int) {
	hl.mtx.RLock()
	defer hl.mtx.RUnlock()

	rLen = len(hl.hl)
	for _, sl := range hl.hl {
		rLen += len(*sl)
	}

	return
} // LenTotal()

// `list()` returns a list of IDs associated with `aMapIdx`.
//
//	`aDelim` The start of words to search (i.e. either '@' or '#').
//	`aMapIdx` Name of the sources list to lookup.
func (hl *THashList) list(aDelim byte, aMapIdx string) (rList []string) {
	if 0 == len(aMapIdx) {
		return
	}
	// We need a R/W lock here since we change the used `tSourceList`.
	hl.mtx.Lock()
	defer hl.mtx.Unlock()

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

// Load reads the configured file returning the data structure
// read from the file and a possible error condition.
//
// If the hash file doesn't exist that is not considered an error.
// If there is an error, it will be of type `*PathError`.
func (hl *THashList) Load() (*THashList, error) {
	hl.mtx.Lock()
	defer hl.mtx.Unlock()

	file, err := os.OpenFile(hl.fn, os.O_RDONLY, 0)
	if nil != err {
		if os.IsNotExist(err) {
			return hl, nil
		}
		return hl, err
	}
	defer file.Close()
	if UseBinaryStorage {
		return hl.loadBinary(file)
	}

	return hl.loadText(file)
} // Load()

// `loadBinary()` reads a file written by `store()` returning
// the modified list and a possible error.
func (hl *THashList) loadBinary(aFile *os.File) (*THashList, error) {
	// The mutex.Lock is done by the caller

	var decodedMap tHashMap
	decoder := gob.NewDecoder(aFile)
	if err := decoder.Decode(&decodedMap); err != nil {
		// `decoder.Decode()` returns `io.EOF` if the input
		// is at EOF which we do not consider an error here.
		if !errors.Is(err, io.EOF) && !errors.Is(err, io.ErrUnexpectedEOF) {
			return hl, err
		}
		decodedMap = make(tHashMap, 64)
	}
	hl.hl = decodedMap

	return hl, nil
} // loadBinary()

// `loadText()` parses a file written by `store()` returning
// the modified list and a possible error.
//
// This method reads one line of the file at a time.
func (hl *THashList) loadText(aFile *os.File) (*THashList, error) {
	// The mutex.Lock is done by the caller
	var (
		mapIdx string
		rRead  int
	)
	scanner := bufio.NewScanner(aFile)
	hl.clear()
	for lineRead := scanner.Scan(); lineRead; lineRead = scanner.Scan() {
		line := scanner.Text()
		rRead += len(line) + 1 // add trailing LF

		line = strings.TrimSpace(line)
		if 0 == len(line) {
			continue
		}

		if matches := htHashHeadRE.FindStringSubmatch(line); nil != matches {
			mapIdx = strings.ToLower(strings.TrimSpace(matches[1]))
		} else {
			hl.add0(mapIdx, line)
		}
	}
	atomic.StoreUint32(&hl.µChange, 0)

	return hl, scanner.Err()
} // loadText()

// MentionAdd appends `aID` to the list of `aMention`.
//
// If either `aMention` or `aID` are empty strings they are
// silently ignored (i.e. this method does nothing).
//
//	`aMention` is the list index to lookup.
//	`aID` is to be added to the hash list.
func (hl *THashList) MentionAdd(aMention, aID string) *THashList {
	if (0 == len(aMention)) || (0 == len(aID)) {
		return hl
	}
	hl.mtx.Lock()
	defer hl.mtx.Unlock()

	return hl.add('@', aMention, aID)
} // MentionAdd()

// MentionCount returns the number of mentions in the list.
func (hl *THashList) MentionCount() int {
	return hl.count('@')
} // MentionCount()

// MentionLen returns the number of IDs stored for `aMention`.
//
//	`aMention` identifies the ID list to lookup.
func (hl *THashList) MentionLen(aMention string) int {
	return hl.idxLen('@', aMention)
} // MentionLen()

// MentionList returns a list of IDs associated with `aMention`.
//
//	`aMention` identifies the ID list to lookup.
func (hl *THashList) MentionList(aMention string) []string {
	return hl.list('@', aMention)
} // MentionList()

// MentionRemove deletes `aID` from the list of `aMention`.
//
//	`aMention` identifies the ID list to lookup.
//	`aID` is the source to remove from the list.
func (hl *THashList) MentionRemove(aMention, aID string) *THashList {
	return hl.remove('@', aMention, aID)
} // MentionRemove()

var (
	// RegEx to identify a numeric HTML entity.
	htEntityRE = regexp.MustCompile(`#[0-9]+;`)

	// match: [#Hashtag|@Mention]
	htHashHeadRE = regexp.MustCompile(`^\[\s*([#@][^\]]*?)\s*\]$`)
	//                                        11111111111

	// match: #hashtag|@mention
	htHashMentionRE = regexp.MustCompile(
		`(?ims)(?:^|\s|[^\p{L}\d_])?([@#][\p{L}’'\d_§-]+)(?:[^\p{L}\d_]|$)`)
	//	                             1111111111111111111  222222222222222

	// RegEx to match texts like `#----`.
	htHyphenRE = regexp.MustCompile(`#[^-]*--`)
)

// `parseID()` checks whether `aText` contains strings starting
// with `[@|#]` and – if found – adds them to the respective list.
//
//	`aID` is the ID to add to the list.
//	`aText` is the text to search.
func (hl *THashList) parseID(aID string, aText []byte) *THashList {
	// The mutex.Lock is done by the caller
	matches := htHashMentionRE.FindAllSubmatch(aText, -1)
	if (nil == matches) || (0 == len(matches)) {
		return hl
	}
	for _, sub := range matches {
		match0 := string(sub[0])
		hash := string(sub[1])
		if '_' == hash[len(hash)-1] {
			// '_' can be both, part of the hashtag and italic
			// markup so we must remove it if it's at the end:
			hash = hash[:len(hash)-1]
		}
		if '#' == hash[0] {
			// `match0` is the match including prefix and postfix
			switch match0[len(match0)-1] {
			case '"':
				// Double quote following a possible hashtag:
				// most probably an URL#fragment, so check
				// whether it's a quoted string:
				if '"' != match0[0] {
					continue // URL#fragment
				}
			case ')':
				// This is a tricky one: it can either be a
				// normal right round bracket or the end of
				// a Markdown link. Here we assume that it's
				// the latter one and ignore this match:
				continue
			case '-':
				// A hyphen at the end of a hashtag:
				// that's not part of an acceptable tag.
				continue
			case ';':
				if htEntityRE.MatchString(match0) {
					// leave HTML entities as is
					continue
				}
			}
			if htHyphenRE.MatchString(hash) {
				continue
			}
		}
		hl.add(hash[0], hash, aID)
	}

	return hl
} // parseID()

// `remove()` deletes `aID` from the list of `aMapIdx`.
//
//	`aDelim` is the start character of words to use (i.e. either '@' or '#').
//	`aMapIdx` identifies the sources list to lookup.
//	`aID` is the source to remove from the list.
func (hl *THashList) remove(aDelim byte, aMapIdx, aID string) *THashList {
	if (0 == len(aMapIdx)) || (0 == len(aID)) {
		return hl
	}
	hl.mtx.Lock()
	defer hl.mtx.Unlock()

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

// `removeID()` deletes all #hashtags/@mentions associated with `aID`.
//
//	`aID` is to be deleted from all lists.
func (hl *THashList) removeID(aID string) *THashList {
	// The mutex.Lock is done by the callers

	oldCRC := hl.checksum()
	defer func() {
		if oldCRC != atomic.LoadUint32(&hl.µChange) {
			_, _ = hl.store()
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

// SetFilename sets `aFilename` to use by this list.
func (hl *THashList) SetFilename(aFilename string) *THashList {
	hl.mtx.Lock()
	defer hl.mtx.Unlock()

	hl.fn = aFilename

	return hl
} // SetFilename()

// `store()` writes the whole list to the configured file
// returning the number of bytes written and a possible error.
//
// If there is an error, it will be of type `*PathError`.
func (hl *THashList) store() (int, error) {
	// the mutex.Lock is done by the callers

	file, err := os.OpenFile(hl.fn, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0660) //#nosec G302
	if nil != err {
		return 0, err
	}
	defer file.Close()

	if UseBinaryStorage {
		encoder := gob.NewEncoder(file)
		if err = encoder.Encode(hl.hl); nil != err {
			return 0, err
		}
		size, err := file.Seek(0, os.SEEK_END)

		return int(size), err
	}

	return file.Write([]byte(hl.string()))
} // store()

// Store writes the whole list to the configured file
// returning the number of bytes written and a possible error.
//
// If there is an error, it will be of type `*PathError`.
func (hl *THashList) Store() (int, error) {
	hl.mtx.RLock()
	defer hl.mtx.RUnlock()

	return hl.store()
} // Store()

// `string()` returns the whole list as a newline separated string.
func (hl *THashList) string() string {
	// the mutex.Lock is done by the caller

	tmp := make(tSourceList, 0, len(hl.hl))
	for hash := range hl.hl {
		tmp = append(tmp, hash)
	}
	// sort the order of hashtags to get a reproducible result
	sort.Slice(tmp, func(i, j int) bool {
		// ignore leading [@#] when sorting
		return (tmp[i][1:] < tmp[j][1:]) // ascending
	})

	var result string
	for _, hash := range tmp {
		sl := hl.hl[hash]
		result += "[" + hash + "]\n" + sl.String() + "\n"
	}

	return result
} // string()

// String returns the whole list as a linefeed separated string.
func (hl *THashList) String() string {
	hl.mtx.RLock()
	defer hl.mtx.RUnlock()

	return hl.string()
} // String()

// `updateID()` checks `aText` removing all #hashtags/@mentions no longer
// present and adds #hashtags/@mentions new in `aText`.
//
//	`aID` is the ID to update.
//	`aText` is the text to use.
func (hl *THashList) updateID(aID string, aText []byte) *THashList {
	// the mutex.Lock is done by the caller
	hl.removeID(aID)

	return hl.parseID(aID, aText)
} // updateID()

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
	//	`aHash` is the list index to lookup.
	//	`aID` is to be added to the hash list.
	THashWalker interface {
		Walk(aHash, aID string) bool
	}
)

// Walk traverses through all entries in the #hashtag/@mention lists
// calling `aFunc` for each entry.
//
// If `aFunc` returns `false` when called the respective ID
// will be removed from the associated #hashtag/@mention.
//
//	`aFunc` is the function called for each ID in all lists.
func (hl *THashList) Walk(aFunc TWalkFunc) {
	oldCRC := hl.checksum()
	defer func() {
		if oldCRC != atomic.LoadUint32(&hl.µChange) {
			_, _ = hl.Store()
		}
	}()

	for hash, sl := range hl.hl {
		for _, id := range *sl {
			if aFunc(hash, id) {
				continue
			}
			sl.removeID(id)
			if 0 == len(*sl) {
				delete(hl.hl, hash)
			}
			atomic.StoreUint32(&hl.µChange, 0)
		}
	}
} // Walk()

// Walker traverses through all entries in the INI list sections
// calling `aWalker` for each entry.
//
//	`aWalker` is an object implementing the `TIniWalker` interface.
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
//	`aFilename` is the name of the file to use for reading and storing.
func New(aFilename string) (*THashList, error) {
	result := &THashList{
		fn:  aFilename,
		hl:  make(tHashMap, 64),
		mtx: new(sync.RWMutex),
	}

	return result.Load()
} // New()

/* EoF */
