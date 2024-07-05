/*
Copyright © 2023, 2024  M.Watermann, 10247 Berlin, Germany

		All rights reserved
	EMail : <support@mwat.de>
*/
package hashtags

//lint:file-ignore ST1017 - I prefer Yoda conditions

/* * /
// does not compile: tt.ci.compareTo undefined (type struct{Count int; Tag string} has no field or method compareTo)
func TestTCountItem_compareTo(t *testing.T) {
	ci1 := TCountItem{11, "one"}
	ci2 := TCountItem{222, "two"}

	tests := []struct {
		name string
		ci   TCountItem
		item TCountItem
		want int
	}{
		{"1", ci1, ci2, -1},
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
/* */

/* * /
// does not compile: tt.ci.Less undefined (type struct{Count int; Tag string} has no field or method Less)
func TestTCountItem_Less(t *testing.T) {

	tests := []struct {
		name string
		ci   TCountItem
		item TCountItem
		want bool
	}{
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
/* */

/* EoF */
