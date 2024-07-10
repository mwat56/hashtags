/*
Copyright © 2023, 2024  M.Watermann, 10247 Berlin, Germany

		All rights reserved
	EMail : <support@mwat.de>
*/
package hashtags

import "testing"

//lint:file-ignore ST1017 - I prefer Yoda conditions

func TestTCountItem_compareTo(t *testing.T) {
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
			if got := tt.ci.compareTo(tt.item); got != tt.want {
				t.Errorf("%q: TCountItem.compareTo() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTCountItem_compareTo()

func TestTCountItem_Less(t *testing.T) {
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
} // TestTCountItem_Less()

/* EoF */
