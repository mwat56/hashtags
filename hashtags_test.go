/*
Copyright © 2019, 2025  M.Watermann, 10247 Berlin, Germany

	    All rights reserved
	EMail : <support@mwat.de>
*/
package hashtags

import (
	"os"
	"path/filepath"
	"testing"
)

//lint:file-ignore ST1017 - I prefer Yoda conditions

const (
	testHtStore = "testHtStore.db"
)

func TestNew(t *testing.T) {
	testDir := t.TempDir()
	validFile := filepath.Join(testDir, "valid.db")

	tests := []struct {
		name     string
		filename string
		safe     bool
		wantErr  bool
	}{
		{"valid new file", validFile, true, false},
		{"valid new file unsafe", validFile, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.filename, tt.safe)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got == nil {
					t.Error("New() returned nil, want non-nil THashTags")
					return
				}
				if got.Filename() != tt.filename {
					t.Errorf("New().Filename() = %v, want %v", got.Filename(), tt.filename)
				}
			}
		})
	}
} // TestNew()

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

func TestTHashTags_SetFilename(t *testing.T) {
	ht, _ := New("", false)
	tmpDir := t.TempDir() // Creates a temporary directory for testing

	tests := []struct {
		name     string
		filename string
		wantErr  bool
	}{
		{"empty filename", "", true},
		{"valid filename", filepath.Join(tmpDir, "test.db"), false},
		{"invalid directory", filepath.Join(tmpDir, "non-existent", "test.db"), true},
		{"whitespace filename", "   ", true},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ht.SetFilename(tt.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("THashTags.SetFilename() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && ht.fn != tt.filename {
				t.Errorf("THashTags.SetFilename() filename = %v, want %v", ht.fn, tt.filename)
			}
		})
	}
} // TestTHashTags_SetFilename()

/* EoF */
