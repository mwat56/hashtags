/*
Copyright © 2019, 2024  M.Watermann, 10247 Berlin, Germany

		All rights reserved
	EMail : <support@mwat.de>
*/
package hashtags

import (
	"fmt"
	"slices"
	"strings"
)

//lint:file-ignore ST1017 - I prefer Yoda conditions

type (
	// `tSourceList` is storing the IDs using a certain #hashtag/@mention.
	tSourceList []uint64
)

// --------------------------------------------------------------------------
// constructor function

// `newSourceList()` creates and returns a new instance of `tSourceList`.
//
// The initial capacity of the list is set to 32 to optimise memory usage.
//
// Returns:
// - `*tSourceList`: A pointer to the newly created instance.
func newSourceList() *tSourceList {
	sl := make(tSourceList, 0, 32)

	return &sl
} // NewSourceList()

// -------------------------------------------------------------------------
// methods of tSourceList

// `clear()` removes all entries in this list.
func (sl *tSourceList) clear() *tSourceList {
	if nil != sl {
		(*sl) = (*sl)[:0]
	}

	return sl
} // clear()

// `equals()` returns whether the current source list is equal
// to the provided source list.
//
// Parameters:
// - `aList`: The source list to compare with.
//
// Returns:
// - `bool`: Whether the source lists are equal.
func (sl tSourceList) equals(aList tSourceList) bool {
	return slices.Equal(sl, aList)
} // equals()

// `findIndex()` returns the list index of `aID`.
//
// Parameters:
// - `aID` is the list element to look up.
//
// Returns:
// - `int`: The index of `aID` in the list.
func (sl tSourceList) findIndex(aID uint64) int {
	sLen := len(sl)
	if 0 == sLen { // empty list
		return -1
	}

	// Find the index of the given ID:
	idx, ok := slices.BinarySearch(sl, aID)
	if !ok {
		return -1
	}

	if (idx < sLen) && (sl[idx] == aID) {
		return idx
	}

	return -1 // aID not found
} // findIndex()

// `insert()` adds `aID` to the list while keeping the list sorted.
//
// Parameters:
// - `aID` the source ID to insert to the list.
//
// Returns:
// - `bool`: `true` if `aID` was inserted, or `false` otherwise.
func (sl *tSourceList) insert(aID uint64) bool {
	if nil == sl {
		return false
	}
	sLen := len(*sl)
	if 0 == sLen { // empty list
		*sl = append(*sl, aID)
		return true
	}

	// Find the index of the given ID:
	idx, _ := slices.BinarySearch(*sl, aID)

	if sLen == idx { // key not found
		*sl = append(*sl, aID) // add new ID
		return true
	}
	if (*sl)[idx] != aID {
		*sl = append(*sl, 0) // make room to insert new ID
		copy((*sl)[idx+1:], (*sl)[idx:])
		(*sl)[idx] = aID
		return true
	}

	return false
} // insert()

// `remove()` deletes the list entry of `aID`.
//
// NOTE: The method's result is an change indicator.
//
// Parameters:
// - `aID`: The ID to look up and delete.
//
// Returns:
// - `bool`: `true` if `aID` was removed, or `false` otherwise.
func (sl *tSourceList) remove(aID uint64) bool {
	sLen := len(*sl)
	if 0 == sLen { // empty list
		return false
	}

	// Find the index of the given ID:
	idx, _ := slices.BinarySearch(*sl, aID)
	if (idx < sLen) && ((*sl)[idx] == aID) {
		// `aID` found at index `idx`
		if 0 == idx {
			if 1 == sLen { // the only element
				*sl = *newSourceList()
			} else { // a longer list
				*sl = (*sl)[1:] // remove the first element
			}
		} else if (sLen - 1) == idx { // remove the last element
			*sl = (*sl)[:idx]
		} else { // remove element in the middle
			*sl = append((*sl)[:idx], (*sl)[idx+1:]...)
		}
		return true
	}

	return false
} // remove()

// `rename()` replaces all occurrences of `aOldID` by `aNewID`.
//
// If `aOldID` equals `aNewID`, or aOldID` doesn't exist then they are
// silently ignored (i.e. this method does nothing), returning `false`.
//
// This method is intended for rare cases when the ID of a document
// gets changed.
//
// Parameters:
// - `aOldID`: ID to be replaced in this list.
// - `aNewID`: The replacement ID in this list.
//
// Returns:
// - `bool`: `true` if the the renaming was successful, or `false` otherwise.
func (sl *tSourceList) rename(aOldID, aNewID uint64) bool {
	if nil == sl {
		return false
	}
	if (0 == len(*sl)) || (aOldID == aNewID) {
		return false
	}

	idx := sl.findIndex(aOldID)
	if 0 > idx { // ID not found
		return sl.insert(aNewID)
	}

	if !sl.insert(aNewID) {
		// This should only happen it there's an OOM problem.
		// Hence we just replace the aOldID by aNewID and sort
		// the list again.
		if (*sl)[idx] != aNewID {
			(*sl)[idx] = aNewID
			sl.sort()
			return true
		}
		return false
	}

	return sl.remove(aOldID)
} // rename()

// `sort()` sorts the list in ascending order.
//
// This method uses the `slices.Sort` function from the standard
// library to sort the list.
//
// Returns:
// - `*tSourceList`: The sorted `tSourceList` instance.
func (sl *tSourceList) sort() *tSourceList {
	if nil != sl {
		slices.Sort(*sl) // ascending
	}

	return sl
} // sort()

// `String()` implements the `fmt.Stringer` interface.
//
// The method returns the list as a linefeed separated string.
// All IDs are represented as strings of 16 hexadecimal characters.
//
// Returns:
// - `string`: The list's contents as a string.
func (sl tSourceList) String() string {
	var result string

	for _, id := range sl {
		strID := fmt.Sprintf("%x\n", id)
		if 17 > len(strID) { // 16 hex chars + LF
			strID = strings.Repeat("0", 17-len(strID)) + strID
		}
		result += strID
	}

	return result
} // String()

/* EoF */
