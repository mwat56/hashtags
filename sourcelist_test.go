/*
Copyright © 2019, 2024  M.Watermann, 10247 Berlin, Germany

		All rights reserved
	EMail : <support@mwat.de>
*/
package hashtags

import (
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
			if got := tt.sl.clear(); !got.compareTo(*tt.want) {
				t.Errorf("%q: tSourceList.clear() = %v\n>>>> want: >>>>\n%v",
					tt.name, got, tt.want)
			}
		})
	}
} // Test_tSourceList_clear()

func Test_tSourceList_compareTo(t *testing.T) {
	type args struct {
		aList tSourceList
	}
	tests := []struct {
		name string
		sl   tSourceList
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sl.compareTo(tt.args.aList); got != tt.want {
				t.Errorf("%q: tSourceList.compareTo() = %v\n>>>> want: >>>>\n%v",
					tt.name, got, tt.want)
			}
		})
	}
}

func Test_tSourceList_indexOf(t *testing.T) {
	sl1 := &tSourceList{
		1, 2, 3, 4, 5,
	}

	tests := []struct {
		name string
		sl   *tSourceList
		id   uint64
		want int
	}{
		{"1", sl1, 1, 0}, // first
		// {" 2", sl1, 2, 1},
		{"3", sl1, 3, 2}, // middle
		// {" 4", sl1, 4, 3},
		{"5", sl1, 5, 4},  // last
		{"6", sl1, 6, -1}, //not found
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sl.indexOf(tt.id); got != tt.want {
				t.Errorf("%q: tSourceList.indexOf() =\n %v,\n>>>> want: >>>>\n%v",
					tt.name, got, tt.want)
			}
		})
	}
} // Test_tSourceList_indexOf()

func Test_tSourceList_insert(t *testing.T) {
	sl := tSourceList{}

	wl0 := tSourceList{
		1,
	}
	wl1 := tSourceList{
		1, 3,
	}
	wl2 := tSourceList{
		1, 3, 5,
	}
	wl3 := tSourceList{
		1, 2, 3, 5,
	}
	wl4 := tSourceList{
		1, 2, 3, 4, 5,
	}

	tests := []struct {
		name string
		id   uint64
		want tSourceList
	}{
		{"0", 1, wl0},
		{"1", 3, wl1},
		{"2", 5, wl2},
		{"3", 2, wl3},
		{"4", 4, wl4},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sl.insert(tt.id); !got.compareTo(tt.want) {
				t.Errorf("%q: tSourceList.insert() =\n%v\n>>>> want: >>>>\n%v",
					tt.name, got, tt.want)
			}
		})
	}
} // Test_tSourceList_insert()

func Test_tSourceList_removeID(t *testing.T) {
	sl1 := &tSourceList{
		1, 2, 3, 4, 5,
	}
	wl1 := &tSourceList{
		2, 3, 4, 5,
	}
	wl2 := &tSourceList{
		2, 3, 4,
	}
	wl3 := &tSourceList{
		2, 4,
	}

	tests := []struct {
		name string
		sl   *tSourceList
		id   uint64
		want *tSourceList
	}{
		{" 1", sl1, 1, wl1},
		{" 2", sl1, 5, wl2},
		{" 3", sl1, 3, wl3},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sl.removeID(tt.id); !got.compareTo(*tt.want) {
				t.Errorf("%q: tSourceList.IDremove() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // Test_tSourceList_removeID()

func Test_tSourceList_renameID(t *testing.T) {
	sl := &tSourceList{
		1, 2, 3,
	}
	wl1 := tSourceList{
		1, 2, 4,
	}
	wl2 := tSourceList{
		1, 2, 6,
	}
	type args struct {
		aOldID uint64
		aNewID uint64
	}
	tests := []struct {
		name string
		ids  args
		want tSourceList
	}{
		{"1", args{3, 4}, wl1},
		{"2", args{4, 6}, wl2},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sl.renameID(tt.ids.aOldID, tt.ids.aNewID); !got.compareTo(tt.want) {
				t.Errorf("%q: tSourceList.renameID() =\n%v\n>>>> want: >>>>\n%v",
					tt.name, got, tt.want)
			}
		})
	}
} // Test_tSourceList_renameID()

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
