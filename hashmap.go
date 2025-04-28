/*
Copyright © 2019, 2025  M.Watermann, 10247 Berlin, Germany

	....All rights reserved
	EMail : <support@mwat.de>
*/
package hashtags

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"

	se "github.com/mwat56/sourceerror"
)

//lint:file-ignore ST1017 - I prefer Yoda conditions

type (
	// `tHashMap` is indexed by `#hashtags` and `@mentions`
	// pointing to a `tSourceList`.
	tHashMap map[string]*tSourceList
)

// --------------------------------------------------------------------------
// constructor function

func newHashMap() *tHashMap {
	hm := make(tHashMap, 64)

	return &hm
} // newHashMap()

// --------------------------------------------------------------------------
// helper for sorting the hash strings; used by `keys()`

// `cmp4sort()` helps to sort a slice in ascending order based on the
// comparison of the substrings after the leading hash or mention mark.
//
// It is used in the `sort()` method to ensure that the hash map is sorted,
// which can improve the performance of certain operations on the hash map,
// such as searching for a specific key.
//
// The function takes two strings as input and returns an integer indicating
// their relative order. If `a` is less than `b`, it returns `-1`. If a is
// greater than `b`, it returns `1`. If `a` and `b` are equal, it returns `0`.
//
// The function first checks if the first character of `a` and `b` is a hash
// or mention mark. If so, it removes the leading character from both strings.
//
// Returns:
//   - `int`: The result of the comparison the two strings as described above.
func cmp4sort(a, b string) int {
	switch a[0] {
	case MarkHash, MarkMention:
		a = a[1:]
	}

	switch b[0] {
	case MarkHash, MarkMention:
		b = b[1:]
	}

	if a < b {
		return -1
	}
	if a > b {
		return 1
	}

	return 0
} // cmp4sort()

// -------------------------------------------------------------------------
// methods of `tHashMap`

var (
	// `gCRCtable` should be considered `const`.
	gCRCtable = crc32.MakeTable(0) // other table types may not be reproducible
)

// `checksum()` computes the list's CRC32 checksum.
//
// Returns:
//   - `uint32`: The computed checksum.
func (hm *tHashMap) checksum() uint32 {
	// We use `String()` because it sorts internally
	// thus generating reproducible results
	return crc32.Update(0, gCRCtable, []byte(hm.String()))
} // checksum()

// `clear()` empties the internal data structures:
// all `#hashtags` and `@mentions` are deleted.
//
// Returns:
//   - `*tHashMap`: The cleared hash map.
func (hm *tHashMap) clear() *tHashMap {
	if 0 == len(*hm) {
		return hm
	}

	var (
		hash string
		sl   *tSourceList
	)
	for hash, sl = range *hm {
		sl.clear()
		delete(*hm, hash)
	}

	return hm
} // clear()

// `count()` returns the number of `#hashtags` (if `aDelim == '#'`) or
// `@mentions` (if `aDelim == '@'`).
//
// Parameters:
//   - `aDelim`: The start of words to search (i.e. either '@' or '#').
//
// Returns:
//   - `int`: The number of `#hashtags` and `@mentions`.
func (hm *tHashMap) count(aDelim byte) int {
	if 0 == len(*hm) {
		return 0
	}
	var (
		hash   string
		result int
	)

	for hash = range *hm {
		if hash[0] == aDelim {
			result++
		}
	}

	return result
} // count()

// `countedList()` returns a list of `#hashtags` and `@mentions` with
// their respective count of associated IDs.
//
// Returns:
//   - `TCountList`: A list of `#hashtags` and `@mentions` with their respective count of associated IDs.
func (hm *tHashMap) countedList() TCountList {
	if 0 == len(*hm) {
		return nil
	}

	result := TCountList{}
	var (
		name string
		sl   *tSourceList
	)
	for name, sl = range *hm {
		result.Insert(TCountItem{len(*sl), name})
	}

	return result
} // countedList()

// `equals()` returns whether the current hash map is equal to the
// provided hash map.
//
// Parameters:
//   - `aMap`: The hash map to compare with.
//
// Returns:
//   - `bool`: Whether the hash maps are equal.
func (hm *tHashMap) equals(aMap tHashMap) bool {
	if len(*hm) != len(aMap) {
		return false
	}

	var (
		hash  string
		other *tSourceList
		sl    *tSourceList
		ok    bool
	)

	for hash, sl = range *hm {
		if other, ok = aMap[hash]; !ok {
			return false
		}
		if !sl.equals(*other) {
			return false
		}
	}

	return true
} // equals()

// `idList()` returns a list of `#hashtags` and `@mentions` associated
// with `aID`.
//
// Parameters:
//   - `aID`: The referenced object to lookup.
//
// Returns:
//   - `[]string`: List of `#hashtags` and `@mentions` associated with `aID`.
func (hm *tHashMap) idList(aID int64) []string {
	var (
		hash   string
		result []string
		sl     *tSourceList
	)
	hLen := len(*hm)
	if 0 == hLen {
		return result
	}

	result = make([]string, 0, hLen)
	for hash, sl = range *hm {
		if 0 > sl.findIndex(aID) {
			continue // ID not found
		}
		result = append(result, hash)
	}

	if 0 < len(result) {
		sort.Slice(result, func(i, j int) bool {
			return (result[i] < result[j]) // ascending
		})
	}

	return result
} // idList()

// `idxLen()` returns the number of IDs stored for `aName`.
//
// If `aName` is empty it is silently ignored (i.e. this method
// does nothing), returning `-1`.
//
// Parameters:
//   - `aDelim`: The first character of words to use (i.e. either '@' or '#').
//   - `aName`: The hash to lookup.
//
// Returns:
//   - `int: The number of references of `aName`, or `-1` if not found.
func (hm *tHashMap) idxLen(aDelim byte, aName string) int {
	if aName = strings.ToLower(strings.TrimSpace((aName))); 0 == len(aName) {
		return -1
	}

	if aName[0] != aDelim {
		aName = string(aDelim) + aName
	}

	if sl, ok := (*hm)[aName]; ok {
		return len(*sl)
	}

	return -1
} // idxLen()

// `insert()` adds `aID` to the sources list associated with `aName`.
//
// Parameters:
//   - `aName`: The list index to lookup.
//   - `aID`: The ID to be added to the hash list.
//
// Returns:
//   - `bool`: `true` if `aID` was added, or `false` otherwise.
func (hm *tHashMap) insert(aName string, aID int64) bool {
	aName = strings.ToLower(strings.TrimSpace((aName)))
	if 0 == len(aName) {
		return false
	}

	if sl, ok := (*hm)[aName]; ok {
		if sl.insert(aID) { // changes in place
			return true
		}
	} else {
		sl := newSourceList()
		if sl.insert(aID) {
			(*hm)[aName] = sl // assign the ID list to the hash
			return true
		}
	}

	return false
} // insert()

// `keys()` returns a slice of all keys in the hash map.
// If the hash map is empty, it returns an empty slice.
//
// The function does not modify the hash map.
//
// Returns:
//   - `[]string`: A sorted slice of all keys in the hash map.
func (hm *tHashMap) keys() []string {
	hLen := len(*hm)
	if 0 == hLen {
		return []string{}
	}

	var key string
	// Create a slice to hold the keys
	keys := make([]string, 0, hLen)
	for key = range *hm {
		keys = append(keys, key)
	}
	// Sort the keys ignoring the respective leading hash/mention mark
	slices.SortFunc(keys, cmp4sort)

	return keys
} // keys()

// `list()` returns a list of object IDs associated with `aName`.
//
// If `aName` is empty it is silently ignored (i.e. this method
// does nothing), returning an empty slice.
//
// Parameters:
//   - `aDelim`: The start of words to search (i.e. either '@' or '#').
//   - `aName`: The hash to lookup.
//
// Returns:
//   - `[]int64`: The number of references of `aName`.
func (hm *tHashMap) list(aDelim byte, aName string) (rList []int64) {
	aName = strings.ToLower(strings.TrimSpace((aName)))
	if "" == aName {
		return
	}
	if 0 == len(*hm) {
		return
	}

	if aName[0] != aDelim {
		aName = string(aDelim) + aName
	}

	if sl, ok := (*hm)[aName]; ok {
		rList = []int64(*sl)
	}

	return
} // list()

// `load()` reads the configured file returning the data structure
// read from the file and a possible error condition.
//
// NOTE: An empty filename or the hash file doesn't exist that
// is not considered an error.
//
// Parameters:
//   - `aFilename`: Name of the file to load.
//
// Returns:
//   - `*tHashMap`: The loaded hash map.
//   - `error`: A possible I/O error.
func (hm *tHashMap) load(aFilename string) (*tHashMap, error) {
	if aFilename = strings.TrimSpace(aFilename); "" == aFilename {
		return nil, se.New(errors.New("empty filename"), 1)
	}

	var (
		err  error
		file *os.File
	)

	file, err = os.OpenFile(aFilename, os.O_RDONLY, 0) //#nosec G304
	if nil != err {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, se.New(err, 5)
	}
	defer file.Close()

	if UseBinaryStorage {
		err = hm.loadBinary(file)
	} else {
		err = hm.loadText(file)
	}

	return hm, err
} // load()

// `loadBinary()` reads a file written by `store()` returning the modified
// list and a possible error.
//
// NOTE: This method updates the list in place.
//
// Parameters:
//   - `aFile`: The file to read from.
//
// Returns:
//   - `error`: A possible I/O error.
func (hm *tHashMap) loadBinary(aFile *os.File) error {
	iMap, iErr := loadBinaryInts(aFile)
	if nil != iErr {
		sMap, sErr := loadBinaryStrings(aFile)
		if nil == sErr {
			*hm = *sMap
			return nil
		}
		return iErr
	}
	*hm = *iMap

	return nil
} // loadBinary()

// `loadBinaryInts()` reads a binary encoded integer map from `aFile`
// and converts it into a `tHashMap`.
//
// Parameters:
//   - `aFile`: The file handle to read from.
//
// Returns:
//   - `*tHashMap`: The decoded and converted hash map.
//   - `error`: A possible decoding or conversion error.
func loadBinaryInts(aFile *os.File) (*tHashMap, error) {
	var decodedMap tHashMap

	_, _ = aFile.Seek(0, io.SeekStart)
	decoder := gob.NewDecoder(aFile)
	if err := decoder.Decode(&decodedMap); nil != err {
		// `decoder.Decode()` returns `io.EOF` if the input
		// is at EOF which we do not consider an error here.
		if !errors.Is(err, io.EOF) && !errors.Is(err, io.ErrUnexpectedEOF) {
			return nil, se.New(err, 4)
		}

		// some other error occurred
		return nil, se.New(err, 8)
	}

	return decodedMap.sort(), nil
} // loadBinaryInts()

// `loadBinaryStrings()` reads a binary encoded string map from `aFile`
// and converts it into a `tHashMap`.
//
// Parameters:
//   - `aFile`: The file handle to read from.
//
// Returns:
//   - `*tHashMap`: The decoded and converted hash map.
//   - `error`: A possible decoding or conversion error.
func loadBinaryStrings(aFile *os.File) (*tHashMap, error) {
	var decodedMap map[string][]string

	_, _ = aFile.Seek(0, io.SeekStart) //#nosec G104
	decoder := gob.NewDecoder(aFile)
	if err := decoder.Decode(&decodedMap); nil != err {
		// `decoder.Decode()` returns `io.EOF` if the input
		// is at EOF which we do not consider an error here.
		if !errors.Is(err, io.EOF) && !errors.Is(err, io.ErrUnexpectedEOF) {
			return nil, se.New(err, 4)
		}

		// some other error occurred
		return nil, se.New(err, 8)
	}

	result := newHashMap()
	var (
		key  string
		sArr []string
		i64  int64
		err  error
	)
	for key, sArr = range decodedMap {
		for _, str := range sArr {
			if i64, err = strconv.ParseInt(str, 16, 64); nil == err {
				result.insert(key, i64)
			}
		}
	}

	return result, nil
} // loadBinaryStrings()

// `loadText()` parses a file written by `store()` returning a possible error.
//
// This method reads one line of the file at a time.
//
// NOTE: This method updates the list in place.
//
// Parameters:
//   - `aFile`: The file to read from.
//
// Returns:
//   - `error`: A possible I/O error.
func (hm *tHashMap) loadText(aFile *os.File) error {
	hm.clear()

	var (
		err      error
		hash     string
		line     string
		matches  []string
		read     int
		i64      int64
		lineRead bool
	)
	scanner := bufio.NewScanner(aFile)
	for lineRead = scanner.Scan(); lineRead; lineRead = scanner.Scan() {
		line = scanner.Text()
		read += len(line) + 1 // add trailing LF

		line = strings.TrimSpace(line)
		if 0 == len(line) {
			continue
		}

		if matches = htHashHeadRE.FindStringSubmatch(line); nil != matches {
			hash = strings.ToLower(strings.TrimSpace(matches[1]))
		} else if i64, err = strconv.ParseInt(line, 16, 64); nil == err {
			hm.insert(hash, i64)
		}
	}
	if err = scanner.Err(); nil != err {
		return se.New(err, 1)
	}

	return nil
} // loadText()

// `removeID()` deletes all `#hashtags` and `@mentions` associated with `aID`.
//
// Parameters:
//   - `aID`: The object to remove from all references list.
//
// Returns:
//   - `bool`: `true` if `aID` was removed, or `false` otherwise.
func (hm *tHashMap) removeID(aID int64) bool {
	if 0 == len(*hm) {
		return false
	}

	var (
		hash   string
		sl     *tSourceList
		result bool
	)
	for hash, sl = range *hm {
		// remove the ID from every list
		if sl.remove(aID) {
			if 0 == len(*sl) {
				delete(*hm, hash)
			}
			result = true
		}
	}

	return result
} // removeID()

// `removeHM()` deletes `aID` from the list of `aName`.
//
// Parameters:
//   - `aDelim`: The start character of words to use (either '@' or '#').
//   - `aName`: The hash/mention to lookup.
//   - `aID`: The referenced object to removeHM from the list.
//
// Returns:
//   - `bool`: `true` if `aID` was removed, or `false` otherwise.
func (hm *tHashMap) removeHM(aDelim byte, aName string, aID int64) bool {
	if aName = strings.ToLower(strings.TrimSpace((aName))); "" == aName {
		return false
	}

	if aName[0] != aDelim {
		aName = string(aDelim) + aName
	}

	result := false
	if sl, ok := (*hm)[aName]; ok {
		if ok = sl.remove(aID); ok {
			if 0 == len(*sl) {
				delete(*hm, aName)
			}
			result = true
		}
	}

	return result
} // removeHM()

// `renameID()` replaces all occurrences of `aOldID` by `aNewID`.
//
// Parameters:
//   - `aOldID`: The ID to be replaced in this list.
//   - `aNewID`: The replacement ID in this list.
//
// Returns:
//   - `bool`: `true` if the the renaming was successful, or `false` otherwise.
func (hm *tHashMap) renameID(aOldID, aNewID int64) bool {
	if (0 == len(*hm)) || (aOldID == aNewID) {
		return false
	}

	var (
		sl         *tSourceList
		ok, result bool
	)
	for _, sl = range *hm {
		if ok = sl.rename(aOldID, aNewID); ok {
			result = true
		}
	}

	return result
} // renameID()

// `sort()` ensures that the hash map is sorted, which can improve the
// performance of certain operations on the hash map, such as searching
// for a specific key.
// Additionally, sorting the keys can make the hash map easier to read
// and understand, as it presents the keys in a consistent order.
//
// Returns:
//   - `*tHashMap`: The sorted hash map.
func (hm *tHashMap) sort() *tHashMap {
	hLen := len(*hm)
	if 0 == hLen {
		return hm
	}

	var (
		key string
		sl  *tSourceList
	)

	keys := hm.keys()
	// Create a new map to store sorted key-value pairs
	sortedMap := make(tHashMap, len(keys))
	// Iterate through sorted keys and create a new sorted map
	for _, key = range keys {
		sl = (*hm)[key]
		sortedMap[key] = sl.sort()
	}
	*hm = sortedMap

	return hm
} // sort()

// `store()` writes the whole hash/mention list to `aFilename`.
//
// Parameters:
//   - `aFileName`: Name of the file to use for storing the current hash map.
//
// Returns:
//   - `int`: Number of bytes written to storage.
//   - `error`: A possible I/O error.
func (hm *tHashMap) store(aFilename string) (int, error) {
	if aFilename = strings.TrimSpace(aFilename); "" == aFilename {
		return 0, se.New(errors.New("empty filename"), 1)
	}

	file, err := os.OpenFile(aFilename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0660) //#nosec G302 #nosec G304
	if nil != err {
		return 0, se.New(err, 3)
	}
	defer file.Close()

	if !UseBinaryStorage {
		// use plain text storage
		return file.Write([]byte(hm.String()))
	}

	encoder := gob.NewEncoder(file)
	if err = encoder.Encode(hm); nil != err {
		return 0, se.New(err, 1)
	}
	size, err := file.Seek(0, io.SeekEnd)
	if nil != err {
		return 0, se.New(err, 2)
	}

	return int(size), nil
} // store()

// `String()` is used to generate a footprint of the hash map.
//
// It is also used to generate a list of `#hashtags` and `@mentions`
// mostly for debugging purposes.
//
// Returns:
//   - `string`: The string representation of this hash map.
func (hm *tHashMap) String() string {
	if 0 == len(*hm) {
		return ""
	}

	// Pre-allocate buffer to avoid multiple allocations
	var buf bytes.Buffer
	buf.Grow(len(*hm) * 64) // Estimate average size

	var (
		hash string
		sl   *tSourceList
	)
	keys := hm.keys()
	// Iterate through sorted keys and create a new sorted string
	for _, hash = range keys {
		sl = (*hm)[hash]
		buf.WriteString(fmt.Sprintf("[%s]\n%s", hash, sl.String()))
	}

	return buf.String()
} // String()

/* EoF */
