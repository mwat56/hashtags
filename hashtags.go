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
		µCRC    uint32
		µCounts TCountList
	}

	// `THashTags` is a list of `#hashtags` and `@mentions`
	// pointing to sources (i.e. IDs).
	THashTags struct {
		fn      string       // the filename to use
		hl      tHashList    // the actual map list of sources/IDs
		mtx     sync.RWMutex // safeguard against concurrent accesses
		µChange uint32       // internal change flag
		µCC     tCountCache  // cache for `CountedList()`
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

// `NewHashTags()` returns a new `THashTags` instance after reading
// the given file.
//
// If the hash file doesn't exist that is not considered an error.
// If there is an error, it will be of type *PathError.
//
// Parameters:
// - `aFilename` is the name of the file to use for reading and storing.
func NewHashTags(aFilename string) (*THashTags, error) {
	hashlist, err := newHashList(aFilename)
	if nil != err {
		return nil, err
	}
	ht := &THashTags{
		fn: aFilename,
		hl: *hashlist,
	}

	return ht, nil
} // NewHashTags()

// -------------------------------------------------------------------------
// methods of THashTags

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
// - `*THashTags`: This hash list.
func (ht *THashTags) add(aDelim byte, aName string, aID uint64) *THashTags {
	// prepare for case-insensitive search:
	aName = strings.ToLower(strings.TrimSpace(aName))
	if 0 == len(aName) {
		return ht
	}

	oldCRC := ht.hl.checksum()
	defer func() {
		if oldCRC != atomic.LoadUint32(&ht.µChange) {
			go func() {
				ht.hl.Store(ht.fn)
			}()
		}
	}()

	if aName[0] != aDelim {
		aName = string(aDelim) + aName
	}
	ht.hl.add0(aName, aID)
	atomic.StoreUint32(&ht.µChange, 0)

	return ht
} // add()

// `Checksum()` returns the list's CRC32 checksum.
//
// This method can be used to get a kind of 'footprint'.
//
// Returns:
// - `uint32`: The computed checksum.
func (ht *THashTags) Checksum() uint32 {
	ht.mtx.RLock()
	defer ht.mtx.RUnlock()

	if 0 == atomic.LoadUint32(&ht.µChange) {
		atomic.StoreUint32(&ht.µChange, ht.hl.checksum())
	}

	return atomic.LoadUint32(&ht.µChange)
} // Checksum()

// `Clear()` empties the internal data structures:
// all `#hashtags` and `@mentions` are deleted.
//
// Returns:
// - `*THashTags`: This cleared list.
func (ht *THashTags) Clear() *THashTags {
	ht.mtx.Lock()
	defer ht.mtx.Unlock()

	ht.hl.clear()
	atomic.StoreUint32(&ht.µChange, 0)

	return ht
} // Clear()

// ` compareTo()` compares the current list with another list.
//
// Parameters:
// - `aList`: The list to compare with.
//
// Returns:
// - `bool`: True if the lists are identical, false otherwise.
func (ht *THashTags) compareTo(aList *THashTags) bool {
	ht.mtx.Lock()
	defer ht.mtx.Unlock()

	return ht.hl.compareTo(&aList.hl)
} // compareTo()

// `Filename()` returns the configured filename for reading/storing
// this list's contents.
//
// Returns:
// - `string`: The filename for reading/storing this list.
func (ht *THashTags) Filename() string {
	return ht.fn
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
// - `*THashTags`: This hash list.
func (ht *THashTags) HashAdd(aHash string, aID uint64) *THashTags {
	if 0 == len(aHash) {
		return ht
	}
	ht.mtx.Lock()
	defer ht.mtx.Unlock()

	ht.add(MarkHash, aHash, aID)

	return ht
} // HashAdd()

// `HashCount()` counts the number of hashtags in the list.
//
// Returns:
// - `int`: The number of hashes in the list.
func (ht *THashTags) HashCount() int {
	ht.mtx.RLock()
	defer ht.mtx.RUnlock()

	return ht.hl.HashCount()
} // HashCount()

// `HashLen()` returns the number of IDs stored for `aHash`.
//
// Parameters:
// - `aHash` The list key to lookup.
//
// Returns:
// - `int`: The number of `aHash` in the list.
func (ht *THashTags) HashLen(aHash string) int {
	ht.mtx.RLock()
	defer ht.mtx.RUnlock()

	return ht.hl.HashLen(aHash)
} // HashLen()

// `HashList()` returns a list of IDs associated with `aHash`.
//
// Parameters:
// - `aName`: The hash to lookup.
//
// Returns:
// - `[]uint64`: The number of references of `aName`.
func (ht *THashTags) HashList(aHash string) []uint64 {
	ht.mtx.RLock()
	defer ht.mtx.RUnlock()

	return ht.hl.HashList(aHash)
} // HashList()

// `HashRemove()` deletes `aID` from the list of `aHash`.
//
// Parameters:
// - `aHash`: The hash to lookup.
// - `aID`: The referenced object to remove from the list.
//
// Returns:
// - `*THashTags`: The current hash list.
func (ht *THashTags) HashRemove(aHash string, aID uint64) *THashTags {
	if 0 == len(aHash) {
		return ht
	}
	ht.mtx.RLock()
	defer ht.mtx.RUnlock()

	return ht.removeHM(MarkHash, aHash, aID)
} // HashRemove()

// `IDlist()` returns a list of #hashtags and @mentions associated with `aID`.
//
// Parameters:
// - `aID`: The referenced object to lookup.
//
// Returns:
// - `[]string`: The list of #hashtags and @mentions associated with `aID`.
func (ht *THashTags) IDlist(aID uint64) []string {
	ht.mtx.RLock()
	defer ht.mtx.RUnlock()

	return ht.hl.idList(aID)
} // IDlist()

// `IDparse()` checks whether `aText` associated with `aID` contains
// strings starting with `[@|#]` and - if found - adds them to the
// respective list.
//
// Parameters:
// - `aID`: the ID to add to the list.
// - `aText:` The text to search.
//
// Returns:
// - `*THashTags`: The updated list.
func (ht *THashTags) IDparse(aID uint64, aText []byte) *THashTags {
	if 0 == len(aText) {
		return ht
	}
	ht.mtx.Lock()
	defer ht.mtx.Unlock()

	oldCRC := ht.hl.checksum()
	defer func() {
		if oldCRC != atomic.LoadUint32(&ht.µChange) {
			ht.hl.Store(ht.fn) //TODO: call ht.hl.method
		}
	}()

	ht.hl.parseID(aID, aText)
	atomic.StoreUint32(&ht.µChange, 0)

	return ht
} // IDparse()

// `IDremove()` deletes all @hashtags/@mentions associated with `aID`.
//
// Parameters:
// - `aID` is to be deleted from all lists.
//
// Returns:
// - `*THashTags`: The updated list.
func (ht *THashTags) IDremove(aID uint64) *THashTags {
	if nil == ht {
		return ht
	}
	ht.mtx.Lock()
	defer ht.mtx.Unlock()

	oldCRC := ht.hl.checksum()
	defer func() {
		if oldCRC != atomic.LoadUint32(&ht.µChange) {
			go func() {
				ht.hl.Store(ht.fn)
			}()
		}
	}()

	ht.hl.IDremove(aID)
	atomic.StoreUint32(&ht.µChange, 0)

	return ht
} // IDremove()

// `IDrename()` replaces all occurrences of `aOldID` by `aNewID`.
//
// This method is intended for rare cases when the ID of a document
// needs to get changed.
//
// Parameters:
// - `aOldID` is to be replaced in all lists.
// - `aNewID` is the replacement in all lists.
//
// Returns:
// - `*THashTags`: The updated list.
func (ht *THashTags) IDrename(aOldID, aNewID uint64) *THashTags {
	if aOldID == aNewID {
		return ht
	}
	ht.mtx.Lock()
	defer ht.mtx.Unlock()

	oldCRC := ht.hl.checksum()
	defer func() {
		if oldCRC != atomic.LoadUint32(&ht.µChange) {
			go func() {
				ht.hl.Store(ht.fn)
			}()
		}
	}()

	ht.hl.IDrename(aOldID, aNewID)
	atomic.StoreUint32(&ht.µChange, 0)

	return ht
} // IDrename()

// `IDupdate()` checks `aText` removing all #hashtags/@mentions no longer
// present and adding #hashtags/@mentions new in `aText`.
//
// Parameters:
// - `aID`: The ID to update.
// - `aText`: The new text to use.
//
// Returns:
// - `*THashTags`: The updated list.
func (ht *THashTags) IDupdate(aID uint64, aText []byte) *THashTags {
	ht.mtx.Lock()
	defer ht.mtx.Unlock()

	oldCRC := ht.hl.checksum()
	defer func() {
		if oldCRC != atomic.LoadUint32(&ht.µChange) {
			go func() {
				ht.hl.Store(ht.fn)
			}()
		}
	}()

	ht.hl.IDupdate(aID, aText)
	atomic.StoreUint32(&ht.µChange, 0)

	return ht
} // IDupdate()

// `Len()` returns the current length of the list i.e. how many #hashtags
// and @mentions are currently stored in the list.
//
// Returns:
// - `int`: The length of all #hashtag/@mention list.
func (ht *THashTags) Len() int {
	ht.mtx.RLock()
	defer ht.mtx.RUnlock()

	return ht.hl.Len()
} // Len()

// `LenTotal()` returns the length of all #hashtag/@mention lists stored
// in the hash list.
//
// Returns:
// - `int`: The total length of all #hashtag/@mention lists.
func (ht *THashTags) LenTotal() (rLen int) {
	ht.mtx.RLock()
	defer ht.mtx.RUnlock()

	return ht.hl.LenTotal()
} // LenTotal()

// `List()` returns a list of #hashtags/@mentions with
// their respective count of associated IDs.
//
// Returns:
// - `TCountList`: A list of #hashtags/@mentions with their counts of IDs.
func (ht *THashTags) List() TCountList {
	if (ht.hl.checksum() == ht.µCC.µCRC) && (0 < len(ht.µCC.µCounts)) {
		return ht.µCC.µCounts
	}
	ht.mtx.RLock()
	defer ht.mtx.RUnlock()

	ht.µCC.µCounts = nil
	ht.µCC.µCRC = ht.hl.checksum()
	ht.µCC.µCounts = ht.hl.List()

	return ht.µCC.µCounts
} // List()

// `Load()` reads the configured file returning the data structure
// read from the file and a possible error condition.
//
// If the hash file doesn't exist that is not considered an error.
//
// Returns:
// - `*THashTags`: The updated list.
// - `error`: If there is an error, it will be of type `*PathError`.
func (ht *THashTags) Load() (*THashTags, error) {
	ht.mtx.Lock()
	defer ht.mtx.Unlock()

	_, err := ht.hl.Load(ht.fn)

	return ht, err
} // Load()

// `MentionAdd()` appends `aID` to the list of `aMention`.
//
// If either `aMention` or `aID` are empty strings they are
// silently ignored (i.e. this method does nothing).
//
// Parameters:
// - `aMention`: The list index to lookup.
// - `aID`: The ID to be added to the hash list.
//
// Returns:
// - `*THashTags`: The updated list.
func (ht *THashTags) MentionAdd(aMention string, aID uint64) *THashTags {
	if 0 == len(aMention) {
		return ht
	}
	ht.mtx.Lock()
	defer ht.mtx.Unlock()

	ht.add(MarkMention, aMention, aID)

	return ht
} // MentionAdd()

// `MentionCount()` returns the number of mentions in the list.
//
// Returns:
// - `int`: The number of mentions in the list.
func (ht *THashTags) MentionCount() int {
	ht.mtx.RLock()
	defer ht.mtx.RUnlock()

	return ht.hl.MentionCount()
} // MentionCount()

// `MentionLen()` returns the number of IDs stored for `aMention`.
//
// Parameters:
// - `aMention` identifies the ID list to lookup.
//
// Returns:
// - `int`: The number of `aMention` in the list.
func (ht *THashTags) MentionLen(aMention string) int {
	ht.mtx.RLock()
	defer ht.mtx.RUnlock()

	return ht.hl.MentionLen(aMention)
} // MentionLen()

// `MentionList()` returns a list of IDs associated with `aMention`.
//
// Parameters:
// - `aMention`: The mention to lookup.
//
// Returns:
// - `[]uint64`: The number of references of `aName`.
func (ht *THashTags) MentionList(aMention string) []uint64 {
	ht.mtx.Lock()
	defer ht.mtx.Unlock()

	return ht.hl.MentionList(aMention)
} // MentionList()

// `MentionRemove()` deletes `aID` from the list of `aMention`.
//
// Parameters:
// - `aMention`: The mention to lookup.
// - `aID`: The referenced object to remove from the list.
//
// Returns:
// - `*THashTags`: The current hash list.
func (ht *THashTags) MentionRemove(aMention string, aID uint64) *THashTags {
	if 0 == len(aMention) {
		return ht
	}
	ht.mtx.RLock()
	defer ht.mtx.RUnlock()

	return ht.removeHM(MarkMention, aMention, aID)
} // MentionRemove()

// `removeHM()` deletes `aID` from the list of `aName`.
//
// Parameters:
// - `aDelim` is the start character of words to use (i.e. either '@' or '#').
// - `aName`: The hash/mention to lookup for `aID`.
// - `aID`: The source to remove from the list.
//
// Returns:
// - `*THashList`: The current hash list.
func (ht *THashTags) removeHM(aDelim byte, aName string, aID uint64) *THashTags {
	aName = strings.ToLower(strings.TrimSpace(aName))
	if (0 == len(aName)) || (0 == aID) {
		return ht
	}

	oldCRC := ht.hl.checksum()
	defer func() {
		if oldCRC != atomic.LoadUint32(&ht.µChange) {
			go func() {
				ht.hl.Store(ht.fn)
			}()
		}
	}()

	ht.hl.removeHM(aDelim, aName, aID)
	atomic.StoreUint32(&ht.µChange, 0)

	return ht
} // removeHM()

// `SetFilename()` sets `aFilename` to be used by this list.
//
// Parameters:
// - `aFilename`: The name of the file to use for storage.
//
// Returns:
// - `*THashList`: The current hash list.
func (ht *THashTags) SetFilename(aFilename string) *THashTags {
	ht.mtx.Lock()
	defer ht.mtx.Unlock()

	ht.fn = aFilename

	return ht
} // SetFilename()

// `Store()` writes the whole list to the configured file
// returning the number of bytes written and a possible error.
//
// If there is an error, it will be of type `*PathError`.
//
// Returns:
// - `int`: Number of bytes written to storage.
// - `error`: A possible storage error, or `nil` in case of success.
func (ht *THashTags) Store() (int, error) {
	ht.mtx.RLock()
	defer ht.mtx.RUnlock()

	return ht.hl.Store(ht.fn)
} // Store()

// `String()` returns the whole list as a linefeed separated string.
//
// Returns:
// - `string`: The string representation of this hash list.
func (ht *THashTags) String() string {
	ht.mtx.RLock()
	defer ht.mtx.RUnlock()

	return ht.hl.String()
} // String()

/* EoF */
