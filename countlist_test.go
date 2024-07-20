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

func TestTCountList_Compare(t *testing.T) {
	cl0 := TCountList{}
	wl0 := TCountList{}

	cl2 := TCountList{TCountItem{2, "#two"}}
	wl2 := TCountList{TCountItem{2, "@two"}}

	wl4 := TCountList{TCountItem{1, "@two"}}
	wl5 := TCountList{TCountItem{1, "zero"}}

	tests := []struct {
		name string
		cl   TCountList
		ol   TCountList
		want int
	}{
		{"0", cl0, wl0, 0},
		{"1", cl0, wl2, -1},
		{"2", cl2, wl0, 1},
		{"3", cl2, wl2, 0},
		{"4", cl2, wl4, 1},
		{"2", wl4, wl5, -1},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cl.Compare(tt.ol); got != tt.want {
				t.Errorf("%q: TCountList.Compare() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTCountList_Compare()

func TestTCountList_Equal(t *testing.T) {
	cl1 := TCountList{}
	wl1 := TCountList{}

	cl2 := TCountList{TCountItem{2, "two"}}
	wl2 := TCountList{TCountItem{2, "two"}}
	wl3 := wl2

	cl4 := TCountList{
		TCountItem{1, "one"}, TCountItem{2, "two"}}
	wl4 := TCountList{TCountItem{2, "two"}, TCountItem{1, "one"}}

	cl5 := cl4
	wl5 := TCountList{TCountItem{11, "one"}, TCountItem{22, "two"}}

	tests := []struct {
		name string
		cl   TCountList
		list TCountList
		want bool
	}{
		{"1", cl1, wl1, true},
		{"2", cl2, wl2, true},
		{"3", cl1, wl3, false},
		{"4", cl4, wl4, false},
		{"5", cl5, wl5, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cl.Equal(tt.list); got != tt.want {
				t.Errorf("%q: TCountList.compareTo() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTCountList_Equal()

func TestTCountList_Insert(t *testing.T) {
	cl := TCountList{}
	i1 := TCountItem{1, "one"}
	wl1 := &TCountList{i1}

	i2 := TCountItem{2, "two"}
	wl2 := &TCountList{i1, i2}

	i3 := TCountItem{3, "part3"}
	wl3 := &TCountList{i1, i3, i2}

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
			if got := cl.Insert(tt.item); !got.Equal(*tt.want) {
				t.Errorf("%q: TCountList.Insert() = \n%v\n>>>> want: >>>\n%v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTCountList_Insert()

func TestTCountList_Len(t *testing.T) {
	cl0 := TCountList{}
	cl1 := TCountList{TCountItem{1, "one"}}

	tests := []struct {
		name string
		cl   TCountList
		want int
	}{
		{"0", cl0, 0},
		{"1", cl1, 1},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cl.Len(); got != tt.want {
				t.Errorf("%q: TCountList.len() = %d, want %d",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTCountList_Len()

func TestTCountList_sort(t *testing.T) {
	cl0 := &TCountList{}
	cl1 := &TCountList{
		TCountItem{345, "three"},
		TCountItem{234, "@pure"},
		TCountItem{123, "#one"},
	}
	wl1 := TCountList{
		TCountItem{123, "#one"},
		TCountItem{234, "@pure"},
		TCountItem{345, "three"},
	}
	cl4 := &TCountList{
		TCountItem{123, "#one"},
		TCountItem{234, "@pure"},
		TCountItem{345, "three"},
		TCountItem{678, "#one"},
	}
	wl4 := TCountList{
		TCountItem{123, "#one"},
		TCountItem{678, "#one"},
		TCountItem{234, "@pure"},
		TCountItem{345, "three"},
	}
	cl5 := &TCountList{
		TCountItem{123, "#one"},
		TCountItem{234, "@pure"},
		TCountItem{234, "#one"},
		TCountItem{345, "three"},
		TCountItem{345, "#one"},
	}
	wl5 := TCountList{
		TCountItem{123, "#one"},
		TCountItem{234, "#one"},
		TCountItem{345, "#one"},
		TCountItem{234, "@pure"},
		TCountItem{345, "three"},
	}
	cl6 := &TCountList{
		TCountItem{987, "#one"},
		TCountItem{234, "@pure"},
		TCountItem{654, "#one"},
		TCountItem{345, "three"},
		TCountItem{321, "#one"},
	}
	wl6 := TCountList{
		TCountItem{321, "#one"},
		TCountItem{654, "#one"},
		TCountItem{987, "#one"},
		TCountItem{234, "@pure"},
		TCountItem{345, "three"},
	}
	cl7 := &TCountList{
		TCountItem{987, "#one"},
		TCountItem{235, "two"},
		TCountItem{654, "#one"},
		TCountItem{235, "two"},
		TCountItem{321, "#one"},
	}
	wl7 := TCountList{
		TCountItem{321, "#one"},
		TCountItem{654, "#one"},
		TCountItem{987, "#one"},
		TCountItem{235, "two"},
		TCountItem{235, "two"},
	}

	tests := []struct {
		name string
		cl   *TCountList
		want TCountList
	}{
		{"0", cl0, TCountList{}},
		{"1", cl1, wl1},
		{"2", cl1, *cl1},
		{"3", &wl1, *cl1},
		{"4", cl4, wl4},
		{"5", cl5, wl5},
		{"6", cl6, wl6},
		{"7", cl7, wl7},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cl.sort(); !got.Equal(tt.want) {
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
			if got := cl.Swap(tt.args.i, tt.args.j); !got.Equal(tt.want) {
				t.Errorf("%q: TCountList.Swap() =\n%v\n>>>> want: >>>\n%v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTCountList_Swap()

/* EoF */
