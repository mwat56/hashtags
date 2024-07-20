/*
Copyright © 2023, 2024  M.Watermann, 10247 Berlin, Germany

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
	// A list of `TCountItems
	TCountList []TCountItem
)

// -------------------------------------------------------------------------
// methods of TCountList

// `Compare()` compares the current list with another list.
//
// Parameters:
//   - `aList`: The list to compare with.
//
// Returns:
//   - `-1` if the current instance is less than `aList`.
//   - ` 0` if the current instance is equal to `aList`.
//   - `+1` if the current instance is greater than `aList`.
func (cl TCountList) Compare(aList TCountList) int {
	sLen, aLen := len(cl), len(aList)
	if sLen < aLen {
		return -1
	}
	if sLen > aLen {
		return 1
	}

	for idx, ci := range cl {
		oci := aList[idx]
		if cmp := ci.Compare(oci); 0 != cmp {
			return cmp
		}
	}

	return 0
} // Compare()

// `Equal()` compares the current list with another list.
//
// Parameters:
//   - `aList`: The list to compare with.
//
// Returns:
//   - `bool`: True if the lists are identical, false otherwise.
func (cl TCountList) Equal(aList TCountList) bool {
	// if len(cl) != len(aList) {
	// 	return false
	// }

	// for idx, ci := range cl {
	// 	oci := aList[idx]
	// 	if ci.Tag != oci.Tag {
	// 		return false
	// 	}
	// 	if ci.Count != oci.Count {
	// 		return false
	// 	}
	// }

	return (0 == cl.Compare(aList))
} // Equal()

// `Insert()` appends `aItem` to the list.
//
// Parameters:
//   - `aItem`: The source ID to insert into the list.
func (cl *TCountList) Insert(aItem TCountItem) *TCountList {
	sLen := len(*cl)
	if 0 == sLen { // empty list
		*cl = append(*cl, aItem)
		return cl
	}

	// find the insertion index using binary search
	idx := sort.Search(sLen, func(i int) bool {
		return (*cl)[i].Tag >= aItem.Tag
	})

	if sLen == idx { // item not found
		// add new ID
		*cl = append(*cl, aItem)
		return cl
	}

	if (*cl)[idx] != aItem {
		// make room to insert new item
		*cl = append(*cl, TCountItem{})
		copy((*cl)[idx+1:], (*cl)[idx:])
		(*cl)[idx] = aItem
	}

	return cl
} // insert()

func (cl TCountList) Len() int {
	return len(cl)
} // Len()

// `Less()` checks whether this `TCountItem` is less than `aList`.
//
// Parameters:
//   - `aList`: The other `TCountList` instance to compare with.
//
// Returns:
//   - `bool`: Whether the current instance is less than the other item.
func (cl TCountList) Less(aList TCountList) bool {
	return (-1 == cl.Compare(aList))
} // Less()

// `sort()` sorts the list in ascending order based on the count of
// occurrences and the tag name.
//
// The sorting is stable, meaning that equal elements preserve their
// original order.
//
// Returns:
//   - `*TCountList`: A pointer to the sorted list.
func (cl *TCountList) sort() *TCountList {
	if 0 == len(*cl) {
		return cl
	}
	// `cmpF()` is a comparison function that compares two
	// `TCountItem` instances.
	cmpF := func(a, b TCountItem) int {
		return a.Compare(b)
	}
	slices.SortStableFunc(*cl, cmpF)

	return cl
} // sort()

// `String()` returns the list as a linefeed separated string.
//
// (Implements `Stringer` interface)
// Returns:
//   - `string`: The string representation of this list.
func (cl TCountList) String() (rStr string) {
	for _, tc := range cl {
		rStr += fmt.Sprintf("%s: %d\n", tc.Tag, tc.Count)
	}

	return
} // String()

// `Swap()` swaps the elements at the specified indices in the list.
//
// If the list is empty, or the old and new indices are the same, or
// the function, or either of the indices is out of bounds, the
// function returns the list unchanged.
//
// Parameters:
//   - `aOldIdx`: The index of the first element to swap.
//   - `aNewIdx`: The index of the second element to swap.
//
// Returns:
//   - `*TCountList`: A pointer to the list with the swapped elements.
func (cl *TCountList) Swap(aOldIdx, aNewIdx int) *TCountList {
	sLen := len(*cl)
	if 0 == sLen || aOldIdx == aNewIdx ||
		0 > aOldIdx || aOldIdx >= sLen ||
		0 > aNewIdx || aNewIdx >= sLen {
		return cl
	}
	(*cl)[aOldIdx], (*cl)[aNewIdx] = (*cl)[aNewIdx], (*cl)[aOldIdx]

	return cl
} // Swap()

/* EoF */
