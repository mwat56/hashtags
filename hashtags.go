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
	// Simple type alias
	tHash = string

	// tSourceList is a slice of strings
	tSourceList []string

	// A map indexed by `tHash` holding a `tStrList`
	tHashMap map[tHash]*tSourceList

	// THashList is a list of hashes pointing to sources
	THashList tHashMap
)

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
} // Clear()

// `indexOf()` returns the list indexOf of `aID`.
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

// `remove()` deletes the source at index `aIdx`.
//
// `aIdx` the list index of the source to delete.
func (sl *tSourceList) remove(aIdx int) *tSourceList {
	slen := len(*sl) - 1
	if 0 > slen {
		// can't remove from empty list …
		return sl
	}
	if 0 == aIdx {
		*sl = (*sl)[1:]
	} else if slen == aIdx {
		*sl = (*sl)[:slen]
	} else {
		*sl = append((*sl)[:aIdx], (*sl)[aIdx+1:]...)
	}

	return sl
} // remove()

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

// Add appends `aID` to the list indexed by `aHash`.
//
// If either `aHash` or `aID` are empty strings htey are silently
// ignored (i.a. this method does nothing).
//
// `aHash` is the list index to lookup.
//
// `aID` is to be added to the hash list.
func (hl *THashList) Add(aHash, aID string) *THashList {
	if (0 == len(aHash)) || (0 == len(aID)) {
		return hl
	}
	if sl, ok := (*hl)[aHash]; ok {
		(*hl)[aHash] = sl.add(aID)
	} else {
		sl := make(tSourceList, 1)
		sl[0] = aID
		(*hl)[aHash] = &sl
	}

	return hl
} // Add()

// Clear empties the internal data structures.
//
// This method can be called once the program has used the config values
// stored in the INI file to setup the application. Emptying these data
// structures should help the garbage collector do release the data not
// needed anymore.
func (hl *THashList) Clear() bool {
	for hash, sl := range *hl {
		sl.clear()
		delete(*hl, hash)
	}

	return (0 == len(*hl))
} // Clear()

// DeleteSource removes `aID` from the list of `aHash`.
//
// `aHash` identifies the sources list to lookup.
//
// `aID` is the source to remove from the list.
func (hl *THashList) DeleteSource(aHash, aID string) *THashList {
	if (0 == len(aHash)) || (0 == len(aID)) {
		return hl
	}
	sl, ok := (*hl)[aHash]
	if !ok {
		return hl
	}

	if idx := sl.indexOf(aID); 0 <= idx {
		sl.remove(idx)
	}

	return hl
} // DeleteSource()

// HashLen returns the number of sources stored for `aHash`.
//
// `aHash` identifies the sources list to lookup.
func (hl *THashList) HashLen(aHash string) int {
	if sl, ok := (*hl)[aHash]; ok {
		return len(*sl)
	}

	return -1
} // HashLen()

// HashList returns a list of strings associates with 'aHash`.
//
// `aHash` identifies the sources list to lookup.
func (hl *THashList) HashList(aHash string) (rList []string) {
	if sl, ok := (*hl)[aHash]; ok {
		rList = []string(*sl)
	}

	return
} // HashList()

// HashParse searches `aText` for #hashtags and if found
// adds it with `aID` to the list.
//
// `aID` is the ID to add to the list.
//
// `aText` is the text to search.
func (hl *THashList) HashParse(aID string, aText []byte) *THashList {
	return hl.parse('#', aID, aText)
} // HashParse()

// Len returns the current length of the list.
func (hl *THashList) Len() int {
	return len(*hl)
} // Len()

// Load reads the given `aFilename` returning the data structure
// read from the file and a possible error condition.
//
// This method reads one line at a time of the INI file skipping both
// empty lines and comments (identified by '#' or ';' at line start).
//
// `aFilename` is the name of the INI file to read.
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

// MentionParse searches `aText` for @mentions and if found
// adds it with `aID` to the list.
//
// `aID` is the ID to add to the list.
//
// `aText` is the text to search.
func (hl *THashList) MentionParse(aID string, aText []byte) *THashList {
	return hl.parse('@', aID, aText)
} // MentionParse()

// `parse()` checks whether `aText` contains strings starting with `aDelim`
// and – if found – adds it to the list.
//
// `aDelim` is the start of words to search (i.e. either '@' or '#').
//
// `aID` is the ID to add to the list.
//
// `aText` is the text to search.
func (hl *THashList) parse(aDelim rune, aID string, aText []byte) *THashList {
	re, err := regexp.Compile(`(?s)\W(\` + string(aDelim) + `\w+)`)
	if nil != err {
		return hl
	}
	matches := re.FindAllSubmatch(aText, -1)
	if (nil == matches) || (0 >= len(matches)) {
		return hl
	}
	for _, sub := range matches {
		if 0 < len(sub[1]) {
			hl = hl.Add(string(sub[1]), aID)
		}
	}

	return hl
} // parse()

var (
	// match: [section]
	sectionRE = regexp.MustCompile(`^\[\s*([^\]]*?)\s*]$`)
)

// `read()` parses a file written by `Store()` returning the
// number of bytes read and a possible error.
//
// This method reads one line of the INI file at a time skipping both
// empty lines and comments (identified by '#' or ';' at line start).
func (hl *THashList) read(aScanner *bufio.Scanner) (rRead int, rErr error) {
	var section string

	for lineRead := aScanner.Scan(); lineRead; lineRead = aScanner.Scan() {
		line := aScanner.Text()
		rRead += len(line) + 1 // add trailing LF

		line = strings.TrimSpace(line)
		if 0 == len(line) {
			continue
		}

		if matches := sectionRE.FindStringSubmatch(line); nil != matches {
			section = strings.TrimSpace(matches[1])
		} else {
			hl.Add(section, line)
		}
	}
	rErr = aScanner.Err()

	return
} // read()

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

// NewList returns a new `THashList` instance.
func NewList() *THashList {
	result := make(THashList, 32)

	return &result
} // NewList()

/* EoF */
