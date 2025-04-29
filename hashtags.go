/*
Copyright © 2019, 2025  M.Watermann, 10247 Berlin, Germany

	    All rights reserved
	EMail : <support@mwat.de>
*/
package hashtags

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"

	se "github.com/mwat56/sourceerror"
)

//lint:file-ignore ST1017 - I prefer Yoda conditions

const (
	// `MarkHash` is the first character in a `#hashtag`.
	MarkHash = byte('#')

	// `MarkMention` is the first character in a `@mention`.
	MarkMention = byte('@')
)

type (
	// `tCountCache` is a data cache for `CountedList()`.
	tCountCache struct {
		crc uint32     // current CRC
		cl  TCountList // last list of counted items
	}

	// `THashTags` is a list of `#hashtags` and `@mentions`
	// pointing to sources (i.e. IDs).
	THashTags struct {
		mtx     sync.RWMutex // safeguard against concurrent accesses
		hm      *tHashMap    // the actual map list of sources/IDs
		fn      string       // the filename to use
		cc      tCountCache  // cache for `CountedList()`
		changed uint32       // internal change flag
		safe    bool         // flag for optional thread safety
	}

	// `THashTagError` is a custom error.
	// Deprecated: Use [sourceerror.ErrSource] instead.
	THashTagError = se.ErrSource
)

var (
	// `UseBinaryStorage` determines whether to use binary storage
	// or not (i.e. plain text).
	//
	// Loading/storing binary data is about three times as fast with
	// the `THashTags` data than reading and parsing plain text data.
	UseBinaryStorage = true

	// RegEx to identify a numeric HTML entity.
	htEntityRE = regexp.MustCompile(`#[0-9]+;`)

	// match: #hashtag|@mention
	htHashMentionRE = regexp.MustCompile(
		`(?ims)(?:^|\s|[^\p{L}\d_])?([@#][\p{L}’'\d_§-]+)(?:[^\p{L}\d_]|$)`)
	//	                             1111111111111111111  222222222222222

	// RegEx to match texts like `#----`.
	htHyphenRE = regexp.MustCompile(`#[^-]*--`)
)

// --------------------------------------------------------------------------
// constructor function:

// `New()` returns a new `THashTags` instance after reading
// the given file.
//
// NOTE: An empty filename or if the hash file doesn't exist is not
// considered an error.
//
// Parameters:
//   - `aFilename`: The name of the file to use for loading and storing.
//   - `aSafe`: Whether to use thread-safe operations or not.
//
// Returns:
//   - `*THashTags`: The new `THashTags` instance.
//   - `error`: `nil` in case of success, otherwise an error.
func New(aFilename string, aSafe bool) (*THashTags, error) {

	ht := &THashTags{
		hm:   newHashMap(),
		safe: aSafe,
	}
	if aFilename = strings.TrimSpace(aFilename); "" == aFilename {
		return ht, nil
	}
	ht.fn = aFilename

	_, err := ht.hm.load(aFilename) // err already wrapped

	return ht, err
} // New()

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

// -------------------------------------------------------------------------
// methods of `THashTags`:

// `checksum()` returns the list's CRC32 checksum.
//
// This internal method is used by the public `Checksum()` method to get
// a kind of 'footprint' of the current contents of the handled data.
//
// Returns:
//   - `uint32`: The computed checksum.
func (ht *THashTags) checksum() uint32 {
	if 0 == atomic.LoadUint32(&ht.changed) {
		atomic.StoreUint32(&ht.changed, ht.hm.checksum())
	}

	return atomic.LoadUint32(&ht.changed)
} // checksum()

// `Checksum()` returns the list's CRC32 checksum.
//
// This method can be used to get a kind of 'footprint' of the current
// contents of the handled data.
//
// Returns:
//   - `uint32`: The computed checksum.
func (ht *THashTags) Checksum() uint32 {
	if ht.safe {
		ht.mtx.RLock()
		defer ht.mtx.RUnlock()
	}

	return ht.checksum()
} // Checksum()

// `Clear()` empties the internal data structures:
// all `#hashtags` and `@mentions` are deleted.
//
// Returns:
//   - `*THashTags`: This cleared list.
func (ht *THashTags) Clear() *THashTags {
	if ht.safe {
		ht.mtx.Lock()
		defer ht.mtx.Unlock()
	}

	ht.hm.clear()
	atomic.StoreUint32(&ht.changed, 0)

	return ht
} // Clear()

// `deferredStore()` returns a closure that, when executed, checks whether
// the list's contents have changed and if so, asynchronously stores
// the current state to the configured file.
//
// This method is meant to be used internally with the `defer` statement.
//
// Returns:
//   - `func()`: A closure that handles deferred storage operations.
func (ht *THashTags) deferredStore() func() {
	oldCRC := ht.hm.checksum()

	return func() {
		if oldCRC != atomic.LoadUint32(&ht.changed) {
			go ht.hm.store(ht.fn)
		}
	}
} // deferredStore()

// `Filename()` returns the configured filename for reading/storing
// this list's contents.
//
// Returns:
//   - `string`: The filename for reading/storing this list.
func (ht *THashTags) Filename() string {
	if ht.safe {
		ht.mtx.RLock()
		defer ht.mtx.RUnlock()
	}

	return ht.fn
} // Filename()

// `HashAdd()` appends `aID` to the list of `aHash`.
//
// If `aHash` is empty it is silently ignored (i.e. this method
// does nothing) returning `false`.
//
// Parameters:
//   - `aHash`: The hash list index to use.
//   - `aID`: The object to be added to the hash list.
//
// Returns:
//   - `bool`: `true` if `aID` was added, or `false` otherwise.
func (ht *THashTags) HashAdd(aHash string, aID int64) bool {
	if aHash = strings.TrimSpace(aHash); "" == aHash {
		return false
	}

	if ht.safe {
		ht.mtx.Lock()
		defer ht.mtx.Unlock()
	}

	return ht.insert(MarkHash, aHash, aID)
} // HashAdd()

// `HashCount()` counts the number of hashtags in the list.
//
// Returns:
//   - `int`: The number of hashes in the list.
func (ht *THashTags) HashCount() int {
	if ht.safe {
		ht.mtx.RLock()
		defer ht.mtx.RUnlock()
	}

	return ht.hm.count(MarkHash)
} // HashCount()

// `HashLen()` returns the number of IDs stored for `aHash`.
//
// If `aHash` is empty it is silently ignored (i.e. this method
// does nothing), returning `-1`.
//
// Parameters:
//   - `aHash`: The list key to lookup.
//
// Returns:
//   - `int`: The number of `aHash` in the list.
func (ht *THashTags) HashLen(aHash string) int {
	if aHash = strings.TrimSpace(aHash); "" == aHash {
		return 0
	}

	if ht.safe {
		ht.mtx.RLock()
		defer ht.mtx.RUnlock()
	}

	return ht.hm.idxLen(MarkHash, aHash)
} // HashLen()

// `HashList()` returns a list of IDs associated with `aHash`.
//
// If `aHash` is empty it is silently ignored (i.e. this method
// does nothing), returning an empty slice.
//
// Parameters:
//   - `aHash`: The hash to lookup.
//
// Returns:
//   - `[]int64`: The number of references of `aHash`.
func (ht *THashTags) HashList(aHash string) []int64 {
	if aHash = strings.TrimSpace(aHash); "" == aHash {
		return []int64{}
	}

	if ht.safe {
		ht.mtx.RLock()
		defer ht.mtx.RUnlock()
	}

	return ht.hm.list(MarkHash, aHash)
} // HashList()

// `HashRemove()` deletes `aID` from the list of `aHash`.
//
// Parameters:
//   - `aHash`: The hash to lookup.
//   - `aID`: The referenced object to remove from the list.
//
// Returns:
//   - `bool`: `true` if `aID` was removed, or `false` otherwise.
func (ht *THashTags) HashRemove(aHash string, aID int64) bool {
	if aHash = strings.TrimSpace(aHash); "" == aHash {
		return false
	}

	if ht.safe {
		ht.mtx.Lock()
		defer ht.mtx.Unlock()
	}

	return ht.removeHM(MarkHash, aHash, aID)
} // HashRemove()

// `IDlist()` returns a list of `#hashtags` and `@mentions` associated
// with `aID`.
//
// Parameters:
//   - `aID`: The referenced object to lookup.
//
// Returns:
//   - `[]string`: The list of `#hashtags` and `@mentions` associated with `aID`.
func (ht *THashTags) IDlist(aID int64) []string {
	if ht.safe {
		ht.mtx.RLock()
		defer ht.mtx.RUnlock()
	}

	return ht.hm.idList(aID)
} // IDlist()

// `IDparse()` checks whether `aText` associated with `aID` contains
// strings starting with `[@|#]` and - if found - adds them to the
// respective list.
//
// If `aText` is empty it is silently ignored (i.e. this method
// does nothing), returning `false`.
//
// Parameters:
//   - `aID`: The ID to add to the list.
//   - `aText:` The text to search.
//
// Returns:
//   - `bool`: `true` if `aID` was updated from `aText`, or `false` otherwise.
func (ht *THashTags) IDparse(aID int64, aText []byte) bool {
	if aText = bytes.TrimSpace(aText); 0 == len(aText) {
		return false
	}

	if ht.safe {
		ht.mtx.Lock()
		defer ht.mtx.Unlock()
	}
	defer ht.deferredStore()

	if ht.parseID(aID, aText) {
		atomic.StoreUint32(&ht.changed, 0)
		return true
	}

	return false
} // IDparse()

// `IDremove()` deletes all `#hashtags` and `@mentions` associated with `aID`.
//
// Parameters:
//   - `aID`: The ID to be deleted from all lists.
//
// Returns:
//   - `bool`: `true` if `aID` was removed, or `false` otherwise.
func (ht *THashTags) IDremove(aID int64) bool {
	if ht.safe {
		ht.mtx.Lock()
		defer ht.mtx.Unlock()
	}
	defer ht.deferredStore()

	if ht.hm.removeID(aID) {
		atomic.StoreUint32(&ht.changed, 0)
		return true
	}

	return false
} // IDremove()

// `IDrename()` replaces all occurrences of `aOldID` by `aNewID`.
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
//   - `bool`: `true` if `aOldID` was renamed, or `false` otherwise.
func (ht *THashTags) IDrename(aOldID, aNewID int64) bool {
	if aOldID == aNewID {
		return false
	}

	if ht.safe {
		ht.mtx.Lock()
		defer ht.mtx.Unlock()
	}
	defer ht.deferredStore()

	if ht.hm.renameID(aOldID, aNewID) {
		atomic.StoreUint32(&ht.changed, 0)
		return true
	}

	return false
} // IDrename()

// `IDupdate()` checks `aText` removing all `#hashtags` and `@mentions`
// no longer present and adding `#hashtags` and `@mentions` new in `aText`.
//
// Parameters:
//   - `aID`: The ID to update.
//   - `aText`: The new text to use.
//
// Returns:
//   - `bool`: `true` if `aID` was updated, or `false` otherwise.
func (ht *THashTags) IDupdate(aID int64, aText []byte) bool {
	if aText = bytes.TrimSpace(aText); 0 == len(aText) {
		return ht.IDremove(aID)
	}

	if ht.safe {
		ht.mtx.Lock()
		defer ht.mtx.Unlock()
	}
	defer ht.deferredStore()

	rr := ht.hm.removeID(aID)
	rp := ht.parseID(aID, aText)

	if rr || rp {
		atomic.StoreUint32(&ht.changed, 0)
		return true
	}

	return false
} // IDupdate()

// `insert()` appends `aID` to the list associated with `aName`.
//
// If `aName` is empty it is silently ignored (i.e. this method
// does nothing) returning `false`.
//
// Parameters:
//   - `aDelim`: The start character of words to use (i.e. either '@' or '#').
//   - `aName`: The `#hashtag` or `@mention` to lookup.
//   - `aID`: The referencing object to be added to the hash list.
//
// Returns:
//   - `bool`: `true` if `aID` was added, or `false` otherwise.
func (ht *THashTags) insert(aDelim byte, aName string, aID int64) bool {
	if aName = strings.TrimSpace(aName); "" == aName {
		return false
	}

	if aName[0] != aDelim {
		aName = string(aDelim) + aName
	}
	defer ht.deferredStore()

	if ht.hm.insert(aName, aID) {
		atomic.StoreUint32(&ht.changed, 0)
		return true
	}

	return false
} // insert()

// `Len()` returns the current length of the list i.e. how many
// `#hashtags` and `@mentions` are currently stored in the list.
//
// Returns:
//   - `int`: The number of all `#hashtag` and `@mention` lists.
func (ht *THashTags) Len() int {
	if ht.safe {
		ht.mtx.RLock()
		defer ht.mtx.RUnlock()
	}

	return len(*ht.hm)
} // Len()

// `LenTotal()` returns the length of all `#hashtag` and `@mention`
// lists stored in the list.
//
// Returns:
//   - `int`: The total length of all `#hashtag` and `@mention` lists.
func (ht *THashTags) LenTotal() int {
	if ht.safe {
		ht.mtx.RLock()
		defer ht.mtx.RUnlock()
	}

	return ht.hm.lenTotal()
} // LenTotal()

// `List()` returns a list of `#hashtags` and `@mentions` with their
// respective count of associated IDs.
//
// Returns:
//   - `TCountList`: A list of `#hashtags` and `@mentions` with their counts of IDs.
func (ht *THashTags) List() TCountList {
	if ht.safe {
		ht.mtx.RLock()
		defer ht.mtx.RUnlock()
	}

	if (ht.hm.checksum() == ht.cc.crc) && (0 < len(ht.cc.cl)) {
		return ht.cc.cl
	}

	ht.cc.cl = nil
	ht.cc.crc = ht.hm.checksum()
	ht.cc.cl = ht.hm.countedList()

	return ht.cc.cl
} // List()

// `Load()` reads the configured file returning the data structure
// read from the file and a possible error condition.
//
// NOTE: An empty filename or the hash file doesn't exist that is not
// considered an error.
//
// Returns:
//   - `*THashTags`: The updated list.
//   - `error`: `nil` in case of success, otherwise an error.
func (ht *THashTags) Load() (*THashTags, error) {
	if ht.safe {
		ht.mtx.Lock()
		defer ht.mtx.Unlock()
	}
	defer ht.deferredStore()

	if _, err := ht.hm.load(ht.fn); nil != err {
		return ht, err
	}
	atomic.StoreUint32(&ht.changed, 0)

	return ht, nil
} // Load()

// `MentionAdd()` appends `aID` to the list of `aMention`.
//
// If `aMention` is empty it is silently ignored (i.e. this method
// does nothing) returning `false`.
//
// Parameters:
//   - `aMention`: The list index to lookup.
//   - `aID`: The ID to be added to the hash list.
//
// Returns:
//   - `bool`: `true` if `aID` was added, or `false` otherwise.
func (ht *THashTags) MentionAdd(aMention string, aID int64) bool {
	if aMention = strings.TrimSpace(aMention); "" == aMention {
		return false
	}

	if ht.safe {
		ht.mtx.Lock()
		defer ht.mtx.Unlock()
	}

	return ht.insert(MarkMention, aMention, aID)
} // MentionAdd()

// `MentionCount()` returns the number of mentions in the list.
//
// Returns:
//   - `int`: The number of mentions in the list.
func (ht *THashTags) MentionCount() int {
	if ht.safe {
		ht.mtx.RLock()
		defer ht.mtx.RUnlock()
	}

	return ht.hm.count(MarkMention)
} // MentionCount()

// `MentionLen()` returns the number of IDs stored for `aMention`.
//
// If `aMention` is empty it is silently ignored (i.e. this method
// does nothing) returning `-1`.
//
// Parameters:
//   - `aMention`: Identifies the ID list to lookup.
//
// Returns:
//   - `int`: The number of `aMention` in the list.
func (ht *THashTags) MentionLen(aMention string) int {
	if aMention = strings.TrimSpace(aMention); "" == aMention {
		return 0
	}

	if ht.safe {
		ht.mtx.RLock()
		defer ht.mtx.RUnlock()
	}

	return ht.hm.idxLen(MarkMention, aMention)
} // MentionLen()

// `MentionList()` returns a list of IDs associated with `aMention`.
//
// If `aMention` is empty it is silently ignored (i.e. this method
// does nothing), returning an empty slice.
//
// Parameters:
//   - `aMention`: The mention to lookup.
//
// Returns:
//   - `[]int64`: The number of references of `aMention`.
func (ht *THashTags) MentionList(aMention string) []int64 {
	if aMention = strings.TrimSpace(aMention); "" == aMention {
		return []int64{}
	}

	if ht.safe {
		ht.mtx.RLock()
		defer ht.mtx.RUnlock()
	}

	return ht.hm.list(MarkMention, aMention)
} // MentionList()

// `MentionRemove()` deletes `aID` from the list of `aMention`.
//
// If `aMention` is empty it is silently ignored (i.e. this method
// does nothing) returning `false`.
//
// Parameters:
//   - `aMention`: The mention to lookup.
//   - `aID`: The referenced object to remove from the list.
//
// Returns:
//   - `bool`: `true` if `aID` was removed, or `false` otherwise.
func (ht *THashTags) MentionRemove(aMention string, aID int64) bool {
	if aMention = strings.TrimSpace(aMention); "" == aMention {
		return false
	}

	if ht.safe {
		ht.mtx.Lock()
		defer ht.mtx.Unlock()
	}

	return ht.removeHM(MarkMention, aMention, aID)
} // MentionRemove()

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
func (ht *THashTags) parseID(aID int64, aText []byte) bool {
	matches := htHashMentionRE.FindAllSubmatch(aText, -1)
	if (nil == matches) || (0 == len(matches)) {
		return false
	}

	var (
		tag, match0 string
		sub         [][]byte
		result      bool
	)
	for _, sub = range matches {
		match0 = string(sub[0])
		tag = string(sub[1])

		if '_' == tag[len(tag)-1] {
			// '_' can be both, part of the hashtag and italic
			// markup so we must remove it if it's at the end:
			tag = tag[:len(tag)-1]
		}
		if MarkHash == tag[0] {
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
			if htHyphenRE.MatchString(tag) {
				continue
			}
		} else if MarkMention == tag[0] {
			if '.' == match0[len(match0)-1] {
				// we assume that it's an email address
				continue
			}
		}
		if ht.insert(tag[0], tag, aID) {
			result = true
		}
	}

	return result
} // parseID()

// `removeHM()` deletes `aID` from the list of `aName`.
//
// If `aName` is empty it is silently ignored (i.e. this method
// does nothing) returning `false`.
//
// Parameters:
//   - `aDelim`: The start character of words to use (i.e. either '@' or '#').
//   - `aName`: The hash/mention to lookup for `aID`.
//   - `aID`: The source to remove from the list.
//
// Returns:
//   - `bool`: `true` if `aID` was updated, or `false` otherwise.
func (ht *THashTags) removeHM(aDelim byte, aName string, aID int64) bool {
	if aName = strings.TrimSpace(aName); "" == aName {
		return false
	}

	defer ht.deferredStore()

	if ht.hm.removeHM(aDelim, aName, aID) {
		atomic.StoreUint32(&ht.changed, 0)
		return true
	}

	return false
} // removeHM()

// `SetFilename()` sets `aFilename` to be used by this list.
//
// Parameters:
//   - `aFilename`: The name of the file to use for storage.
//
// Returns:
//   - `error`: `nil` in case of success, otherwise an error.
func (ht *THashTags) SetFilename(aFilename string) error {
	if aFilename = strings.TrimSpace(aFilename); "" == aFilename {
		return se.New(errors.New("empty filename not allowed"), 1)
	}

	// Check if directory exists and is writeable
	dir := filepath.Dir(aFilename)
	if _, err := os.Stat(dir); nil != err {
		return se.New(err, 1)
	}

	if ht.safe {
		ht.mtx.Lock()
		defer ht.mtx.Unlock()
	}
	ht.fn = aFilename

	return nil
} // SetFilename()

// `Store()` writes the whole list to the configured file
// returning the number of bytes written and a possible error.
//
// Returns:
//   - `int`: Number of bytes written to storage.
//   - `error`: A possible storage error, or `nil` in case of success.
func (ht *THashTags) Store() (int, error) {
	if ht.safe {
		ht.mtx.RLock()
		defer ht.mtx.RUnlock()
	}

	return ht.hm.store(ht.fn)
} // Store()

// `String()` returns the whole list as a linefeed separated string.
//
// Returns:
//   - `string`: The string representation of this hash list.
func (ht *THashTags) String() string {
	if ht.safe {
		ht.mtx.RLock()
		defer ht.mtx.RUnlock()
	}

	return ht.hm.String()
} // String()

/* EoF */
