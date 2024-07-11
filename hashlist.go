/*
Copyright © 2019, 2024  M.Watermann, 10247 Berlin, Germany

		All rights reserved
	EMail : <support@mwat.de>
*/
package hashtags

import (
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
)

//lint:file-ignore ST1017 - I prefer Yoda conditions

const (
	// `MarkHash` is the first character in asy hash tag.
	MarkHash = byte('#')

	// `MarkMention` is the first character in asy mention tag.
	MarkMention = byte('@')
)

type (
	// Data cache for `CountedList()`
	tCountCache struct {
		µCRC    uint32
		µCounts TCountList
	}

	// `THashList` is a list of `#hashtags` and `@mentions`
	// pointing to sources (i.e. IDs).
	THashList struct {
		fn      string       // the filename to use
		hm      tHashMap     // the actual map list of sources/IDs
		mtx     sync.RWMutex // safeguard against concurrent accesses
		µChange uint32       // internal change flag
		µCC     tCountCache  // cache for `CountedList()`
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
func newHashList(aFilename string) (*THashList, error) {
	result := &THashList{
		fn: aFilename,
		hm: make(tHashMap, 64),
	}
	if 0 == len(aFilename) {
		return result, nil
	}

	return result.Load()
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
func (hl *THashList) add(aDelim byte, aName string, aID uint64) *THashList {
	// the mutex.Lock is done by the callers

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
func (hl *THashList) add0(aName string, aID uint64) *THashList {
	// the mutex.Lock is done by the callers

	hl.hm.add(aName, aID)
	atomic.StoreUint32(&hl.µChange, 0)

	return hl
} // add0()

func (hl *THashList) checksum() uint32 {
	// the mutex.Lock is done by the callers

	if 0 == atomic.LoadUint32(&hl.µChange) {
		atomic.StoreUint32(&hl.µChange, hl.hm.checksum())
	}

	return atomic.LoadUint32(&hl.µChange)
} // checksum()

// `Checksum()` returns the list's CRC32 checksum.
//
// This method can be used to get a kind of 'footprint'.
//
// Returns:
// - `uint32`: The computed checksum.
func (hl *THashList) Checksum() uint32 {
	hl.mtx.RLock()
	defer hl.mtx.RUnlock()

	return hl.checksum()
} // Checksum()

func (hl *THashList) clear() *THashList {
	// the mutex.Lock is done by the callers

	hl.hm.clear()
	atomic.StoreUint32(&hl.µChange, 0)

	return hl
} // clear()

// `Clear()` empties the internal data structures:
// all `#hashtags` and `@mentions` are deleted.
//
// Returns:
// - `*tHashList`: This hash list.
func (hl *THashList) Clear() *THashList {
	hl.mtx.Lock()
	defer hl.mtx.Unlock()

	return hl.clear()
} // Clear()

func (hl *THashList) compare2(aList *THashList) bool {
	// the mutex.Lock is done by the caller

	if len((*hl).hm) != len((*aList).hm) {
		return false
	}

	return hl.hm.compareTo(aList.hm)
} // compare2()

// ` compareTo()` compares the current list with another list.
//
// Parameters:
// - `aList`: The list to compare with.
//
// Returns:
// - `bool`: True if the lists are identical, false otherwise.
func (hl *THashList) compareTo(aList *THashList) bool {
	hl.mtx.Lock()
	defer hl.mtx.Unlock()

	return hl.compare2(aList)
} // compareTo()

// `Filename()` returns the configured filename for reading/storing
// this list's contents.
//
// Returns:
// - `string`: The filename for reading/storing this list.
func (hl *THashList) Filename() string {
	return hl.fn
} // Filename()

// `HashAdd()` appends `aID` to the list of `aHash`.
//
// If `aHash` is an empty string it is silently ignored
// (i.e. this method does nothing).
//
// Parameters:
// - `aHash`: The hash list index to use.
// - `aID`: The object to be added to the hash list.
//
// Returns:
// - `*tHashList`: This hash list.
func (hl *THashList) HashAdd(aHash string, aID uint64) *THashList {
	if 0 == len(aHash) {
		return hl
	}
	hl.mtx.Lock()
	defer hl.mtx.Unlock()

	return hl.add(MarkHash, aHash, aID)
} // HashAdd()

// `HashCount()` counts the number of hashtags in the list.
//
// Returns:
// - `int`: The number of hashes in the list.
func (hl *THashList) HashCount() int {
	hl.mtx.RLock()
	defer hl.mtx.RUnlock()

	return hl.hm.count(MarkHash)
} // HashCount()

// `HashLen()` returns the number of IDs stored for `aHash`.
//
// Parameters:
// - `aHash` The list key to lookup.
//
// Returns:
// - `int`: The number of `aHash` in the list.
func (hl *THashList) HashLen(aHash string) int {
	hl.mtx.RLock()
	defer hl.mtx.RUnlock()

	return hl.hm.idxLen(MarkHash, aHash)
} // HashLen()

// `HashList()` returns a list of IDs associated with `aHash`.
//
// Parameters:
// - `aName`: The hash to lookup.
//
// Returns:
// - `[]uint64`: The number of references of `aName`.
func (hl *THashList) HashList(aHash string) []uint64 {
	hl.mtx.RLock()
	defer hl.mtx.RUnlock()

	return hl.hm.list(MarkHash, aHash)
} // HashList()

// `HashRemove()` deletes `aID` from the list of `aHash`.
//
// Parameters:
// - `aHash`: The hash to lookup.
// - `aID`: The referenced object to remove from the list.
//
// Returns:
// - `*tHashList`: The current hash list.
func (hl *THashList) HashRemove(aHash string, aID uint64) *THashList {
	if 0 == len(aHash) {
		return hl
	}

	hl.mtx.RLock()
	defer hl.mtx.RUnlock()

	hl.removeHM(MarkHash, aHash, aID)

	return hl
} // HashRemove()

func (hl *THashList) idList(aID uint64) []string {
	if 0 == len(hl.hm) {
		return nil
	}

	return hl.hm.idList(aID)
} // ifList()

// `IDlist()` returns a list of #hashtags and @mentions associated with `aID`.
//
// Parameters:
// - `aID`: The referenced object to lookup.
//
// Returns:
// - `[]string`: The list of #hashtags and @mentions associated with `aID`.
func (hl *THashList) IDlist(aID uint64) []string {
	hl.mtx.RLock()
	defer hl.mtx.RUnlock()

	return hl.idList(aID)
} // IDlist()

// `IDparse()` checks whether `aText` contains strings starting with
// `[@|#]` and - if found - adds them to the respective list.
//
// Parameters:
// - `aID`: the ID to add to the list.
// - `aText:` The text to search.
func (hl *THashList) IDparse(aID uint64, aText []byte) *THashList {
	if 0 == len(aText) {
		return hl
	}
	hl.mtx.Lock()
	defer hl.mtx.Unlock()

	oldCRC := hl.checksum()
	defer func() {
		if oldCRC != atomic.LoadUint32(&hl.µChange) {
			hl.hm.store(hl.fn)
		}
	}()

	return hl.parseID(aID, aText)
} // IDparse()

// `IDremove()` deletes all @hashtags/@mentions associated with `aID`.
//
// Parameters:
// - `aID` is to be deleted from all lists.
func (hl *THashList) IDremove(aID uint64) *THashList {
	if (nil == hl) || (0 == len(hl.hm)) {
		return hl
	}
	hl.mtx.Lock()
	defer hl.mtx.Unlock()

	return hl.removeID(aID)
} // IDremove()

// `IDrename()` replaces all occurrences of `aOldID` by `aNewID`.
//
// This method is intended for rare cases when the ID of a document
// needs to get changed.
//
// Parameters:
// - `aOldID` is to be replaced in all lists.
// - `aNewID` is the replacement in all lists.
func (hl *THashList) IDrename(aOldID, aNewID uint64) *THashList {
	if (aOldID == aNewID) || (0 == len(hl.hm)) {
		return hl
	}
	hl.mtx.Lock()
	defer hl.mtx.Unlock()

	return hl.renameID(aOldID, aNewID)
} // IDrename()

// `IDupdate()` checks `aText` removing all #hashtags/@mentions no longer
// present and adding #hashtags/@mentions new in `aText`.
//
// Parameters:
// - `aID` is the ID to update.
// - `aText` is the text to use.
func (hl *THashList) IDupdate(aID uint64, aText []byte) *THashList {
	hl.mtx.Lock()
	defer hl.mtx.Unlock()

	return hl.updateID(aID, aText)
} // IDupdate()

// `Len()` returns the current length of the list i.e. how many #hashtags
// and @mentions are currently stored in the list.
func (hl *THashList) Len() int {
	hl.mtx.RLock()
	defer hl.mtx.RUnlock()

	return len(hl.hm)
} // Len()

// `LenTotal()` returns the length of all #hashtag/@mention lists stored
// in the hash list.
//
// Returns:
// - `int`: The total length of all #hashtag/@mention lists.
func (hl *THashList) LenTotal() (rLen int) {
	hl.mtx.RLock()
	defer hl.mtx.RUnlock()

	rLen = len(hl.hm)
	for _, sl := range hl.hm {
		rLen += len(*sl)
	}

	return
} // LenTotal()

func (hl *THashList) list() TCountList {
	if (hl.checksum() == hl.µCC.µCRC) && (0 < len(hl.µCC.µCounts)) {
		return hl.µCC.µCounts
	}

	hl.µCC.µCounts = nil
	hl.µCC.µCRC = hl.checksum()
	hl.µCC.µCounts = hl.hm.countedList()

	return hl.µCC.µCounts
} // countedList()

// `List()` returns a list of #hashtags/@mentions with
// their respective count of associated IDs.
//
// Returns:
// - `TCountList`: A list of #hashtags/@mentions with their
// respective counts of associated IDs.
func (hl *THashList) List() TCountList {
	hl.mtx.Lock()
	defer hl.mtx.Unlock()

	return hl.list()
} // List()

// `Load()` reads the configured file returning the data structure
// read from the file and a possible error condition.
//
// If the hash file doesn't exist that is not considered an error.
//
// Returns:
// - `*THashList`: The updated list.
// - `error`: If there is an error, it will be of type `*PathError`.
func (hl *THashList) Load() (*THashList, error) {
	hl.mtx.Lock()
	defer hl.mtx.Unlock()

	_, err := hl.hm.Load(hl.fn)
	return hl, err
} // Load()

// `MentionAdd()` appends `aID` to the list of `aMention`.
//
// If either `aMention` or `aID` are empty strings they are
// silently ignored (i.e. this method does nothing).
//
// Parameters:
// - `aMention` is the list index to lookup.
// - `aID` is to be added to the hash list.
func (hl *THashList) MentionAdd(aMention string, aID uint64) *THashList {
	if 0 == len(aMention) {
		return hl
	}
	hl.mtx.Lock()
	defer hl.mtx.Unlock()

	return hl.add(MarkMention, aMention, aID)
} // MentionAdd()

// `MentionCount()` returns the number of mentions in the list.
//
// Returns:
// - `int`: The number of mentions in the list.
func (hl *THashList) MentionCount() int {
	hl.mtx.RLock()
	defer hl.mtx.RUnlock()

	return hl.hm.count(MarkMention)
} // MentionCount()

// `MentionLen()` returns the number of IDs stored for `aMention`.
//
// Parameters:
// - `aMention` identifies the ID list to lookup.
//
// Returns:
// - `int`: The number of `aMention` in the list.
func (hl *THashList) MentionLen(aMention string) int {
	hl.mtx.RLock()
	defer hl.mtx.RUnlock()

	return hl.hm.idxLen(MarkMention, aMention)
} // MentionLen()

// `MentionList()` returns a list of IDs associated with `aMention`.
//
// Parameters:
// - `aMention`: The mention to lookup.
//
// Returns:
// - `[]uint64`: The number of references of `aName`.
func (hl *THashList) MentionList(aMention string) []uint64 {
	hl.mtx.Lock()
	defer hl.mtx.Unlock()

	return hl.hm.list(MarkMention, aMention)
} // MentionList()

// `MentionRemove()` deletes `aID` from the list of `aMention`.
//
// Parameters:
// - `aMention`: The mention to lookup.
// - `aID`: The referenced object to remove from the list.
//
// Returns:
// - `*tHashList`: The current hash list.
func (hl *THashList) MentionRemove(aMention string, aID uint64) *THashList {
	if 0 == len(aMention) {
		return hl
	}
	hl.mtx.RLock()
	defer hl.mtx.RUnlock()

	hl.removeHM(MarkMention, aMention, aID)

	return hl
} // MentionRemove()

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
func (hl *THashList) parseID(aID uint64, aText []byte) *THashList {
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
func (hl *THashList) removeHM(aDelim byte, aName string, aID uint64) *THashList {
	aName = strings.ToLower(strings.TrimSpace(aName))
	if (0 == len(aName)) || (0 == aID) {
		return hl
	}

	oldCRC := hl.checksum()
	defer func() {
		if oldCRC != atomic.LoadUint32(&hl.µChange) {
			go func() {
				hl.hm.store(hl.fn)
			}()
		}
	}()

	hl.hm.remove(aDelim, aName, aID)
	atomic.StoreUint32(&hl.µChange, 0)

	return hl
} // removeHM()

// `removeID()` deletes all #hashtags/@mentions associated with `aID`.
//
// Parameters:
// - `aID`: The object to remove from all references list.
//
// Returns:
// - `*THashList`: The modified hash list.
func (hl *THashList) removeID(aID uint64) *THashList {
	if (nil == hl) || (0 == len(hl.hm)) {
		return hl
	}

	oldCRC := hl.checksum()
	defer func() {
		if oldCRC != atomic.LoadUint32(&hl.µChange) {
			go func() {
				hl.hm.store(hl.fn)
			}()
		}
	}()

	hl.hm.removeID((aID))
	atomic.StoreUint32(&hl.µChange, 0)

	return hl
} // removeID()

func (hl *THashList) renameID(aOldID, aNewID uint64) *THashList {
	if (aOldID == aNewID) || (nil == hl) || (0 == len(hl.hm)) {
		return hl
	}

	oldCRC := hl.checksum()
	defer func() {
		if oldCRC != atomic.LoadUint32(&hl.µChange) {
			go func() {
				hl.hm.store(hl.fn)
			}()
		}
	}()

	for _, sl := range hl.hm {
		sl.rename(aOldID, aNewID)
	}
	atomic.StoreUint32(&hl.µChange, 0)

	return hl
} // renameID()

// `SetFilename()` sets `aFilename` to be used by this list.
//
// Parameters:
// - `aFilename`: The name of the file to use for storage.
//
// Returns:
// - `*THashList`: The current hash list.
func (hl *THashList) SetFilename(aFilename string) *THashList {
	hl.mtx.Lock()
	defer hl.mtx.Unlock()

	hl.fn = aFilename

	return hl
} // SetFilename()

// `Store()` writes the whole list to the configured file
// returning the number of bytes written and a possible error.
//
// If there is an error, it will be of type `*PathError`.
//
// Returns:
// - `int`: Number of bytes written to storage.
// - `error`: A possible storage error, or `nil` in case of success.
func (hl *THashList) Store() (int, error) {
	hl.mtx.RLock()
	defer hl.mtx.RUnlock()

	return hl.hm.store(hl.fn)
} // Store()

// `String()` returns the whole list as a linefeed separated string.
//
// Returns:
// - `string`: The string representation of this hash list.
func (hl *THashList) String() string {
	hl.mtx.RLock()
	defer hl.mtx.RUnlock()

	return hl.hm.String()
} // String()

// `updateID()` checks `aText` removing all #hashtags/@mentions no longer
// present and adds #hashtags/@mentions new in `aText`.
//
// Parameters:
// - `aID` is the ID to update.
// - `aText` is the text to use.
//
// Returns:
// - `*THashList`: The current hash list.
func (hl *THashList) updateID(aID uint64, aText []byte) *THashList {
	if (nil == hl) || (0 == len(aText)) || (0 == len(hl.hm)) {
		return hl
	}

	oldCRC := hl.checksum()
	defer func() {
		if oldCRC != atomic.LoadUint32(&hl.µChange) {
			hl.hm.store(hl.fn)
		}
	}()

	return hl.removeID(aID).parseID(aID, aText)
} // updateID()

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
func (hl *THashList) Walk(aFunc TWalkFunc) {
	hl.mtx.RLock()
	defer hl.mtx.RUnlock()

	oldCRC := hl.checksum()
	defer func() {
		if oldCRC != atomic.LoadUint32(&hl.µChange) {
			_, _ = hl.hm.store(hl.fn)
		}
	}()

	changed := hl.hm.walk(aFunc)
	if changed {
		atomic.StoreUint32(&hl.µChange, 0)
	}
} // Walk()

// `Walker()` traverses through all entries in the hash lists
// calling `aWalker` for each entry.
//
// Parameters:
// - `aWalker` is an object implementing the `IHashWalker` interface.
func (hl *THashList) Walker(aWalker IHashWalker) {
	hl.Walk(aWalker.Walk)
} // Walker()

/* EoF */
