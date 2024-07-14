/*
Copyright © 2019, 2024  M.Watermann, 10247 Berlin, Germany

		All rights reserved
	EMail : <support@mwat.de>
*/
package hashtags

import (
	"reflect"
	"testing"
)

//lint:file-ignore ST1017 - I prefer Yoda conditions

func Test_tSourceList_clear(t *testing.T) {
	sl1 := &tSourceList{
		1,
		2,
		3,
	}
	wl1 := &tSourceList{}

	tests := []struct {
		name string
		sl   *tSourceList
		want *tSourceList
	}{
		{"1", sl1, wl1},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sl.clear(); !got.equals(*tt.want) {
				t.Errorf("%q: tSourceList.clear() = %v\n>>>> want: >>>>\n%v",
					tt.name, got, tt.want)
			}
		})
	}
} // Test_tSourceList_clear()

func Test_tSourceList_equals(t *testing.T) {
	sl1 := tSourceList{
		1,
		2,
		3,
	}
	sl2 := tSourceList{
		3,
		2,
		1,
	}

	tests := []struct {
		name string
		sl   tSourceList
		list tSourceList
		want bool
	}{
		{"0", sl1, nil, false},
		{"1", sl1, sl1, true},
		{"2", sl2, sl1, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sl.equals(tt.list); got != tt.want {
				t.Errorf("%q: tSourceList.equals() = %v\n>>>> want: >>>>\n%v",
					tt.name, got, tt.want)
			}
		})
	}
} // Test_tSourceList_equals()

func Test_tSourceList_findIndex(t *testing.T) {
	sl1 := &tSourceList{
		1, 2, 3, 4, 5,
	}

	tests := []struct {
		name string
		sl   *tSourceList
		id   uint64
		want int
	}{
		{"0", &tSourceList{}, 1, -1}, // empty list
		{"1", sl1, 1, 0},             // first
		{"2", sl1, 3, 2},             // middle
		{"3", sl1, 5, 4},             // last
		{"4", sl1, 6, -1},            // not found
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sl.findIndex(tt.id); got != tt.want {
				t.Errorf("%q: tSourceList.findIndex() =\n %v,\n>>>> want: >>>>\n%v",
					tt.name, got, tt.want)
			}
		})
	}
} // Test_tSourceList_findIndex()

func Test_tSourceList_insert(t *testing.T) {
	sl := tSourceList{}

	tests := []struct {
		name string
		id   uint64
		want bool
	}{
		{"0", 1, true}, // beginning
		{"1", 3, true}, // end
		{"2", 5, true}, // end
		{"3", 2, true}, // middle
		{"4", 4, true}, // middle
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sl.insert(tt.id); got != tt.want {
				t.Errorf("tSourceList.insert() = %v, want %v", got, tt.want)
			}
		})
	}
} // Test_tSourceList_insert()

func Test_tSourceList_remove(t *testing.T) {
	sl0 := &tSourceList{}
	sl1 := &tSourceList{1, 2, 3, 4, 5}

	tests := []struct {
		name string
		sl   *tSourceList
		id   uint64
		want bool
	}{
		{"0", sl0, 0, false}, // not found
		{"1", sl1, 1, true},  // beginning
		{"2", sl1, 5, true},  // end
		{"3", sl1, 3, true},  // middle
		{"4", sl1, 2, true},  // (new) beginning
		{"5", sl1, 2, false}, // not found
		{"6", sl1, 4, true},  // beginning == end
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sl.remove(tt.id); got != tt.want {
				t.Errorf("%q: tSourceList.remove() =\n%v\n>>>> want: >>>>\n%v",
					tt.name, got, tt.want)
			}
		})
	}
} // Test_tSourceList_remove()

func Test_tSourceList_rename(t *testing.T) {
	sl := &tSourceList{1, 2, 3}

	type tArgs struct {
		aOldID, aNewID uint64
	}
	tests := []struct {
		name string
		ids  tArgs
		want bool
	}{
		{"0", tArgs{}, false},         // not found
		{"1", tArgs{3, 4}, true},      // end
		{"2", tArgs{4, 6}, true},      // (new) end
		{"3", tArgs{3333, 333}, true}, // only new ID added
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sl.rename(tt.ids.aOldID, tt.ids.aNewID); got != tt.want {
				t.Errorf("%q: tSourceList.rename() =\n%v\n>>>> want: >>>>\n%v",
					tt.name, got, tt.want)
			}
		})
	}
} // Test_tSourceList_rename()

func Test_tSourceList_sort(t *testing.T) {
	sl1 := &tSourceList{}
	wl1 := &tSourceList{}
	sl2 := &tSourceList{
		3, 1, 2,
	}
	wl2 := &tSourceList{
		1, 2, 3,
	}
	tests := []struct {
		name string
		sl   *tSourceList
		want *tSourceList
	}{
		{"0", nil, nil},
		{"1", sl1, wl1},
		{"2", sl2, wl2},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sl.sort(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%q: tSourceList.sort() =\n%v\n>>>> want: >>>>\n%v",
					tt.name, got, tt.want)
			}
		})
	}
} // Test_tSourceList_sort()

func Test_tSourceList_String(t *testing.T) {
	sl1 := &tSourceList{
		1,
		2,
		3,
	}
	wl1 := "1\n2\n3\n"

	sl2 := &tSourceList{}
	wl2 := ""

	tests := []struct {
		name string
		sl   *tSourceList
		want string
	}{
		// TODO: Add test cases.
		{" 1", sl1, wl1},
		{" 2", sl2, wl2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sl.String(); got != tt.want {
				t.Errorf("%q: tSourceList.String() =\n%v\n>>>> want: >>>>\n%v",
					tt.name, got, tt.want)
			}
		})
	}
} // Test_tSourceList_String()

/* EoF */
