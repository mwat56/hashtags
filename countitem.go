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

// `compareTo()` checks whether the current `TCountItem` is equal to `aID`.
//
// Parameters:
// - `aItem`: The other item to compare with.
//
// Returns:
// - `bool`: Whether the two items are equal.
func (ci TCountItem) compareTo(aItem TCountItem) bool {
	if 0 == len(ci.Tag) {
		return 0 == len(aItem.Tag)
	}
	at, bt := ci.Tag+"Z", aItem.Tag+"Z" // avoid empty strings

	if MarkHash == at[0] || MarkMention == at[0] {
		at = at[1:]
	}
	if MarkHash == bt[0] || MarkMention == bt[0] {
		bt = bt[1:]
	}
	if at == bt {
		return ci.Count == aItem.Count
	}

	return false
} // compareTo()

// `Less()` checks whether this `TCountItem` is less than `aItem`.
//
// * Parameters:
// - `aItem`: The other `TCountItem` instance to compare with.
//
// * Returns:
// - `bool`: Whether the current instance's tag is less than the
// other instance's tag.
func (ci TCountItem) Less(aItem TCountItem) bool {
	at, bt := ci.Tag+"Z", aItem.Tag+"Z" // avoid empty strings

	if MarkHash == at[0] || MarkMention == at[0] {
		at = at[1:]
	}
	if MarkHash == bt[0] || MarkMention == bt[0] {
		bt = bt[1:]
	}
	if at < bt {
		return ci.Count < aItem.Count
	}

	return false
} // Less()

/* EoF */
