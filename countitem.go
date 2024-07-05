/*
Copyright © 2023, 2024  M.Watermann, 10247 Berlin, Germany

		All rights reserved
	EMail : <support@mwat.de>
*/
package hashtags

//lint:file-ignore ST1017 - I prefer Yoda conditions

type (
	// `TCountItem` holds a #hashtag/@mention and its number of occurrences.
	TCountItem = struct {
		// number of IDs for this #hashtag/@mention
		Count int
		// name of #hashtag/@mention
		Tag string
	}
)

/* * /
// does not compile: invalid receiver type struct{Count int; Tag string}
func (ci TCountItem) compareTo(aItem TCountItem) int {
	at, bt := ci.Tag, aItem.Tag
	// at, bt := (*ci).Tag, (*aItem).Tag

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
		return 0
	}
} // compareTo()
/* */

/* * /
// does not compile: invalid receiver type struct{Count int; Tag string}
func (ci TCountItem) Less(aItem TCountItem) bool {
	at, bt := ci.Tag, aItem.Tag
	// at, bt := (*ci).Tag, (*aItem).Tag
	if MarkHash == at[0] || MarkMention == at[0] {
		at = at[1:]
	}
	if MarkHash == bt[0] || MarkMention == bt[0] {
		bt = bt[1:]
	}

	return at < bt
} // Less()
/* */

/* EoF */
