/*
Copyright © 2023, 2024  M.Watermann, 10247 Berlin, Germany

		All rights reserved
	EMail : <support@mwat.de>
*/
package hashtags

//lint:file-ignore ST1017 - I prefer Yoda conditions

type (
	// `TCountItem` holds a #hashtag/@mention and its number of occurrences.
	TCountItem struct {
		Count int    // number of IDs for this #hashtag/@mention
		Tag   string // name of #hashtag/@mention
	}
)

// -------------------------------------------------------------------------
// methods of TCountItem

// `Compare()` compares two `TCountItem` instances based on their tags
// and counts.
//
// Parameters:
//   - aItem: The other TCountItem instance to compare with.
//
// Returns:
//   - `-1` if the current instance is less than `aItem`.
//   - ` 0` if the current instance is equal to `aItem`.
//   - `+1` if the current instance is greater than `aItem`.
func (ci TCountItem) Compare(aItem TCountItem) int {
	at, bt := ci.Tag+"Z", aItem.Tag+"Z" // avoid empty strings

	if MarkHash == at[0] || MarkMention == at[0] {
		at = at[1:]
	}
	if MarkHash == bt[0] || MarkMention == bt[0] {
		bt = bt[1:]
	}

	switch true {
	case at < bt:
		return -1
	case at > bt:
		return +1
	default:
		switch true {
		case ci.Count < aItem.Count:
			return -1
		case ci.Count > aItem.Count:
			return 1
		}
	}

	return 0
} // Compare()

// `Equal()` checks whether the current `TCountItem` is equal to `aItem`.
//
// Parameters:
//   - `aItem`: The other item to compare with.
//
// Returns:
//   - `bool`: Whether the two items are equal.
func (ci TCountItem) Equal(aItem TCountItem) bool {
	return (0 == ci.Compare(aItem))
} // Equal()

// `Less()` checks whether this `TCountItem` is less than `aItem`.
//
// Parameters:
//   - `aItem`: The other `TCountItem` instance to compare with.
//
// Returns:
//   - `bool`: Whether the current instance is less than the other item.
func (ci TCountItem) Less(aItem TCountItem) bool {
	return (-1 == ci.Compare(aItem))
} // Less()

/* EoF */
