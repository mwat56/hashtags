/*
Copyright © 2019, 2024  M.Watermann, 10247 Berlin, Germany

		All rights reserved
	EMail : <support@mwat.de>
*/
package hashtags

import (
	"regexp"
	"strings"
)

//lint:file-ignore ST1017 - I prefer Yoda conditions

type (
		// `tHashList` is a list of `#hashtags` and `@mentions`
	// pointing to sources (i.e. IDs).
	tHashList struct {
		hm tHashMap // the actual map list of sources/IDs
	}
)

// --------------------------------------------------------------------------
// constructor function

// `newHashList()` returns a new `THashList` instance after reading
// the given file.
//
// If the hash file doesn't exist that is not considered an error.
// If there is an error, it will be of type *PathError.
//
// Parameters:
// - `aFilename` is the name of the file to use for reading and storing.
func newHashList(aFilename string) (*tHashList, error) {
	result := &tHashList{
		hm: make(tHashMap, 64),
	}
	if 0 == len(aFilename) {
		return result, nil
	}

	return result.Load(aFilename)
} // newHashList()

// -------------------------------------------------------------------------
// methods of THashList

// `add()` appends `aID` to the list associated with `aMapIdx`.
//
// If either `aName` or `aID` are empty they are silently ignored
// (i.e. this method does nothing) returning the current list.
//
// Parameters:
// - `aDelim`: The start character of words to use (i.e. either '@' or '#').
// - `aName`: The hashtag/mention to lookup.
// - `aID`: The referencing object to be added to the hash list.
//
// Returns:
// - `*tHashList`: This hash list.
func (hl *tHashList) add(aDelim byte, aName string, aID uint64) *tHashList {
	// prepare for case-insensitive search:
	aName = strings.ToLower(strings.TrimSpace(aName))
	if 0 == len(aName) {
		return hl
	}

	if aName[0] != aDelim {
		aName = string(aDelim) + aName
	}

	return hl.add0(aName, aID)
} // add()

// `add0()` appends `aID` to the tag list associated with `aMapIdx`.
//
// Parameters:
// - `aName`: The hashtag/mention to lookup.
// - `aID` is to be added to the hash list.
//
// Returns:
// - `*tHashList`: This hash list.
func (hl *tHashList) add0(aName string, aID uint64) *tHashList {
	// the mutex.Lock is done by the callers

	hl.hm.add(aName, aID)

	return hl
} // add0()

// `checksum()` returns the list's CRC32 checksum.
//
// This method can be used to get a kind of 'footprint'.
//
// Returns:
// - `uint32`: The computed checksum.
func (hl *tHashList) checksum() uint32 {
	return hl.hm.checksum()
} // checksum()

func (hl *tHashList) clear() *tHashList {
	hl.hm.clear()

	return hl
} // clear()

// ` compareTo()` compares the current list with another list.
//
// Parameters:
// - `aList`: The list to compare with.
//
// Returns:
// - `bool`: `true` if the lists are identical, `false` otherwise.
func (hl *tHashList) compareTo(aList *tHashList) bool {
	if len((*hl).hm) != len((*aList).hm) {
		return false
	}

	return hl.hm.compareTo(aList.hm)
} // compareTo()

// `HashCount()` counts the number of hashtags in the list.
//
// Returns:
// - `int`: The number of hashes in the list.
func (hl *tHashList) HashCount() int {
	return hl.hm.count(MarkHash)
} // HashCount()

// `HashLen()` returns the number of IDs stored for `aHash`.
//
// Parameters:
// - `aHash` The list key to lookup.
//
// Returns:
// - `int`: The number of `aHash` in the list.
func (hl *tHashList) HashLen(aHash string) int {
	return hl.hm.idxLen(MarkHash, aHash)
} // HashLen()

// `HashList()` returns a list of IDs associated with `aHash`.
//
// Parameters:
// - `aName`: The hash to lookup.
//
// Returns:
// - `[]uint64`: The number of references of `aName`.
func (hl *tHashList) HashList(aHash string) []uint64 {
	return hl.hm.list(MarkHash, aHash)
} // HashList()

// `IDlist()` returns a list of #hashtags and @mentions associated with `aID`.
//
// Parameters:
// - `aID`: The referenced object to lookup.
//
// Returns:
// - `[]string`: The list of #hashtags and @mentions associated with `aID`.
func (hl *tHashList) idList(aID uint64) []string {
	if 0 == len(hl.hm) {
		return nil
	}

	return hl.hm.idList(aID)
} // idList()

// `IDremove()` deletes all #hashtags/@mentions associated with `aID`.
//
// Parameters:
// - `aID`: The object to remove from all references list.
//
// Returns:
// - `*THashList`: The modified hash list.
func (hl *tHashList) IDremove(aID uint64) *tHashList {
	if (nil == hl) || (0 == len(hl.hm)) {
		return hl
	}

	hl.hm.removeID((aID))

	return hl
} // IDremove()

// `IDrename()` replaces all occurrences of `aOldID` by `aNewID`.
//
// This method is intended for rare cases when the ID of a document
// needs to get changed.
//
// Parameters:
// - `aOldID`: The ID to be replaced in all lists.
// - `aNewID`: The replacement in all lists.
//
// Returns:
// - `*THashList`: The modified hash list.
func (hl *tHashList) IDrename(aOldID, aNewID uint64) *tHashList {
	if (aOldID == aNewID) || (0 == len(hl.hm)) {
		return hl
	}
	for _, sl := range hl.hm {
		sl.rename(aOldID, aNewID)
	}

	return hl
} // IDrename()

// `IDupdate()` checks `aText` removing all #hashtags/@mentions no longer
// present and adds #hashtags/@mentions new in `aText`.
//
// Parameters:
// - `aID`: The ID to update.
// - `aText`: The text to use.
//
// Returns:
// - `*THashList`: The current hash list.
func (hl *tHashList) IDupdate(aID uint64, aText []byte) *tHashList {
	if (nil == hl) || (0 == len(aText)) || (0 == len(hl.hm)) {
		return hl
	}

	return hl.IDremove(aID).parseID(aID, aText)
} // IDupdate()

// `Len()` returns the current length of the list i.e. how many #hashtags
// and @mentions are currently stored in the list.
func (hl *tHashList) Len() int {
	// hl.mtx.RLock()
	// defer hl.mtx.RUnlock()

	return len(hl.hm)
} // Len()

// `LenTotal()` returns the length of all #hashtag/@mention lists stored
// in the hash list.
//
// Returns:
// - `int`: The total length of all #hashtag/@mention lists.
func (hl *tHashList) LenTotal() (rLen int) {
	// hl.mtx.RLock()
	// defer hl.mtx.RUnlock()

	rLen = len(hl.hm)
	for _, sl := range hl.hm {
		rLen += len(*sl)
	}

	return
} // LenTotal()

// `List()` returns a list of #hashtags/@mentions with
// their respective count of associated IDs.
//
// Returns:
// - `TCountList`: A list of #hashtags/@mentions with their
// respective counts of associated IDs.
func (hl *tHashList) List() TCountList {
	return hl.hm.countedList()
} // List()

// `Load()` reads the configured file returning the data structure
// read from the file and a possible error condition.
//
// If the hash file doesn't exist that is not considered an error.
//
// Returns:
// - `*THashList`: The updated list.
// - `error`: If there is an error, it will be of type `*PathError`.
func (hl *tHashList) Load(aFilename string) (*tHashList, error) {
	_, err := hl.hm.Load(aFilename)
	return hl, err
} // Load()

// `MentionCount()` returns the number of mentions in the list.
//
// Returns:
// - `int`: The number of mentions in the list.
func (hl *tHashList) MentionCount() int {
	return hl.hm.count(MarkMention)
} // MentionCount()

// `MentionLen()` returns the number of IDs stored for `aMention`.
//
// Parameters:
// - `aMention` identifies the ID list to lookup.
//
// Returns:
// - `int`: The number of `aMention` in the list.
func (hl *tHashList) MentionLen(aMention string) int {
	return hl.hm.idxLen(MarkMention, aMention)
} // MentionLen()

// `MentionList()` returns a list of IDs associated with `aMention`.
//
// Parameters:
// - `aMention`: The mention to lookup.
//
// Returns:
// - `[]uint64`: The number of references of `aName`.
func (hl *tHashList) MentionList(aMention string) []uint64 {
	return hl.hm.list(MarkMention, aMention)
} // MentionList()

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

// `parseID()` checks whether `aText` contains strings starting
// with `[@|#]` and – if found – adds them to the respective list.
//
// Parameters:
// - `aID`: The ID to add to the list of hashes/mention.
// - `aText`: The text to parse for hashtags and mentions.
//
// Returns:
// - `*THashList`: The current hash list.
func (hl *tHashList) parseID(aID uint64, aText []byte) *tHashList {
	if 0 == len(aText) {
		return hl
	}

	matches := htHashMentionRE.FindAllSubmatch(aText, -1)
	if (nil == matches) || (0 == len(matches)) {
		return hl
	}

	for _, sub := range matches {
		match0 := string(sub[0])
		hash := string(sub[1])

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
				// This is a tricky one: it can either be a
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
		}
		hl.add(hash[0], hash, aID)
	}

	return hl
} // parseID()

// `removeHM()` deletes `aID` from the list of `aName`.
//
// Parameters:
// - `aDelim` is the start character of words to use (i.e. either '@' or '#').
// - `aName`: The hash/mention to lookup for `aID`.
// - `aID` is the source to removeHM from the list.
//
// Returns:
// - `*THashList`: The current hash list.
func (hl *tHashList) removeHM(aDelim byte, aName string, aID uint64) *tHashList {
	aName = strings.ToLower(strings.TrimSpace(aName))
	if (0 == len(aName)) || (0 == aID) {
		return hl
	}

	hl.hm.remove(aDelim, aName, aID)

	return hl
} // removeHM()

// `Store()` writes the whole list to the configured file
// returning the number of bytes written and a possible error.
//
// If there is an error, it will be of type `*PathError`.
//
// Parameters:
// - 'aFilename`: The name of the file to write.
//
// Returns:
// - `int`: Number of bytes written to storage.
// - `error`: A possible storage error, or `nil` in case of success.
func (hl *tHashList) Store(aFilename string) (int, error) {
	// hl.mtx.RLock()
	// defer hl.mtx.RUnlock()

	return hl.hm.store(aFilename)
} // Store()

// `String()` returns the whole list as a linefeed separated string.
//
// Returns:
// - `string`: The string representation of this hash list.
func (hl *tHashList) String() string {
	// hl.mtx.RLock()
	// defer hl.mtx.RUnlock()

	return hl.hm.String()
} // String()

// -------------------------------------------------------------------------

type (
	// `TWalkFunc` is used by `Walk()` when visiting an entry
	// in the #hashtag/@mention lists.
	//
	// Parameters:
	// - `aHash`: The hash list index to check.
	// - `aID`: The ID to check.
	//
	// Returns:
	// - `bool`: `true` if the entry was successfully visited,
	// or `false` otherwise
	//
	// see `Walk()`
	TWalkFunc func(aHash string, aID uint64) bool

	// `IHashWalker` is used by `Walker()` when visiting an entry
	// in the #hashtag/@mentions lists.
	IHashWalker interface {
		Walk(aHash string, aID uint64) bool
	}
)

// `Walk()` traverses through all entries in the #hashtag/@mention lists
// calling `aFunc` for each entry.
//
// If `aFunc` returns `false` when called the respective ID
// will be removed from the associated #hashtag/@mention.
//
// Parameters:
// - `aFunc` The function called for each ID in all lists.
// func (hl *tHashList) Walk(aFunc TWalkFunc) {
// 	// hl.mtx.RLock()
// 	// defer hl.mtx.RUnlock()

// 	oldCRC := hl.checksum()
// 	defer func() {
// 		if oldCRC != atomic.LoadUint32(&hl.µChange) {
// 			_, _ = hl.hm.store(hl.fn)
// 		}
// 	}()

// 	changed := hl.hm.walk(aFunc)
// 	if changed {
// 		atomic.StoreUint32(&hl.µChange, 0)
// 	}
// } // Walk()

// `Walker()` traverses through all entries in the hash lists
// calling `aWalker` for each entry.
//
// Parameters:
// - `aWalker` is an object implementing the `IHashWalker` interface.
// func (hl *tHashList) Walker(aWalker IHashWalker) {
// 	hl.Walk(aWalker.Walk)
// } // Walker()

/* EoF */
