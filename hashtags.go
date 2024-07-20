/*
Copyright © 2019, 2024  M.Watermann, 10247 Berlin, Germany

		All rights reserved
	EMail : <support@mwat.de>
*/
package hashtags

import (
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
		crc uint32     // current CRC
		cl  TCountList // last list of counted items
	}

	// `THashTags` is a list of `#hashtags` and `@mentions`
	// pointing to sources (i.e. IDs).
	THashTags struct {
		fn      string       // the filename to use
		hl      tHashList    // the actual map list of sources/IDs
		mtx     sync.RWMutex // safeguard against concurrent accesses
		changed uint32       // internal change flag
		cc      tCountCache  // cache for `CountedList()`
		safe    bool         // flag for optional thread safety
	}
)

var (
	// `UseBinaryStorage` determines whether to use binary storage
	// or not (i.e. plain text).
	//
	// Loading/storing binary data is about three times as fast with
	// the `THashTags` data than reading and parsing plain text data.
	UseBinaryStorage = true
)

// --------------------------------------------------------------------------
// constructor function

// `New()` returns a new `THashTags` instance after reading
// the given file.
//
// If the hash file doesn't exist that is not considered an error.
// If there is an error, it will be of type *PathError.
//
// Parameters:
//   - `aFilename`: The name of the file to use for loading and storing.
//
// Returns:
//   - `*THashTags`: The new `THashTags` instance.
//   - `error`: If there is an error, it will be from reading `aFilename`.
func New(aFilename string, aSafe bool) (*THashTags, error) {
	hashlist, err := newHashList(aFilename)
	if nil != err {
		return nil, err
	}
	ht := &THashTags{
		fn:   aFilename,
		hl:   *hashlist,
		safe: aSafe,
	}

	return ht, nil
} // New()

// -------------------------------------------------------------------------
// methods of THashTags

func (ht *THashTags) checksum() uint32 {
	if 0 == atomic.LoadUint32(&ht.changed) {
		atomic.StoreUint32(&ht.changed, ht.hl.checksum())
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

	ht.hl.clear()
	atomic.StoreUint32(&ht.changed, 0)

	return ht
} // Clear()

// `equals()` compares the current list with another list.
//
// Parameters:
//   - `aList`: The list to compare with.
//
// Returns:
//   - `bool`: True if the lists are identical, false otherwise.
func (ht *THashTags) equals(aList *THashTags) bool {
	if ht.safe {
		ht.mtx.Lock()
		defer ht.mtx.Unlock()
	}

	return ht.hl.equals(&aList.hl)
} // compareTo()

// `Filename()` returns the configured filename for reading/storing
// this list's contents.
//
// Returns:
//   - `string`: The filename for reading/storing this list.
func (ht *THashTags) Filename() string {
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
func (ht *THashTags) HashAdd(aHash string, aID uint64) bool {
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

	return ht.hl.hashCount()
} // HashCount()

// `HashLen()` returns the number of IDs stored for `aHash`.
//
// If `aHash` is empty it is silently ignored (i.e. this method
// does nothing), returning `-1`.
//
// Parameters:
//   - `aHash` The list key to lookup.
//
// Returns:
//   - `int`: The number of `aHash` in the list.
func (ht *THashTags) HashLen(aHash string) int {
	if ht.safe {
		ht.mtx.RLock()
		defer ht.mtx.RUnlock()
	}

	return ht.hl.hashLen(aHash)
} // HashLen()

// `HashList()` returns a list of IDs associated with `aHash`.
//
// If `aHash` is empty it is silently ignored (i.e. this method
// does nothing), returning an empty slice.
//
// Parameters:
//   - `aName`: The hash to lookup.
//
// Returns:
//   - `[]uint64`: The number of references of `aName`.
func (ht *THashTags) HashList(aHash string) []uint64 {
	if ht.safe {
		ht.mtx.RLock()
		defer ht.mtx.RUnlock()
	}

	return ht.hl.hashList(aHash)
} // HashList()

// `HashRemove()` deletes `aID` from the list of `aHash`.
//
// Parameters:
//   - `aHash`: The hash to lookup.
//   - `aID`: The referenced object to remove from the list.
//
// Returns:
//   - `bool`: `true` if `aID` was removed, or `false` otherwise.
func (ht *THashTags) HashRemove(aHash string, aID uint64) bool {
	if ht.safe {
		ht.mtx.Lock()
		defer ht.mtx.Unlock()
	}

	return ht.removeHM(MarkHash, aHash, aID)
} // HashRemove()

// `IDlist()` returns a list of #hashtags and @mentions associated with `aID`.
//
// Parameters:
//   - `aID`: The referenced object to lookup.
//
// Returns:
//   - `[]string`: The list of #hashtags and @mentions associated with `aID`.
func (ht *THashTags) IDlist(aID uint64) []string {
	if ht.safe {
		ht.mtx.RLock()
		defer ht.mtx.RUnlock()
	}

	return ht.hl.idList(aID)
} // IDlist()

// `IDparse()` checks whether `aText` associated with `aID` contains
// strings starting with `[@|#]` and - if found - adds them to the
// respective list.
//
// If `aText` is empty it is silently ignored (i.e. this method
// does nothing), returning `false`.
//
// Parameters:
//   - `aID`: the ID to add to the list.
//   - `aText:` The text to search.
//
// Returns:
//   - `bool`: `true` if `aID` was updated from `aText`, or `false` otherwise.
func (ht *THashTags) IDparse(aID uint64, aText []byte) bool {
	if ht.safe {
		ht.mtx.Lock()
		defer ht.mtx.Unlock()
	}

	oldCRC := ht.hl.checksum()
	defer func() {
		if oldCRC != atomic.LoadUint32(&ht.changed) {
			ht.hl.store(ht.fn) //TODO: call ht.hl.method
		}
	}()

	result := ht.hl.parseID(aID, aText)
	if result {
		atomic.StoreUint32(&ht.changed, 0)
	}

	return result
} // IDparse()

// `IDremove()` deletes all @hashtags/@mentions associated with `aID`.
//
// Parameters:
//   - `aID` is to be deleted from all lists.
//
// Returns:
//   - `bool`: `true` if `aID` was removed, or `false` otherwise.
func (ht *THashTags) IDremove(aID uint64) bool {
	if ht.safe {
		ht.mtx.Lock()
		defer ht.mtx.Unlock()
	}

	oldCRC := ht.hl.checksum()
	defer func() {
		if oldCRC != atomic.LoadUint32(&ht.changed) {
			go func() {
				ht.hl.store(ht.fn)
			}()
		}
	}()

	result := ht.hl.removeID(aID)
	if result {
		atomic.StoreUint32(&ht.changed, 0)
	}

	return result
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
//   - `aOldID` is to be replaced in all lists.
//   - `aNewID` is the replacement in all lists.
//
// Returns:
//   - `bool`: `true` if `aOldID` was renamed, or `false` otherwise.
func (ht *THashTags) IDrename(aOldID, aNewID uint64) bool {
	if aOldID == aNewID {
		return false
	}

	if ht.safe {
		ht.mtx.Lock()
		defer ht.mtx.Unlock()
	}

	oldCRC := ht.hl.checksum()
	defer func() {
		if oldCRC != atomic.LoadUint32(&ht.changed) {
			go func() {
				ht.hl.store(ht.fn)
			}()
		}
	}()

	result := ht.hl.renameID(aOldID, aNewID)
	if result {
		atomic.StoreUint32(&ht.changed, 0)
	}

	return result
} // IDrename()

// `IDupdate()` checks `aText` removing all #hashtags/@mentions no longer
// present and adding #hashtags/@mentions new in `aText`.
//
// Parameters:
//   - `aID`: The ID to update.
//   - `aText`: The new text to use.
//
// Returns:
//   - `bool`: `true` if `aID` was updated, or `false` otherwise.
func (ht *THashTags) IDupdate(aID uint64, aText []byte) bool {
	if ht.safe {
		ht.mtx.Lock()
		defer ht.mtx.Unlock()
	}

	oldCRC := ht.hl.checksum()
	defer func() {
		if oldCRC != atomic.LoadUint32(&ht.changed) {
			go func() {
				ht.hl.store(ht.fn)
			}()
		}
	}()

	result := ht.hl.updateID(aID, aText)
	if result {
		atomic.StoreUint32(&ht.changed, 0)
	}

	return result
} // IDupdate()

// `insert()` appends `aID` to the list associated with `aName`.
//
// If `aName` is empty it is silently ignored (i.e. this method
// does nothing) returning `false`.
//
// Parameters:
//   - `aDelim`: The start character of words to use (i.e. either '@' or '#').
//   - `aName`: The hashtag/mention to lookup.
//   - `aID`: The referencing object to be added to the hash list.
//
// Returns:
//   - `bool`: `true` if `aID` was added, or `false` otherwise.
func (ht *THashTags) insert(aDelim byte, aName string, aID uint64) bool {
	// prepare for case-insensitive search:
	aName = strings.ToLower(strings.TrimSpace(aName))
	if 0 == len(aName) {
		return false
	}

	oldCRC := ht.checksum()
	defer func() {
		if oldCRC != atomic.LoadUint32(&ht.changed) {
			go func() {
				ht.hl.store(ht.fn)
			}()
		}
	}()

	result := ht.hl.insert(aDelim, aName, aID)
	if result {
		atomic.StoreUint32(&ht.changed, 0)
	}

	return result
} // insert()

// `Len()` returns the current length of the list i.e. how many
// #hashtags and @mentions are currently stored in the list.
//
// Returns:
//   - `int`: The length of all #hashtag/@mention list.
func (ht *THashTags) Len() int {
	if ht.safe {
		ht.mtx.RLock()
		defer ht.mtx.RUnlock()
	}

	return ht.hl.len()
} // Len()

// `LenTotal()` returns the length of all #hashtag/@mention lists
// stored in the list.
//
// Returns:
//   - `int`: The total length of all #hashtag/@mention lists.
func (ht *THashTags) LenTotal() int {
	if ht.safe {
		ht.mtx.RLock()
		defer ht.mtx.RUnlock()
	}

	return ht.hl.lenTotal()
} // LenTotal()

// `List()` returns a list of #hashtags/@mentions with
// their respective count of associated IDs.
//
// Returns:
//   - `TCountList`: A list of #hashtags/@mentions with their counts of IDs.
func (ht *THashTags) List() TCountList {
	if ht.safe {
		ht.mtx.RLock()
		defer ht.mtx.RUnlock()
	}

	if (ht.hl.checksum() == ht.cc.crc) && (0 < len(ht.cc.cl)) {
		return ht.cc.cl
	}

	ht.cc.cl = nil
	ht.cc.crc = ht.hl.checksum()
	ht.cc.cl = ht.hl.countedList()

	return ht.cc.cl
} // List()

// `Load()` reads the configured file returning the data structure
// read from the file and a possible error condition.
//
// If the hash file doesn't exist that is not considered an error.
//
// Returns:
//   - `*THashTags`: The updated list.
//   - `error`: If there is an error, it will be of type `*PathError`.
func (ht *THashTags) Load() (*THashTags, error) {
	if ht.safe {
		ht.mtx.Lock()
		defer ht.mtx.Unlock()
	}

	oldCRC := ht.hl.checksum()
	defer func() {
		if oldCRC != atomic.LoadUint32(&ht.changed) {
			go func() {
				ht.hl.store(ht.fn)
			}()
		}
	}()

	_, err := ht.hl.load(ht.fn)
	if nil == err {
		atomic.StoreUint32(&ht.changed, 0)
	}

	return ht, err
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
func (ht *THashTags) MentionAdd(aMention string, aID uint64) bool {
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

	return ht.hl.mentionCount()
} // MentionCount()

// `MentionLen()` returns the number of IDs stored for `aMention`.
//
// If `aMention` is empty it is silently ignored (i.e. this method
// does nothing) returning `-1`.
//
// Parameters:
//   - `aMention` identifies the ID list to lookup.
//
// Returns:
//   - `int`: The number of `aMention` in the list.
func (ht *THashTags) MentionLen(aMention string) int {
	if ht.safe {
		ht.mtx.RLock()
		defer ht.mtx.RUnlock()
	}

	return ht.hl.mentionLen(aMention)
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
//   - `[]uint64`: The number of references of `aName`.
func (ht *THashTags) MentionList(aMention string) []uint64 {
	if ht.safe {
		ht.mtx.RLock()
		defer ht.mtx.RUnlock()
	}

	return ht.hl.mentionList(aMention)
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
func (ht *THashTags) MentionRemove(aMention string, aID uint64) bool {
	if ht.safe {
		ht.mtx.Lock()
		defer ht.mtx.Unlock()
	}

	return ht.removeHM(MarkMention, aMention, aID)
} // MentionRemove()

// `removeHM()` deletes `aID` from the list of `aName`.
//
// If `aName` is empty it is silently ignored (i.e. this method
// does nothing) returning `false`.
//
// Parameters:
//   - `aDelim` is the start character of words to use (i.e. either '@' or '#').
//   - `aName`: The hash/mention to lookup for `aID`.
//   - `aID`: The source to remove from the list.
//
// Returns:
//   - `bool`: `true` if `aID` was updated, or `false` otherwise.
func (ht *THashTags) removeHM(aDelim byte, aName string, aID uint64) bool {
	aName = strings.ToLower(strings.TrimSpace(aName))
	if 0 == len(aName) {
		return false
	}

	oldCRC := ht.hl.checksum()
	defer func() {
		if oldCRC != atomic.LoadUint32(&ht.changed) {
			go func() {
				ht.hl.store(ht.fn)
			}()
		}
	}()

	result := ht.hl.removeHM(aDelim, aName, aID)
	if result {
		atomic.StoreUint32(&ht.changed, 0)
	}

	return result
} // removeHM()

// `SetFilename()` sets `aFilename` to be used by this list.
//
// Parameters:
//   - `aFilename`: The name of the file to use for storage.
//
// Returns:
//   - `*THashList`: The current hash list.
func (ht *THashTags) SetFilename(aFilename string) *THashTags {
	if ht.safe {
		ht.mtx.Lock()
		defer ht.mtx.Unlock()
	}

	ht.fn = aFilename

	return ht
} // SetFilename()

// `Store()` writes the whole list to the configured file
// returning the number of bytes written and a possible error.
//
// If there is an error, it will be of type `*PathError`.
//
// Returns:
//   - `int`: Number of bytes written to storage.
//   - `error`: A possible storage error, or `nil` in case of success.
func (ht *THashTags) Store() (int, error) {
	if ht.safe {
		ht.mtx.RLock()
		defer ht.mtx.RUnlock()
	}

	return ht.hl.store(ht.fn)
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

	return ht.hl.String()
} // String()

/* EoF */
