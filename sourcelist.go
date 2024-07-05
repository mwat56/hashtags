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

	// `tSourceList` is storing the IDs of a single #hashtag/@mention.
	tSourceList []uint64
)

// --------------------------------------------------------------------------
// constructor function

// `NewSourceList()` creates and returns a new instance of tSourceList.
//
// The initial capacity of the list is set to 32 to optimize memory usage.
//
// Returns:
// - `*tSourceList`: A pointer to the newly created tSourceList instance.
func NewSourceList() *tSourceList {
	sl := make(tSourceList, 0, 32)

	return &sl
} // NewSourceList()

// -------------------------------------------------------------------------
// methods of tSourceList

// `clear()` removes all entries in this list.
func (sl *tSourceList) clear() *tSourceList {
	(*sl) = (*sl)[:0]

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
	for result, id := range sl {
		if id == aID {
			return result
		}
	}

	return -1
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
		*sl = append(*sl, aID)
		return sl
	}

	// find the insertion index using binary search
	idx := sort.Search(sLen, func(i int) bool {
		return (*sl)[i] >= aID
	})

	if sLen == idx { // key not found
		// add new ID
		*sl = append(*sl, aID)
		return sl
	}

	if (*sl)[idx] != aID {
		// make room to insert new ID
		*sl = append(*sl, 0)
		copy((*sl)[idx+1:], (*sl)[idx:])
		(*sl)[idx] = aID
	}

	return sl
} // insert()

// `removeID()` deletes the list entry of `aID`.
//
// Parameters:
// - `aID` is the string to look up.
//
// Returns:
// - `*tSourceList`: the current list.
func (sl *tSourceList) removeID(aID uint64) *tSourceList {
	sLen := len(*sl)
	if 0 == sLen { // empty list
		return sl
	}

	// Find the index of the old value
	idx := sort.Search(sLen, func(i int) bool {
		return (*sl)[i] >= aID
	})

	// If the old value is not found, return the original slice
	if idx == sLen || (*sl)[idx] != aID {
		return sl
	}

	if idx == 0 {
		*sl = (*sl)[1:]
	} else {
		// Remove the old value
		*sl = append((*sl)[:idx], (*sl)[idx+1:]...)
	}

	/*
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
				*sl = (*sl)[:0] // clear the list
			} else {
				*sl = (*sl)[1:]
			}

		case sLen:
			*sl = (*sl)[:sLen]

		default:
			*sl = append((*sl)[:idx], (*sl)[idx+1:]...)
		}
	*/
	return sl
} // removeID()

/*
// `renameID()` replaces all occurrences of `aOldID` by `aNewID`.
//
// This method is intended for rare cases when the ID of a document
// gets changed.
//
// - `aOldID` is to be replaced in this list.
// - `aNewID` is the replacement in this list.
func (sl *tSourceList) renameID(aOldID, aNewID uint64) *tSourceList {

	for idx, id := range *sl {
		if id == aOldID {
			(*sl)[idx] = aNewID
			return sl //.sort()
		}
	}

	return sl
} // renameID()
*/

// `renameID()` replaces all occurrences of `aOldID` by `aNewID`.
//
// This method is intended for rare cases when the ID of a document
// gets changed.
//
// Parameters:
// - `aOldID` is to be replaced in this list.
// - `aNewID` is the replacement in this list.
//
// Returns:
// - `*tSourceList`: the current list.
func (sl *tSourceList) renameID(aOldID, aNewID uint64) *tSourceList {
	if aOldID == aNewID {
		return sl
	}
	return sl.removeID(aOldID).insert(aNewID)
} // renameID()

/* */
// `sort()` returns the sorted list.
func (sl *tSourceList) sort() *tSourceList {
	sort.Slice(*sl, func(i, j int) bool {
		return ((*sl)[i] < (*sl)[j]) // ascending
	})

	return sl
} // sort()
/* */

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
