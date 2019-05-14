/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
                  All rights reserved
              EMail : <support@mwat.de>
*/

package hashtags

import (
	"bufio"
	"os"
	"regexp"
	"sort"
	"strings"
)

type (
	// `tSourceList` is a slice of strings
	tSourceList []string
)

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

// `add()` appends 'aID` to the list
//
// `aID` the source ID to add to the list.
func (sl *tSourceList) add(aID string) *tSourceList {
	for _, s := range *sl {
		if s == aID {
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

type (
	// A map indexed by a string holding a `tStrList`
	tHashMap map[string]*tSourceList

	// THashList is a list of `#hashtags` and `@mentions`
	// pointing to sources (occurances).
	THashList tHashMap
)

var (
	// match: [#Hashtag|@mention]
	hashHeadRE = regexp.MustCompile(`^\[\s*([#@][^\]]*?)\s*\]$`)

	// match: #hashtag|@mention
	hashMentionRE = regexp.MustCompile(`(^|\W)([\@\#]\w+)(\W|$)`)
)

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
	if sl, ok := (*hl)[aMapIdx]; ok {
		(*hl)[aMapIdx] = sl.add(aID).sort()
	} else {
		sl := make(tSourceList, 1, 32)
		sl[0] = aID
		(*hl)[aMapIdx] = &sl
	}

	return hl
} // add0()

// Clear empties the internal data structures:
// all `#hashtags` and `@mentions` are deleted.
func (hl *THashList) Clear() *THashList {
	for mapIdx, sl := range *hl {
		sl.clear()
		delete(*hl, mapIdx)
	}

	return hl
} // Clear()

// HashAdd appends `aID` to the list indexed by `aHash`.
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

// HashLen returns the number of sources stored for `aHash`.
//
// `aHash` identifies the sources list to lookup.
func (hl *THashList) HashLen(aHash string) int {
	return hl.idxLen('#', aHash)
} // HashLen()

// HashList returns a list of IDs associated with `aHash`.
//
// `aHash` identifies the sources list to lookup.
func (hl *THashList) HashList(aHash string) []string {
	return hl.list('#', aHash)
} // HashList()

// HashRemove deletes `aID` from the list of `aHash`.
//
// `aHash` identifies the sources list to lookup.
//
// `aID` is the source to remove from the list.
func (hl *THashList) HashRemove(aHash, aID string) *THashList {
	return hl.remove('#', aHash, aID)
} // HashRemove()

// IDlist returns a list of #hashtags and @mentions associated with `aID`.
func (hl *THashList) IDlist(aID string) (rList []string) {
	for mapIdx, sl := range *hl {
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
// and – if found – adds them to the list.
//
// `aID` is the ID to add to the list.
//
// `aText` is the text to search.
func (hl *THashList) IDparse(aID string, aText []byte) *THashList {
	matches := hashMentionRE.FindAllSubmatch(aText, -1)
	if (nil == matches) || (0 >= len(matches)) {
		return hl
	}
	for _, sub := range matches {
		if 0 < len(sub[2]) {
			hl.add0(strings.ToLower(string(sub[2])), aID)
		}
	}

	return hl
} // IDparse()

// IDremove deletes all @hashtags/@mentions associated with `aID`.
//
// `aID` is to be deleted from all lists.
func (hl *THashList) IDremove(aID string) *THashList {
	for mapIdx, sl := range *hl {
		sl.removeID(aID)
		if 0 == len(*sl) {
			delete(*hl, mapIdx)
		}
	}

	return hl
} // IDremove()

// IDrename replaces all occurances of `aOldID` by `aNewID`.
//
// This method is intended for rare cases when the ID of a document
// gets changed.
//
// `aOldID` is to be replaced in all lists.
//
// `aNewID` is the replacement in all lists.
func (hl *THashList) IDrename(aOldID, aNewID string) *THashList {
	for _, sl := range *hl {
		sl.renameID(aOldID, aNewID)
	}

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

	matches := hashMentionRE.FindAllSubmatch(aText, -1)
	if (nil == matches) || (0 >= len(matches)) {
		return hl
	}
	for _, sub := range matches {
		if 0 < len(sub[2]) {
			hl = hl.add(sub[2][0], string(sub[2]), aID)
		}
	}

	return hl
} // IDupdate()

// `idxLen()` returns the number of sources stored for `aMapIdx`.
//
// `aDelim` is the start character of words to use (i.e. either '@' or '#').
//
// `aMapIdx` identifies the sources list to lookup.
func (hl *THashList) idxLen(aDelim byte, aMapIdx string) int {
	if 0 == len(aMapIdx) {
		return -1
	}
	aMapIdx = strings.ToLower(aMapIdx)
	if aMapIdx[0] != aDelim {
		aMapIdx = string(aDelim) + aMapIdx
	}
	if sl, ok := (*hl)[aMapIdx]; ok {
		return len(*sl)
	}

	return -1
} // idxLen()

// Len returns the current length of the list i.e. how many #hashtags
// and @mentions are currently stored in the list.
func (hl *THashList) Len() int {
	return len(*hl)
} // Len()

// LenTotal returns the length of all #hashtag/@mention lists together.
func (hl *THashList) LenTotal() int {
	result := len(*hl)
	for _, sl := range *hl {
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
	if sl, ok := (*hl)[aMapIdx]; ok {
		sl.sort()
		rList = []string(*sl)
	}

	return
} // list()

// Load reads the given `aFilename` returning the data structure
// read from the file and a possible error condition.
//
// If there is an error, it will be of type *PathError.
//
// `aFilename` is the name of the file to read.
func (hl *THashList) Load(aFilename string) (*THashList, error) {
	file, err := os.Open(aFilename)
	if nil != err {
		return hl, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	_, err = hl.read(scanner)
	return hl, err
} // Load()

// MentionAdd appends `aID` to the list indexed by `aMention`.
//
// If either `aMention` or `aID` are empty strings they are silently
// ignored (i.e. this method does nothing).
//
// `aMention` is the list index to lookup.
//
// `aID` is to be added to the hash list.
func (hl *THashList) MentionAdd(aMention, aID string) *THashList {
	return hl.add('@', aMention, aID)
} // MentionAdd()

// MentionLen returns the number of sources stored for `aMention`.
//
// `aMention` identifies the sources list to lookup.
func (hl *THashList) MentionLen(aMention string) int {
	return hl.idxLen('@', aMention)
} // MentionLen()

// MentionList returns a list of IDs associated with `aHash`.
//
// `aMention` identifies the sources list to lookup.
func (hl *THashList) MentionList(aMention string) []string {
	return hl.list('@', aMention)
} // MentionList()

// MentionRemove deletes `aID` from the list of `aMention`.
//
// `aMention` identifies the sources list to lookup.
//
// `aID` is the source to remove from the list.
func (hl *THashList) MentionRemove(aMention, aID string) *THashList {
	return hl.remove('@', aMention, aID)
} // MentionRemove()

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
	if sl, ok := (*hl)[aMapIdx]; ok {
		sl.removeID(aID)
		if 0 == len(*sl) {
			delete(*hl, aMapIdx)
		}
	}

	return hl
} // remove()

// Store writes the whole list to `aFilename`
// returning the number of bytes written and a possible error.
//
// `aFilename` is the name of the file to write.
func (hl *THashList) Store(aFilename string) (int, error) {
	s := hl.String()
	file, err := os.Create(aFilename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	return file.Write([]byte(s))
} // Store()

// String returns the whole list as a linefeed separated string.
func (hl *THashList) String() string {
	var (
		result string
		tmp    []string
	)
	for hash := range *hl {
		tmp = append(tmp, hash)
	}
	// sort the order of hashtags to get a reproducible result
	sort.Slice(tmp, func(i, j int) bool {
		return (tmp[i] < tmp[j]) // ascending
	})
	for _, hash := range tmp {
		sl, _ := (*hl)[hash]
		result += "[" + hash + "]\n" + sl.String() + "\n"
	}

	return result
} // String()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

// LoadList returns a new `THashList` instance after reading
// the given file.
//
// If there is an error, it will be of type *PathError.
//
// `aFilename` is the name of the file to read.
func LoadList(aFilename string) (*THashList, error) {
	result := NewList()

	return result.Load(aFilename)
} // LoadList()

// NewList returns a new `THashList` instance.
func NewList() *THashList {
	result := make(THashList, 32)

	return &result
} // NewList()

/* EoF */
