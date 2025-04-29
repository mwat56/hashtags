/*
Copyright © 2023, 2025  M.Watermann, 10247 Berlin, Germany

	    All rights reserved
	EMail : <support@mwat.de>
*/
package hashtags

import (
	"log"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"sort"
	"strconv"
	"testing"
)

//lint:file-ignore ST1017 - I prefer Yoda conditions

var (
	exts = map[bool]string{
		true:  ".gob",
		false: ".txt",
	}
)

func hmFilename(useBinary bool) string {
	UseBinaryStorage = useBinary

	return filepath.Join(os.TempDir(), "testHmStore"+exts[useBinary])
} // hmFilename()

const baseListLen = 64

func prepHashMap() *tHashMap {
	hm := make(tHashMap, baseListLen*2)
	for i := range baseListLen {
		for j := range baseListLen {
			h, m := "#hash"+strconv.Itoa(j), "@mention"+strconv.Itoa(j)
			hm.insert(h, int64(i*11))
			hm.insert(m, int64(i*11))
		}
	}

	return &hm
} // prepHashMap()

func Test_tHashMap_checksum(t *testing.T) {
	hm1 := newHashMap()
	w1 := hm1.checksum()

	hm2 := prepHashMap()
	w2 := hm2.checksum()

	hm3 := prepHashMap()
	hm3.insert("#hash4", 444)
	w3 := hm3.checksum()

	tests1 := []struct {
		name string
		hm   *tHashMap
		want uint32
	}{
		{"1", hm1, w1},
		{"2", hm2, w2},
		{"3", hm3, w3},

		// TODO: Add test cases.
	}
	for _, tt := range tests1 {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.hm.checksum()
			if got != tt.want {
				t.Errorf("%q: tHashMap.checksum() = '%d', want '%d'",
					tt.name, got, tt.want)
			}
		})
	}

	tests2 := []struct {
		name string
		cs1  uint32
		cs2  uint32
		want bool
	}{
		{"initial vs modified", w1, w2, false},
		{"modified vs final", w2, w3, false},
		{"same checksum", w3, w3, true},

		// TODO: Add test cases.
	}
	for _, tt := range tests2 {
		t.Run(tt.name, func(t *testing.T) {
			if (tt.cs1 == tt.cs2) != tt.want {
				t.Errorf("THashTags.checksum() comparison = '%v', want '%v'",
					(tt.cs1 == tt.cs2), tt.want)
			}
		})
	}
} // Test_tHashMap_checksum()

func Test_tHashMap_clear(t *testing.T) {
	hm1 := prepHashMap()
	wm1 := newHashMap()

	tests := []struct {
		name string
		hm   *tHashMap
		want *tHashMap
	}{
		{"1", hm1, wm1},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hm.clear(); !got.equals(*tt.want) {
				t.Errorf("%q: tHashMap.clear() =\n%v\n>>>> want >>>>\n%v",
					tt.name, got, tt.want)
			}
		})
	}
} // Test_tHashMap_clear()

func Test_tHashMap_count(t *testing.T) {
	hm1 := prepHashMap()

	hm2 := prepHashMap()
	hm2.insert("#hash257", 123)

	hm3 := prepHashMap()
	hm3.insert("@mention257", 123)
	hm3.insert("@mention258", 123)

	tests := []struct {
		name    string
		hm      *tHashMap
		delim   byte
		wantInt int
	}{
		{"1", hm1, MarkHash, baseListLen},
		{"2", hm2, MarkHash, baseListLen + 1},
		{"3", hm3, MarkMention, baseListLen + 2},
		{"4", hm1, 'x', 0},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotInt := tt.hm.count(tt.delim); gotInt != tt.wantInt {
				t.Errorf("%q: tHashMap.count() = %v, want %v",
					tt.name, gotInt, tt.wantInt)
			}
		})
	}
} // Test_tHashMap_count()

func Test_tHashMap_countedList(t *testing.T) {
	// Create a small, controlled hashmap for testing
	hm1 := &tHashMap{}
	hm1.insert("#hash1", 111)
	hm1.insert("#hash2", 222)
	hm1.insert("#hash3", 333)
	hm1.insert("@mention1", 111)
	hm1.insert("@mention2", 222)

	// Empty hashmap for nil test
	hm2 := &tHashMap{}
	// var wc2 TCountList = nil

	// Hashmap with multiple IDs per hash
	hm3 := &tHashMap{}
	hm3.insert("#hash1", 111)
	hm3.insert("#hash1", 222) // Same hash, different ID
	hm3.insert("@mention1", 111)
	hm3.insert("@mention1", 333) // Same mention, different ID

	tests := []struct {
		name    string
		hm      *tHashMap
		wantInt int // Just check the length
	}{
		{"1", hm1, 5},                         // 5 entries
		{"2", hm2, 0},                         // Empty hashmap
		{"3", hm3, 2},                         // 2 entries with multiple IDs
		{"4", prepHashMap(), baseListLen * 2}, // Large hashmap from prepHashMap
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.hm.countedList()
			if len(got) != tt.wantInt {
				t.Errorf("%q: tHashMap.countedList() length = %d, want %d",
					tt.name, len(got), tt.wantInt)
			}

			// For the small hashmaps, we can do more detailed checks
			if "1" == tt.name {
				// Check that each entry has count 1
				for _, item := range got {
					if 1 != item.Count {
						t.Errorf("%q: Expected count 1 for %s, got %d",
							tt.name, item.Tag, item.Count)
					}
				}
			} else if "3" == tt.name {
				// Check that each entry has count 2
				for _, item := range got {
					if 2 != item.Count {
						t.Errorf("%q: Expected count 2 for %s, got %d",
							tt.name, item.Tag, item.Count)
					}
				}
			}
		})
	}
} // Test_tHashMap_countedList()

func Test_tHashMap_equals(t *testing.T) {
	hm1 := prepHashMap()
	om1 := prepHashMap()
	om2 := prepHashMap()
	om2.insert("#hash4", 222)

	tests := []struct {
		name string
		hm   *tHashMap
		oMap *tHashMap
		want bool
	}{
		{"1", hm1, om1, true},
		{"2", hm1, om2, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hm.equals(*tt.oMap); got != tt.want {
				t.Errorf("%q: tHashMap.equals() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // Test_tHashMap_equals()

func Test_tHashMap_idList(t *testing.T) {
	// Small controlled hashmap
	hm1 := newHashMap()
	hm1.insert("#Hash1", 111)
	hm1.insert("@Mention1", 111)
	hm1.insert("#Hash2", 222)
	hm1.insert("@Mention2", 222)

	// Large hashmap from prepHashMap
	hm2 := prepHashMap()
	// Add a unique ID with multiple tags
	hm2.insert("#uniqueHash", 999999)
	hm2.insert("@uniqueMention", 999999)

	tests := []struct {
		name         string
		hm           *tHashMap
		id           int64
		wantLen      int
		wantContains []string
	}{
		{"empty", newHashMap(), 111, 0, nil},
		{"single ID", hm1, 111, 2, []string{"#hash1", "@mention1"}}, // lowercase!
		{"not found", hm1, 999, 0, nil},
		{"unique ID in large map", hm2, 999999, 2, []string{"#uniquehash", "@uniquemention"}}, // lowercase!
		{"common ID in large map", hm2, 0, baseListLen * 2, nil},                              // All entries in prepHashMap have ID 0
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.hm.idList(tt.id)

			if len(got) != tt.wantLen {
				t.Errorf("%q: tHashMap.idList() returned %d items, want %d",
					tt.name, len(got), tt.wantLen)
			}

			if nil != tt.wantContains {
				for _, tag := range tt.wantContains {
					if !slices.Contains(got, tag) {
						t.Errorf("%q: tHashMap.idList() missing expected tag %q",
							tt.name, tag)
					}
				}
			}

			// Check that result is sorted
			if 1 < len(got) {
				sorted := slices.Clone(got)
				sort.Strings(sorted)
				if !slices.Equal(got, sorted) {
					t.Errorf("%q: tHashMap.idList() returned unsorted result %v",
						tt.name, got)
				}
			}
		})
	}
} // Test_tHashMap_idList()

func Test_tHashMap_idxLen(t *testing.T) {
	hm := prepHashMap()
	hm.insert("#hash2", 333)

	type tArgs struct {
		aDelim byte
		aName  string
	}
	tests := []struct {
		name string
		args tArgs
		want int
	}{
		{"0", tArgs{}, -1},
		{"1", tArgs{MarkHash, "#hash2"}, baseListLen + 1},
		{"2", tArgs{MarkHash, "hash2"}, baseListLen + 1},
		{"3", tArgs{MarkMention, "@hash2"}, -1},
		{"4", tArgs{MarkMention, "mention1"}, baseListLen},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hm.idxLen(tt.args.aDelim, tt.args.aName); got != tt.want {
				t.Errorf("%q: tHashMap.idxLen() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // Test_tHashMap_idxLen()

func Test_tHashMap_insert(t *testing.T) {
	hm1 := prepHashMap()
	hm2 := prepHashMap()
	hm2.insert("#hash4", 444)

	type tArgs struct {
		aName string
		aID   int64
	}
	tests := []struct {
		name string
		hm   *tHashMap
		args tArgs
		want bool
	}{
		// {" 0", hm1, tArgs{"", 0}, false},            // empty hash
		{" 1", hm1, tArgs{"@mention2", 220}, false}, // already added
		// {" 2", hm2, tArgs{"#hash4", 222}, true},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hm.insert(tt.args.aName, tt.args.aID); got != tt.want {
				t.Errorf("%q: tHashMap.insert() =\n%v\n>>>> want: >>>>\n%v",
					tt.name, got, tt.want)
			}
		})
	}
} // Test_tHashMap_insert()

func Test_tHashMap_keys(t *testing.T) {
	hm0 := newHashMap()

	hm1 := newHashMap()
	hm1.insert("#hash1", 111)
	hm1.insert("#hash2", 222)
	hm1.insert("@mention1", 333)

	tests := []struct {
		name     string
		hm       *tHashMap
		wantKeys []string
		wantLen  int
	}{
		{"empty", hm0, []string{}, 0},
		{"small", hm1, []string{"#hash1", "#hash2", "@mention1"}, 3},
		{"large", prepHashMap(), nil, baseListLen * 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.hm.keys()

			// Check length for all cases
			if len(got) != tt.wantLen {
				t.Errorf("%q: tHashMap.keys() returned %d keys, want %d",
					tt.name, len(got), tt.wantLen)
			}

			// For small maps, check exact content
			if nil != tt.wantKeys && !reflect.DeepEqual(got, tt.wantKeys) {
				t.Errorf("%q: tHashMap.keys() = %v, want %v",
					tt.name, got, tt.wantKeys)
			}

			// For large maps, check sorting
			if "large" == tt.name && 0 < len(got) {
				// Verify keys are sorted
				sorted := make([]string, len(got))
				copy(sorted, got)
				slices.SortFunc(sorted, cmp4sort)

				if !reflect.DeepEqual(got, sorted) {
					t.Errorf("%q: tHashMap.keys() returned unsorted keys", tt.name)
				}
			}
		})
	}
} // Test_tHashMap_keys()

func Test_tHashMap_list(t *testing.T) {
	// Create a small, controlled hashmap for testing
	hm1 := newHashMap()
	hm1.insert("#hash1", 111)
	hm1.insert("#hash2", 222)
	hm1.insert("#hash3", 333)
	hm1.insert("#hash2", 333) // Add another ID to hash2
	hm1.insert("@mention1", 111)
	hm1.insert("@mention2", 222)

	// Test with the large prepHashMap
	hm2 := prepHashMap()

	wl0 := []int64{}
	wl1 := []int64{111}
	wl2 := []int64{222, 333}

	wl6 := []int64{0, 11, 22, 33, 44, 55, 66, 77, 88, 99, 110, 121, 132, 143, 154, 165, 176, 187, 198, 209, 220, 231, 242, 253, 264, 275, 286, 297, 308, 319, 330, 341, 352, 363, 374, 385, 396, 407, 418, 429, 440, 451, 462, 473, 484, 495, 506, 517, 528, 539, 550, 561, 572, 583, 594, 605, 616, 627, 638, 649, 660, 671, 682, 693}
	wl7 := []int64{0, 11, 22, 33, 44, 55, 66, 77, 88, 99, 110, 121, 132, 143, 154, 165, 176, 187, 198, 209, 220, 231, 242, 253, 264, 275, 286, 297, 308, 319, 330, 341, 352, 363, 374, 385, 396, 407, 418, 429, 440, 451, 462, 473, 484, 495, 506, 517, 528, 539, 550, 561, 572, 583, 594, 605, 616, 627, 638, 649, 660, 671, 682, 693}

	type tArgs struct {
		aDelim byte
		aName  string
	}
	tests := []struct {
		name     string
		hm       *tHashMap
		args     tArgs
		wantList []int64
	}{
		// Tests with small controlled hashmap
		{"0", hm1, tArgs{}, wl0},
		{"1", hm1, tArgs{MarkHash, " "}, wl0},
		{"2", hm1, tArgs{MarkHash, "hash1"}, wl1},
		{"3", hm1, tArgs{MarkHash, "hash2"}, wl2},
		{"4", hm1, tArgs{MarkMention, "mention1"}, wl1},
		{"5", hm1, tArgs{MarkMention, "hash1"}, wl0},

		// Tests with large hashmap - just check length
		{"6", hm2, tArgs{MarkHash, "hash0"}, wl6},
		{"7", hm2, tArgs{MarkMention, "mention0"}, wl7},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotList := tt.hm.list(tt.args.aDelim, tt.args.aName)

			if 1 == len(tt.name) && tt.name[0] >= '6' {
				// For large hashmap tests, just check the length
				if len(gotList) != len(tt.wantList) {
					t.Errorf("%q: tHashMap.list() returned %d items, want %d",
						tt.name, len(gotList), len(tt.wantList))
				}
			} else {
				// For small hashmap tests, check exact content
				if !slices.Equal(gotList, tt.wantList) {
					t.Errorf("%q: tHashMap.list() =\n%v\n>>>> want: >>>>\n%v",
						tt.name, gotList, tt.wantList)
				}
				// if !reflect.DeepEqual(gotList, tt.wantList) {
				// 	t.Errorf("%q: tHashMap.list() =\n%v\n>>>> want: >>>>\n%v",
				// 		tt.name, gotList, tt.wantList)
				// }
			}
		})
	}
} // Test_tHashMap_list()

func Test_tHashMap_load(t *testing.T) {
	saveBinary := UseBinaryStorage
	defer func() {
		UseBinaryStorage = saveBinary
	}()

	hm1 := prepHashMap()
	hm1.insert("@CrashTestDummy", 1)

	wm1 := prepHashMap()
	wm1.insert("@CrashTestDummy", 1)

	tests := []struct {
		name    string
		hm      *tHashMap
		binary  bool
		want    *tHashMap
		wantErr bool
	}{
		{"1", hm1, false, wm1, false},
		{"2", hm1, true, wm1, false},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn := hmFilename(tt.binary)
			// make sure there's actually data in the file:
			tt.hm.store(fn)

			got, err := tt.hm.load(fn)
			if (nil != err) != tt.wantErr {
				t.Errorf("%q: tHashMap.load() =\n%v\n>>>> want >>>>\n'%v'",
					tt.name, err, tt.wantErr)
				return
			}
			if !tt.want.equals(*got) {
				t.Errorf("%q: tHashMap.load() =\n%v\n>>>> want >>>>\n%v",
					tt.name, got, tt.want)
			}
			// fName = hmFilename(!tt.binary)
			// tt.hm.store(fName)

			// os.Remove(fn)
		})
	}
} // Test_tHashMap_load()

func Test_tHashMap_remove(t *testing.T) {
	hm := prepHashMap()
	hm.insert("#nameX", 999)
	hm.insert("#hash3", 333)

	type tArgs struct {
		aDelim byte
		aName  string
		aID    int64
	}
	tests := []struct {
		name string
		args tArgs
		want bool
	}{
		{"0", tArgs{}, false},
		{"1", tArgs{MarkHash, "#nameX", 111}, false},
		{"2", tArgs{MarkHash, "#nameX", 999}, true},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hm.removeHM(tt.args.aDelim, tt.args.aName, tt.args.aID)
			if got != tt.want {
				t.Errorf("%q: tHashMap.remove() =\n%v\n>>>> want >>>>\n%v",
					tt.name, got, tt.want)
			}
		})
	}
} // Test_tHashMap_remove()

func Test_tHashMap_removeHM(t *testing.T) {
	hm := prepHashMap()
	hm.insert("#testTag", 123)
	hm.insert("@testMention", 456)

	type tArgs struct {
		aDelim byte
		aTag   string
		aID    int64
	}
	tests := []struct {
		name string
		hm   *tHashMap
		args tArgs
		want bool
	}{
		{"empty tag", hm, tArgs{MarkHash, "", 123}, false},
		{"whitespace tag", hm, tArgs{MarkHash, "  ", 123}, false},
		{"non-existent tag", hm, tArgs{MarkHash, "nonexistent", 123}, false},
		{"existing tag wrong ID", hm, tArgs{MarkHash, "testTag", 999}, false},
		{"existing tag with ID", hm, tArgs{MarkHash, "testTag", 123}, true},
		{"existing tag with prefix", hm, tArgs{MarkHash, "#testTag", 123}, false},
		{"mention with ID", hm, tArgs{MarkMention, "testMention", 456}, true},
		{"mention with prefix", hm, tArgs{MarkMention, "@testMention", 456}, false},

		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hm.removeHM(tt.args.aDelim, tt.args.aTag, tt.args.aID); got != tt.want {
				t.Errorf("%q: tHashMap.removeHM() = '%v', want '%v'",
					tt.name, got, tt.want)
			}
		})
	}
} // Test_tHashMap_removeHM()

func Test_tHashMap_removeID(t *testing.T) {
	hm0 := newHashMap()

	hm1 := prepHashMap()
	hm1.insert("#hash2", 999)
	hm1.insert("#hash3", 888)
	hm1.insert("#CrashTestDummy", 777)

	tests := []struct {
		name string
		hm   *tHashMap
		id   int64
		want bool
	}{
		{"0", hm0, 999, false},
		{"1", hm1, 999, true},
		{"2", hm1, 888, true},
		{"3", hm1, 777, true},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hm.removeID(tt.id); got != tt.want {
				t.Errorf("%q: tHashMap.removeID() =\n%v\n>>>> want >>>>\n%v",
					tt.name, got, tt.want)
			}
		})
	}
} // Test_tHashMap_removeID()

func Test_tHashMap_renameID(t *testing.T) {
	id1, id2, id3 := int64(11), int64(22), int64(33)

	hm1 := prepHashMap()
	hm1.insert("#hash1", id1)
	hm1.insert("#hash1", id3)
	hm1.insert("@mention1", id3)

	hm3 := prepHashMap()
	hm3.insert("#hash3", id2)
	hm3.insert("#hash4", id2)

	type tArgs struct {
		aOldID, aNewID int64
	}
	tests := []struct {
		name string
		hm   *tHashMap
		args tArgs
		want bool
	}{
		{"0", hm1, tArgs{}, false}, // no change
		{"1", hm1, tArgs{id1, id2}, true},
		{"2", hm3, tArgs{id2, id3}, true},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hm.renameID(tt.args.aOldID, tt.args.aNewID); got != tt.want {
				t.Errorf("%q: tHashMap.renameID() =\n%v\n>>>> want: >>>>\n%v",
					tt.name, got, tt.want)
			}
		})
	}
} // Test_tHashMap_renameID

func Test_tHashMap_sort(t *testing.T) {
	hm1 := &tHashMap{
		"#hash1": &tSourceList{
			int64(111),
		},
		"@mention1": &tSourceList{
			int64(111),
		},
		"#hash2": &tSourceList{
			int64(222),
		},
		"@mention2": &tSourceList{
			int64(333),
			int64(222),
		},
		"#hash3": &tSourceList{
			int64(333),
		},
	}
	wm1 := &tHashMap{
		"#hash1": &tSourceList{
			int64(111),
		},
		"#hash2": &tSourceList{
			int64(222),
		},
		"#hash3": &tSourceList{
			int64(333),
		},
		"@mention1": &tSourceList{
			int64(111),
		},
		"@mention2": &tSourceList{
			int64(222),
			int64(333),
		},
	}

	tests := []struct {
		name string
		hm   *tHashMap
		want *tHashMap
	}{
		{"1", hm1, wm1},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hm.sort(); !got.equals(*tt.want) {
				t.Errorf("%q: tHashMap.sort() = \n%v\n>>>> want: >>>>\n%v",
					tt.name, got, tt.want)
			}
		})
	}
} // Test_tHashMap_sort()

func Test_tHashMap_store(t *testing.T) {
	saveBinary := UseBinaryStorage
	defer func() {
		UseBinaryStorage = saveBinary
	}()

	hm1 := prepHashMap()
	hm1.insert("@alphons", 1)

	tests := []struct {
		name    string
		hm      *tHashMap
		binary  bool
		wantInt int
		wantErr bool
	}{
		{"1", hm1, false, 140744, false}, // expected file size
		{"2", hm1, true, 23653, false},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		fName := hmFilename(tt.binary)
		t.Run(tt.name, func(t *testing.T) {
			gotInt, err := tt.hm.store(fName)
			if (nil != err) != tt.wantErr {
				t.Errorf("%q: tHashMap.store() error = '%v', wantErr '%v'",
					tt.name, err, tt.wantErr)
				return
			}
			if gotInt != tt.wantInt {
				t.Errorf("%q: tHashMap.store() = '%d', want '%d'",
					tt.name, gotInt, tt.wantInt)
			}
		})
		// os.Remove(fName)
	}
} // Test_tHashMap_store()

func Test_tHashMap_String(t *testing.T) {
	// Empty hashmap
	hm1 := newHashMap()
	want1 := ""

	// Hashmap with single entry
	hm2 := newHashMap()
	hm2.insert("#test", 123)

	// Hashmap with multiple entries from prepHashMap()
	hm3 := prepHashMap()

	tests := []struct {
		name string
		hm   *tHashMap
		want string
	}{
		{"empty", hm1, want1},
		{"single entry", hm2, "[#test]\n000000000000007b\n"},
		{"large map", hm3, ""}, // We'll check length instead of exact content
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.hm.String()

			if "empty" == tt.name || "single entry" == tt.name {
				if tt.want != got {
					t.Errorf("%q: tHashMap.String() =\n%q\nwant:\n%q",
						tt.name, got, tt.want)
				}
			} else if "large map" == tt.name {
				// For the large map from prepHashMap(), just check that:
				// 1. The output is not empty
				// 2. The output is reasonably large (at least 1000 chars)
				if "" == got {
					t.Errorf("%q: tHashMap.String() returned empty string",
						tt.name)
				}

				// Calculate approximate length by the
				// number of tags * list length *
				// length of each ID  as hex string:
				approx := len(*tt.hm) * baseListLen * 16
				if approx > len(got) {
					t.Errorf("%q: tHashMap.String() returned string of length %d, expected at least %d",
						tt.name, len(got), approx)
				}
			}
		})
	}
} // Test_tHashMap_String()

func Benchmark_LoadTxT(b *testing.B) {
	saveBinary := UseBinaryStorage
	defer func() {
		UseBinaryStorage = saveBinary
	}()
	fn := hmFilename(false)

	hm := prepHashMap()
	hm.insert("@CrashTestDummy", 1)
	hm.store(fn)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		if _, err := hm.load(fn); nil != err {
			log.Printf("LoadTxt(): %v", err)
		}
	}

	// os.Remove(fn)
} // Benchmark_LoadTxt()

func Benchmark_LoadGob(b *testing.B) {
	saveBinary := UseBinaryStorage
	defer func() {
		UseBinaryStorage = saveBinary
	}()
	fn := hmFilename(true)

	hm := prepHashMap()
	hm.insert("@CrashTestDummy", 1)
	hm.store(fn)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		if _, err := hm.load(fn); nil != err {
			log.Printf("LoadBin(): %v", err)
		}
	}

	// os.Remove(fn)
} // Benchmark_LoadBin()

/*
func Benchmark_LoadCustom(b *testing.B) {
	saveBinary := UseBinaryStorage
	defer func() {
		UseBinaryStorage = saveBinary
	}()
	fn := hmFilename(true) + `.custom`

	hm := prepHashMap()
	hm.insert("@CrashTestDummy", 1)
	hm.storeCustomBinary(fn)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		if _, err := loadCustomBinary(fn); nil != err {
			log.Printf("LoadBin(): %v", err)
		}
	}

	// os.Remove(fn)
} // Benchmark_LoadBin()
*/

func Benchmark_StoreTxt(b *testing.B) {
	saveBinary := UseBinaryStorage
	defer func() {
		UseBinaryStorage = saveBinary
	}()
	fn := hmFilename(false)

	hm := prepHashMap()
	hm.insert("@CrashTestDummy", 1)

	for n := 0; n < b.N; n++ {
		if _, err := hm.store(fn); nil != err {
			log.Printf("StoreTxt(): %v", err)
		}
	}

	// os.Remove(fn)
} // Benchmark_StoreTxt()

func Benchmark_StoreGob(b *testing.B) {
	saveBinary := UseBinaryStorage
	defer func() {
		UseBinaryStorage = saveBinary
	}()
	fn := hmFilename(true)

	hm := prepHashMap()
	hm.insert("@CrashTestDummy", 1)

	for n := 0; n < b.N; n++ {
		if _, err := hm.store(fn); nil != err {
			log.Printf("StoreBin(): %v", err)
		}
	}

	// os.Remove(fn)
} // Benchmark_StoreBin()

/*
func Benchmark_StoreCustom(b *testing.B) {
	saveBinary := UseBinaryStorage
	defer func() {
		UseBinaryStorage = saveBinary
	}()
	fn := hmFilename(true) + `.custom`

	hm := prepHashMap()
	hm.insert("@CrashTestDummy", 1)

	for n := 0; n < b.N; n++ {
		if _, err := hm.storeCustomBinary(fn); nil != err {
			log.Printf("StoreBin(): %v", err)
		}
	}

	// os.Remove(fn)
} // Benchmark_StoreCustom()
*/

/* EoF */
