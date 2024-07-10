/*
Copyright © 2019, 2024  M.Watermann, 10247 Berlin, Germany

		All rights reserved
	EMail : <support@mwat.de>
*/
package hashtags

import (
	"fmt"
	"sort"
)

//lint:file-ignore ST1017 - I prefer Yoda conditions

type (
	// `tSourceList` is storing the IDs using a certain #hashtag/@mention.
	tSourceList []uint64
)

// --------------------------------------------------------------------------
// constructor function

// `newSourceList()` creates and returns a new instance of tSourceList.
//
// The initial capacity of the list is set to 32 to optimize memory usage.
//
// Returns:
// - `*tSourceList`: A pointer to the newly created tSourceList instance.
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

// `compareTo()` returns whether the current source list is equal
// to the provided source list.
//
// Parameters:
// - `aList`: The source list to compare with.
//
// Returns:
// - `bool`: Whether the source lists are equal.
func (sl tSourceList) compareTo(aList tSourceList) bool {
	if len(sl) != len(aList) {
		return false
	}

	// Iterate over the source lists and compare each element.
	for idx, id := range sl {
		// If the element in the current source list does not
		// match the corresponding element in the provided source
		// list, the source lists are not equal.
		if id != aList[idx] {
			return false
		}
	}

	// If no mismatches were found, the source lists are equal.
	return true
} // compareTo()

// `indexOf()` returns the list index of `aID`.
//
// Parameters:
// - `aID` is the string to look up.
//
// Returns:
// - `int`: the index of `aID` in the list.
func (sl tSourceList) indexOf(aID uint64) int {
	sLen := len(sl)
	if 0 == sLen { // empty list
		return -1
	}

	// Find the index of the old value
	result := sort.Search(sLen, func(i int) bool {
		return sl[i] >= aID
	})

	if (result < sLen) && (sl[result] == aID) {
		return result
	}

	return -1 // aID not found
} // indexOf()

// `insert()` adds `aID` to the list while keeping the list sorted.
//
// Parameters:
// - `aID` the source ID to insert to the list.
//
// Returns:
// - `*tSourceList`: the current list.
func (sl *tSourceList) insert(aID uint64) *tSourceList {
	sLen := len(*sl)
	if 0 == sLen { // empty list
		if nil != sl {
			*sl = append(*sl, aID)
		}

		return sl
	}

	// find the insertion index using binary search
	idx := sort.Search(sLen, func(i int) bool {
		return (*sl)[i] >= aID
	})

	if sLen == idx { // key not found
		*sl = append(*sl, aID) // add new ID
	} else if (*sl)[idx] != aID {
		*sl = append(*sl, 0) // make room to insert new ID
		copy((*sl)[idx+1:], (*sl)[idx:])
		(*sl)[idx] = aID
	}

	return sl
} // insert()

// `remove()` deletes the list entry of `aID`.
//
// Parameters:
// - `aID`: The ID to look up and delete.
//
// Returns:
// - `*tSourceList`: The modified list.
func (sl *tSourceList) remove(aID uint64) *tSourceList {
	sLen := len(*sl)
	if 0 == sLen { // empty list
		return sl
	}

	// Find the index of the old value
	idx := sort.Search(sLen, func(i int) bool {
		return (*sl)[i] >= aID
	})

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
	} else {
		idx = -1 // DEBUG breakpoint: aID not found
	}

	return sl
} // remove()

// `rename()` replaces all occurrences of `aOldID` by `aNewID`.
//
// If `aOldID` equals `aNewID`, or aOldID` doesn't exist then nothing
// is changed.
//
// This method is intended for rare cases when the ID of a document
// gets changed.
//
// Parameters:
// - `aOldID`: ID to be replaced in this list.
// - `aNewID`: The replacement ID in this list.
//
// Returns:
// - `*tSourceList`: The current list.
func (sl *tSourceList) rename(aOldID, aNewID uint64) *tSourceList {
	if (0 == len(*sl)) || (aOldID == aNewID) {
		return sl
	}

	if 0 > sl.indexOf(aOldID) { // ID not found
		return sl
	}

	return sl.remove(aOldID).insert(aNewID)
} // rename()

// `sort()` sorts the list in ascending order.
//
// This method uses the `sort.SliceStable` function from the standard
// library to sort the list.
// The sorting is stable, meaning that equal elements will retain their
// original order.
//
// Returns:
// - `*tSourceList`: The sorted `tSourceList` instance.
func (sl *tSourceList) sort() *tSourceList {
	if nil != sl {
		sort.SliceStable(*sl, func(i, j int) bool {
			return ((*sl)[i] < (*sl)[j]) // ascending
		})
	}

	return sl
} // sort()

// `String()` returns the list as a linefeed separated string.
//
// (Implements `Stringer` interface.)
//
// Returns:
// - `string`: The list's contents as a string.
func (sl tSourceList) String() string {
	var result string
	for _, id := range sl {
		result += fmt.Sprintf("%d\n", id)
	}

	return result
} // String()

/* EoF */
