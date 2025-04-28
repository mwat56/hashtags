/*
Copyright © 2019, 2025  M.Watermann, 10247 Berlin, Germany

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
		id   int64
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
		id   int64
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
		id   int64
		want bool
	}{
		{"0", sl0, 0, false},  // empty list
		{"1", sl1, 1, true},   // remove first element
		{"2", sl1, 5, true},   // remove last element
		{"3", sl1, 3, true},   // remove middle element
		{"4", sl1, 2, true},   // remove another element
		{"5", sl1, 2, false},  // try to remove already removed element
		{"6", sl1, 4, true},   // remove last remaining element
		{"7", sl1, 99, false}, // element not in list
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sl.remove(tt.id); got != tt.want {
				t.Errorf("%q: tSourceList.remove() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // Test_tSourceList_remove()

func Test_tSourceList_rename(t *testing.T) {
	sl1 := &tSourceList{1, 2, 3}
	sl2 := &tSourceList{}

	type tArgs struct {
		oldID, newID int64
	}
	tests := []struct {
		name string
		sl   *tSourceList
		args tArgs
		want bool
	}{
		// Empty list
		{"0", sl2, tArgs{1, 2}, false},
		// Same IDs - no change
		{"1", sl1, tArgs{1, 1}, false},
		// Replace existing ID
		{"2", sl1, tArgs{2, 4}, true},
		// Old ID doesn't exist, new ID added
		{"3", sl1, tArgs{99, 5}, true},
		// Replace first element
		{"4", sl1, tArgs{1, 6}, true},
		// Replace last element
		{"5", sl1, tArgs{3, 7}, true},
		// Nil list
		{"6", nil, tArgs{1, 2}, false},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sl.rename(tt.args.oldID, tt.args.newID); got != tt.want {
				t.Errorf("%q: tSourceList.rename() = %v, want %v",
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
	wl1 := "0000000000000001\n0000000000000002\n0000000000000003\n"

	sl2 := &tSourceList{}
	wl2 := ""
	sl3 := &tSourceList{3, 2, 1}
	wl3 := "0000000000000003\n0000000000000002\n0000000000000001\n"

	tests := []struct {
		name string
		sl   *tSourceList
		want string
	}{
		{" 1", sl1, wl1},
		{" 2", sl2, wl2},
		{" 3", sl3, wl3},

		// TODO: Add test cases.
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
