/*
Copyright © 2019, 2024  M.Watermann, 10247 Berlin, Germany

		All rights reserved
	EMail : <support@mwat.de>
*/
package hashtags

import (
	"fmt"
	"slices"
	"sort"
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

	// Find the index of the old value
	result := sort.Search(sLen, func(i int) bool {
		return sl[i] >= aID
	})

	if (result < sLen) && (sl[result] == aID) {
		return result
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

	// find the insertion index using binary search
	idx := sort.Search(sLen, func(i int) bool {
		return (*sl)[i] >= aID
	})

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
		return true
	}

	return false
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
		// sort.SliceStable(*sl, func(i, j int) bool {
		// 	return ((*sl)[i] < (*sl)[j]) // ascending
		// })
		slices.Sort(*sl) // ascending
	}

	return sl
} // sort()

// `String()` implements the `fmt.Stringer` interface.
//
// The method returns the list as a linefeed separated string.
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
