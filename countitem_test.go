/*
Copyright © 2023, 2025  M.Watermann, 10247 Berlin, Germany

	    All rights reserved
	EMail : <support@mwat.de>
*/
package hashtags

import (
	"testing"
)

//lint:file-ignore ST1017 - I prefer Yoda conditions

func Test_TCountItem_Compare(t *testing.T) {
	ci0 := TCountItem{}
	it0 := TCountItem{}

	ci1 := TCountItem{1, "#one"}
	it1 := TCountItem{1, "@one"}

	ci2 := TCountItem{1, "#one"}
	it2 := TCountItem{2, "@one"}

	ci4 := TCountItem{1, "#two"}

	tests := []struct {
		name string
		ci   TCountItem
		item TCountItem
		want int
	}{
		{"0", ci0, it0, 0},
		{"1", ci1, it1, 0},
		{"2", ci2, it2, -1},
		{"3", ci1, it2, -1},
		{"4", ci4, it2, 1},
		{"5", ci1, ci4, -1},
		{"6", it2, ci2, 1},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ci.Compare(tt.item); got != tt.want {
				t.Errorf("%q: TCountItem.Compare() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // Test_TCountItem_Compare()

func Test_TCountItem_Equal(t *testing.T) {
	ci0 := TCountItem{}
	ci1 := TCountItem{11, "one"}
	ci2 := TCountItem{222, "#two"}
	ci3 := TCountItem{222, "#alphons"}
	ci4 := ci3

	tests := []struct {
		name string
		ci   TCountItem
		item TCountItem
		want bool
	}{
		{"1", ci0, ci1, false},
		{"2", ci1, ci2, false},
		{"3", ci2, ci3, false},
		{"4", ci3, ci4, true},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ci.Equal(tt.item); got != tt.want {
				t.Errorf("%q: TCountItem.compareTo() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // Test_TCountItem_Equal()

func Test_TCountItem_Less(t *testing.T) {
	ci0 := TCountItem{}
	ci1 := TCountItem{11, "one"}
	ci2 := TCountItem{222, "#two"}
	ci3 := TCountItem{222, "#Xaver"}

	tests := []struct {
		name string
		ci   TCountItem
		item TCountItem
		want bool
	}{
		{"0", ci0, ci0, false},
		{"1", ci0, ci1, true},
		{"2", ci1, ci0, false},
		{"3", ci1, ci2, true},
		{"4", ci2, ci3, false},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ci.Less(tt.item); got != tt.want {
				t.Errorf("%q: TCountItem.Less() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // Test_TCountItem_Less()

/* EoF */
