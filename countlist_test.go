/*
Copyright © 2023, 2024  M.Watermann, 10247 Berlin, Germany

		All rights reserved
	EMail : <support@mwat.de>
*/
package hashtags

import (
	"testing"
)

//lint:file-ignore ST1017 - I prefer Yoda conditions

func TestTCountList_compareTo(t *testing.T) {
	cl1 := TCountList{}
	wl1 := TCountList{}
	cl2 := TCountList{
		TCountItem{2, "two"},
	}
	wl2 := TCountList{
		TCountItem{2, "two"},
	}

	tests := []struct {
		name string
		cl   TCountList
		list TCountList
		want bool
	}{
		{"1", cl1, wl1, true},
		{"2", cl2, wl2, true},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cl.compareTo(tt.list); got != tt.want {
				t.Errorf("%q: TCountList.compareTo() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTCountList_compareTo()

func TestTCountList_insert(t *testing.T) {
	cl := TCountList{}
	i1 := TCountItem{1, "one"}
	wl1 := &TCountList{
		i1,
	}

	i2 := TCountItem{2, "two"}
	wl2 := &TCountList{
		i1,
		i2,
	}
	i3 := TCountItem{3, "part3"}
	wl3 := &TCountList{
		i1,
		i3,
		i2,
	}

	tests := []struct {
		name string
		item TCountItem
		want *TCountList
	}{
		{"1", i1, wl1},
		{"2", i2, wl2},
		{"3", i3, wl3},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cl.insert(tt.item); !got.compareTo(*tt.want) {
				t.Errorf("%q: TCountList.insert() = \n%v\n>>>> want: >>>\n%v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTCountList_insert()

func TestTCountList_sort(t *testing.T) {
	cl1 := &TCountList{
		TCountItem{345, "three"},
		TCountItem{234, "pure"},
		TCountItem{123, "one"},
	}

	wl1 := TCountList{
		TCountItem{123, "one"},
		TCountItem{234, "pure"},
		TCountItem{345, "three"},
	}

	tests := []struct {
		name string
		cl   *TCountList
		want TCountList
	}{
		{"1", cl1, wl1},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cl.sort(); !got.compareTo(tt.want) {
				t.Errorf("%q: TCountList.sort() =\n%v\n>>>> want: >>>\n%v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTCountList_sort()

func TestTCountList_String(t *testing.T) {
	cl := TCountList{
		TCountItem{123, "one"},
		TCountItem{234, "pure"},
		TCountItem{345, "three"},
	}
	ws := "one: 123\npure: 234\nthree: 345\n"

	tests := []struct {
		name string
		want string
	}{
		{"1", ws},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cl.String(); got != tt.want {
				t.Errorf("%q: TCountList.String() = \n%v\n>>>> want >>>>\n%v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTCountList_String()

func TestTCountList_Swap(t *testing.T) {
	c1, c2, c3 := 123, 234, 345
	cl := &TCountList{
		TCountItem{c3, "three"},
		TCountItem{c2, "pure"},
		TCountItem{c1, "one"},
	}
	wl2 := TCountList{
		TCountItem{c2, "pure"},
		TCountItem{c3, "three"},
		TCountItem{c1, "one"},
	}
	wl3 := TCountList{
		TCountItem{c2, "pure"},
		TCountItem{c1, "one"},
		TCountItem{c3, "three"},
	}

	type tArgs struct {
		i, j int
	}
	tests := []struct {
		name string
		args tArgs
		want TCountList
	}{
		{"0", tArgs{}, *cl},       // no actual change
		{"1", tArgs{c1, c2}, *cl}, // index out of bounds
		{"2", tArgs{0, 1}, wl2},   // actual swap
		{"3", tArgs{1, 2}, wl3},   // actual swap
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cl.Swap(tt.args.i, tt.args.j); !got.compareTo(tt.want) {
				t.Errorf("%q: TCountList.Swap() =\n%v\n>>>> want: >>>\n%v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTCountList_Swap()

/* EoF */
