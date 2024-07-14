/*
Copyright © 2019, 2024  M.Watermann, 10247 Berlin, Germany

		All rights reserved
	EMail : <support@mwat.de>
*/
package hashtags

import (
	"os"
	"testing"
)

//lint:file-ignore ST1017 - I prefer Yoda conditions

const (
	testHtStore = "testHtStore.db"
)

func TestTHashTags_equals(t *testing.T) {
	defer func() {
		os.Remove(testHtStore)
	}()

	ht1, _ := New("", false)
	wt1, _ := New("", false)
	wt2, _ := New("", false)
	wt2.HashAdd("hash1", 0)

	tests := []struct {
		name  string
		list  *THashTags
		other *THashTags
		want  bool
	}{
		{"1", ht1, wt1, true},
		{"s", ht1, wt2, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ht := tt.list
			if got := ht.equals(tt.other); got != tt.want {
				t.Errorf("%q: tHashTags.equals() =\n%v\n>>>> want >>>>\n%v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTHashTags_equals()

/* EoF */
