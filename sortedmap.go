/*
Copyright © 2023, 2024  M.Watermann, 10247 Berlin, Germany

		All rights reserved
	EMail : <support@mwat.de>
*/
package hashtags

import (
	"cmp"
	"fmt"
	"sort"
	"sync"
)

//lint:file-ignore ST1017 - I prefer Yoda conditions

// `TSortedMap` is a generic type that accepts two type parameters:
// - K for the key type (which must be cmp.Ordered)
// - V for the value type
//
// The `cmp.Ordered` interface allows all numeric data types as well
// as pointers and strings.
type TSortedMap[K cmp.Ordered, V any] struct {
	m    map[K]V
	keys []K
	mtx  sync.RWMutex
}

// --------------------------------------------------------------------------
// constructor function

// `NewSortedMap()` creates a new instance of `TSortedMap` with the
// specified key and value types.
//
// The returned map is initially empty.
//
// Parameters:
// - `K`: The type of the keys in the sorted map.
// - `V`: The type of the values in the sorted map.
//
// Returns:
// - `*TSortedMap[K, V]`: A pointer to a new instance with the specified
// key and value types.
func NewSortedMap[K cmp.Ordered, V any]() *TSortedMap[K, V] {
	return &TSortedMap[K, V]{
		m:    make(map[K]V),
		keys: make([]K, 0),
	}
} // NewSortedMap()

// --------------------------------------------------------------------------
// methods of TSortedMap

// `Add()` adds or updates a key/value pair in the sorted map.
//
// Parameters:
// - `aKey`: The key of the entry to be added or updated.
// - `aValue`: The value to be associated with the key.
//
// Returns:
// - `*TSortedMap[K, V]`: A pointer to the updated SortedMap instance.
func (sm *TSortedMap[K, V]) Add(aKey K, aValue V) *TSortedMap[K, V] {
	sm.mtx.Lock()
	defer sm.mtx.Unlock()

	if _, exists := sm.m[aKey]; !exists {
		// There are different situations to consider:
		// 1: the key-list is empty,
		// 2: the key-list doesn't already contain the key,
		// 3: the key-list contains the key but with a different value
		sLen := len(sm.keys)
		if 0 == sLen {
			// 1: empty list: just add the new item
			sm.keys = append(sm.keys, aKey)
		} else {
			// find the insertion index using binary search
			idx := sort.Search(sLen, func(i int) bool {
				return sm.keys[i] >= aKey
			})

			if sLen == idx {
				// 2: key not found: add key at the end
				sm.keys = append(sm.keys, aKey)
			} else if (sm.keys)[idx] != aKey {
				// 3: the search index doesn't point to the required key
				sm.keys = append(sm.keys, aKey)
				copy((sm.keys)[idx+1:], (sm.keys)[idx:])
				(sm.keys)[idx] = aKey
			} else {
				// dummy instruction for debugger
				fmt.Println("\n", sLen)
			}
		}
	}

	sm.m[aKey] = aValue
	return sm
} // Add()

// `Delete()` removes a key/value pair from the map.
//
// Parameters:
// - `aKey`: The key of the entry to be deleted.
//
// Returns:
// - `*TSortedMap[K, V]`: A pointer to the updated map instance.
func (sm *TSortedMap[K, V]) Delete(aKey K) *TSortedMap[K, V] {
	sm.mtx.Lock()
	defer sm.mtx.Unlock()

	if _, exists := sm.m[aKey]; exists {
		delete(sm.m, aKey)
		for idx, key := range sm.keys {
			if key == aKey {
				sm.keys = append(sm.keys[:idx], sm.keys[idx+1:]...)
				break
			}
		}
	}

	return sm
} // Delete()

// `Get()` retrieves a value by its key from the SortedMap
//
// Parameters:
// - `aKey`: The key of the entry to be retrieved.
//
// Returns:
// - `V`: The value associated with the `aKey`.
// - `bool`: An indication whether the key was found in the map.
func (sm *TSortedMap[K, V]) Get(aKey K) (V, bool) {
	sm.mtx.RLock()
	defer sm.mtx.RUnlock()

	value, exists := sm.m[aKey]

	return value, exists
} // Get()

// Keys returns a slice of all keys in sorted order

// `Keys()` returns a slice of all keys in sorted order
//
// This method returns a slice of all keys in the map in sorted order.
// It is thread-safe and can be called concurrently.
//
// Returns:
// - `[]K`: A slice of keys in the sorted map.
func (sm *TSortedMap[K, V]) Keys() []K {
	sm.mtx.RLock()
	defer sm.mtx.RUnlock()

	return append([]K{}, sm.keys...)
} // Keys()

// `Iterate()` allows iteration over the map in sorted key order.
//
// Parameters:
// - `f`: A function that takes a key and its associated value as arguments and performs some operation on them.
//
// Returns:
// - `*TSortedMap[K, V]`: A pointer to the same SortedMap instance, allowing method chaining.
func (sm *TSortedMap[K, V]) Iterate(f func(aKey K, aValue V)) *TSortedMap[K, V] {
	sm.mtx.RLock()
	defer sm.mtx.RUnlock()

	for _, key := range sm.keys {
		f(key, sm.m[key])
	}

	return sm
} // Iterate()

/* EoF */
