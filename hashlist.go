/*
Copyright © 2019, 2025  M.Watermann, 10247 Berlin, Germany

	    All rights reserved
	EMail : <support@mwat.de>
*/
package hashtags

import (
	"bytes"
	"regexp"
	"strings"
)

//lint:file-ignore ST1017 - I prefer Yoda conditions

type (
	// `tHashList` is a list of `#hashtags` and `@mentions`
	// pointing to their respective sources (i.e. IDs).
	tHashList struct {
		hm tHashMap // the actual map list of sources/IDs
	}
)

// --------------------------------------------------------------------------
// constructor function

// `newHashList()` returns a new `tHashList` instance after loading
// the given file.
//
// If the hash file doesn't exist that is not considered an error.
//
// Parameters:
//   - `aFilename`: The name of the file to use for loading and storing.
//
// Returns:
//   - `*tHashList`: The new `tHashList` instance.
//   - `error`: If there is an error, it will be from loading `aFilename`.
func newHashList(aFilename string) (*tHashList, error) {
	result := &tHashList{
		hm: make(tHashMap, 64),
	}

	if aFilename = strings.TrimSpace(aFilename); "" == aFilename {
		return result, nil
	}

	return result.load(aFilename)
} // newHashList()

// -------------------------------------------------------------------------
// methods of `tHashList`

// `checksum()` returns the list's CRC32 checksum.
//
// This method can be used to get a kind of 'footprint' of the current
// contents of the handled data.
//
// Returns:
//   - `uint32`: The computed checksum.
func (hl *tHashList) checksum() uint32 {
	return hl.hm.checksum()
} // checksum()

// `clear()` empties the internal data structures:
// all `#hashtags` and `@mentions` are deleted.
//
// Returns:
//   - `*tHashList`: This cleared list.
func (hl *tHashList) clear() *tHashList {
	hl.hm.clear()

	return hl
} // clear()

// `countedList()` returns a list of `#hashtags` and `@mentions` with
// their respective count of associated IDs.
//
// Returns:
//   - `TCountList`: A list of `#hashtags` and `@mentions` with their respective counts of associated IDs.
func (hl *tHashList) countedList() TCountList {
	return hl.hm.countedList()
} // countedList()

// `equals()` compares the current list with another list.
//
// Parameters:
//   - `aList`: The list to compare with.
//
// Returns:
//   - `bool`: `true` if the lists are identical, `false` otherwise.
func (hl *tHashList) equals(aList *tHashList) bool {
	if nil == aList {
		return false
	}

	if len((*hl).hm) != len((*aList).hm) {
		return false
	}

	return hl.hm.equals(aList.hm)
} // equals()

// `hashCount()` counts the number of hashtags in the list.
//
// Returns:
//   - `int`: The number of hashes in the list.
func (hl *tHashList) hashCount() int {
	return hl.hm.count(MarkHash)
} // hashCount()

// `hashLen()` returns the number of IDs stored for `aHash`.
//
// If `aHash` is empty it is silently ignored (i.e. this method
// does nothing) returning `-1`.
//
// Parameters:
//   - `aHash`: The list key to lookup.
//
// Returns:
//   - `int`: The number of `aHash` in the list.
func (hl *tHashList) hashLen(aHash string) int {
	return hl.hm.idxLen(MarkHash, aHash)
} // hashLen()

// `hashList()` returns a list of IDs associated with `aHash`.
//
// If `aHash` is empty it is silently ignored (i.e. this method
// does nothing), returning an empty slice.
//
// Parameters:
//   - `aHash`: The hash to lookup.
//
// Returns:
//   - `[]int64`: The number of references of `aHash`.
func (hl *tHashList) hashList(aHash string) []int64 {
	return hl.hm.list(MarkHash, aHash)
} // hashList()

// `idList()` returns a list of `#hashtags` and `@mentions` associated
// with `aID`.
//
// Parameters:
//   - `aID`: The referenced object to lookup.
//
// Returns:
//   - `[]string`: The list of `#hashtags` and `@mentions` associated with `aID`.
func (hl *tHashList) idList(aID int64) []string {
	if 0 == len(hl.hm) {
		return []string{}
	}

	return hl.hm.idList(aID)
} // idList()

// `insert()` adds `aID` to the sources list associated with `aName`.
//
// If either `aName` or `aID` are empty they are silently ignored
// (i.e. this method does nothing) returning the current list.
//
// Parameters:
//   - `aDelim`: The start character of words to use (i.e. either '@' or '#').
//   - `aName`: The `#hashtag` or `@mention` to lookup.
//   - `aID`: The referencing object to be added to the list.
//
// Returns:
//   - `bool`: `true` if `aID` was added, or `false` otherwise.
func (hl *tHashList) insert(aDelim byte, aName string, aID int64) bool {
	// prepare for case-insensitive search:
	if aName = strings.ToLower(strings.TrimSpace(aName)); "" == aName {
		return false
	}

	if aName[0] != aDelim {
		aName = string(aDelim) + aName
	}

	return hl.hm.insert(aName, aID)
} // insert()

// `len()` returns the current length of the list i.e. how many
// `#hashtags` and `@mentions` are currently stored in the list.
//
// Returns:
//   - `int`: The length of all `#hashtags` and `@mentions` list.
func (hl *tHashList) len() int {
	return len(hl.hm)
} // len()

// `lenTotal()` returns the length of all `#hashtags` and `@mentions`
// lists stored in the hash list.
//
// Returns:
//   - `int`: The total length of all `#hashtags` and `@mentions` lists.
func (hl *tHashList) lenTotal() (rLen int) {
	if rLen = len(hl.hm); 0 == rLen {
		return
	}
	var sl *tSourceList

	for _, sl = range hl.hm {
		rLen += len(*sl)
	}

	return
} // lenTotal()

// `load()` reads the configured file returning the data structure
// read from the file and a possible error condition.
//
// If the hash file doesn't exist that is not considered an error.
//
// Parameters:
//   - `aFilename`: The name of the file to load.
//
// Returns:
//   - `*tHashList`: The loaded list.
//   - `error`: If there is an error, it will be from loading `aFilename`.
func (hl *tHashList) load(aFilename string) (*tHashList, error) {
	_, err := hl.hm.load(aFilename) // already wrapped

	return hl, err
} // load()

// `mentionCount()` returns the number of mentions in the list.
//
// Returns:
//   - `int`: The number of mentions in the list.
func (hl *tHashList) mentionCount() int {
	return hl.hm.count(MarkMention)
} // mentionCount()

// `mentionLen()` returns the number of IDs stored for `aMention`.
//
// If `aMention` is empty it is silently ignored (i.e. this method
// does nothing) returning `-1`.
//
// Parameters:
//   - `aMention`: Identifies the ID list to lookup.
//
// Returns:
//   - `int`: The number of `aMention` in the list.
func (hl *tHashList) mentionLen(aMention string) int {
	return hl.hm.idxLen(MarkMention, aMention)
} // mentionLen()

// `mentionList()` returns a list of IDs associated with `aMention`.
//
// If `aMention` is empty it is silently ignored (i.e. this method
// does nothing), returning an empty slice.
//
// Parameters:
//   - `aMention`: The mention to lookup.
//
// Returns:
//   - `[]int64`: The number of references of `aMention`.
func (hl *tHashList) mentionList(aMention string) []int64 {
	return hl.hm.list(MarkMention, aMention)
} // mentionList()

var (
	// RegEx to identify a numeric HTML entity.
	htEntityRE = regexp.MustCompile(`#[0-9]+;`)

	// match: [#Hashtag|@Mention]
	htHashHeadRE = regexp.MustCompile(`^\[\s*([#@][^\]]*?)\s*\]$`)
	//                                        11111111111

	// match: #hashtag|@mention
	htHashMentionRE = regexp.MustCompile(
		`(?ims)(?:^|\s|[^\p{L}\d_])?([@#][\p{L}’'\d_§-]+)(?:[^\p{L}\d_]|$)`)
	//	                             1111111111111111111  222222222222222

	// RegEx to match texts like `#----`.
	htHyphenRE = regexp.MustCompile(`#[^-]*--`)
)

// `HashMentionRE()` returns a compiled regular expression used to
// identify `#hashtags` and `@mentions` in a text.
//
// This regular expression matches strings that start with either '@'
// or '#' followed by any number of characters that are not whitespace.
//
// Returns:
//   - `*regexp.Regexp`: A pointer to the compiled regular expression.
func HashMentionRE() *regexp.Regexp {
	return htHashMentionRE
} // HashMentionRE()

// `parseID()` checks whether `aText` contains strings starting with
// `[@|#]` and - if found - adds them to the respective lists with `aID`.
//
// If `aText` is empty it is silently ignored (i.e. this method
// does nothing), returning `false`.
//
// Parameters:
//   - `aID`: The ID to add to the list of hashes/mention.
//   - `aText`: The text to parse for hashtags and mentions.
//
// Returns:
//   - `bool`: `true` if `aID` was updated from `aText`, or `false` otherwise.
func (hl *tHashList) parseID(aID int64, aText []byte) bool {
	if aText = bytes.TrimSpace(aText); 0 == len(aText) {
		return false
	}

	matches := htHashMentionRE.FindAllSubmatch(aText, -1)
	if (nil == matches) || (0 == len(matches)) {
		return false
	}

	var (
		hash, match0 string
		sub          [][]byte
		result       bool
	)
	for _, sub = range matches {
		match0 = string(sub[0])
		hash = string(sub[1])

		if '_' == hash[len(hash)-1] {
			// '_' can be both, part of the hashtag and italic
			// markup so we must remove it if it's at the end:
			hash = hash[:len(hash)-1]
		}
		if MarkHash == hash[0] {
			// `match0` is the match including prefix and postfix
			switch match0[len(match0)-1] {
			case '"':
				// Double quote following a possible hashtag:
				// most probably an URL#fragment, so check
				// whether it's a quoted string:
				if '"' != match0[0] {
					continue // URL#fragment
				}

			case ')':
				// This is a tricky one: It can either be a
				// normal right round bracket or the end of
				// a Markdown link. Here we assume that it's
				// the latter one and ignore this match:
				continue

			case '-':
				// A hyphen at the end of a hashtag:
				// that's not part of an acceptable tag.
				continue

			case ';':
				if htEntityRE.MatchString(match0) {
					// leave HTML entities as is
					continue
				}
			}
			if htHyphenRE.MatchString(hash) {
				continue
			}
		} else if MarkMention == hash[0] {
			if '.' == match0[len(match0)-1] {
				// we assume that it's an email address
				continue
			}
		}
		if hl.insert(hash[0], hash, aID) {
			result = true
		}
	}

	return result
} // parseID()

// `removeHM()` deletes `aID` from the list of `aName`.
//
// Parameters:
//   - `aDelim`: The start character of words to use (i.e. either '@' or '#').
//   - `aName`: The '#hashtag'/'@mention' to lookup for `aID`.
//   - `aID`: The source to remove from the list.
//
// Returns:
//   - `bool`: `true` if `aName` was removed, or `false` otherwise.
func (hl *tHashList) removeHM(aDelim byte, aName string, aID int64) bool {
	if aName = strings.ToLower(strings.TrimSpace(aName)); "" == aName {
		return false
	}

	return hl.hm.removeHM(aDelim, aName, aID)
} // removeHM()

// `removeID()` deletes all `#hashtags` and `@mentions` associated with `aID`.
//
// Parameters:
//   - `aID`: The object to remove from all references list.
//
// Returns:
//   - `bool`: `true` if `aID` was removed, or `false` otherwise.
func (hl *tHashList) removeID(aID int64) bool {
	return hl.hm.removeID(aID)
} // removeID()

// `renameID()` replaces all occurrences of `aOldID` by `aNewID`.
//
// If `aOldID` equals `aNewID` they are silently ignored (i.e. this
// method does nothing), returning `false`.
//
// This method is intended for rare cases when the ID of a document
// needs to get changed.
//
// Parameters:
//   - `aOldID`: The ID to be replaced in all lists.
//   - `aNewID`: The replacement in all lists.
//
// Returns:
//   - `bool`: `true` if the the renaming was successful, or `false` otherwise.
func (hl *tHashList) renameID(aOldID, aNewID int64) bool {
	if (aOldID == aNewID) || (0 == len(hl.hm)) {
		return false
	}

	var (
		result bool
		sl     *tSourceList
	)
	for _, sl = range hl.hm {
		if sl.rename(aOldID, aNewID) {
			result = true
		}
	}

	return result
} // renameID()

// `store()` writes the whole list to the configured file
// returning the number of bytes written and a possible error.
//
// Parameters:
//   - `aFilename`: The name of the file to write.
//
// Returns:
//   - `int`: Number of bytes written to storage.
//   - `error`: A possible storage error, or `nil` in case of success.
func (hl *tHashList) store(aFilename string) (int, error) {
	return hl.hm.store(aFilename)
} // store()

// `String()` returns the whole list as a linefeed separated string.
//
// Returns:
//   - `string`: The string representation of this hash list.
func (hl *tHashList) String() string {
	return hl.hm.String()
} // String()

// `updateID()` checks `aText` removing all `#hashtags` and `@mentions`
// no longer present and adds `#hashtags` and `@mentions` new in `aText`.
//
// Parameters:
//   - `aID`: The ID to update.
//   - `aText`: The text to use.
//
// Returns:
//   - `bool`: `true` if `aID` was updated, or `false` otherwise.
func (hl *tHashList) updateID(aID int64, aText []byte) bool {
	if aText = bytes.TrimSpace(aText); 0 == len(aText) {
		return false
	}

	rr := hl.removeID(aID)
	rp := hl.parseID(aID, aText)

	return (rr || rp)
} // updateID()

/* EoF */
