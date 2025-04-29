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

func Test_New(t *testing.T) {
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
			if (nil != err) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if nil == got {
					t.Error("New() returned nil, want non-nil THashTags")
					return
				}
				if got.Filename() != tt.filename {
					t.Errorf("New().Filename() = %v, want %v", got.Filename(), tt.filename)
				}
			}
		})
	}
} // Test_New()

func Test_THashTags_equals(t *testing.T) {
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
				t.Errorf("%q: tHashTags.equals() = '%v', want '%v'",
					tt.name, got, tt.want)
			}
		})
	}
} // Test_THashTags_equals()

func Test_THashTags_IDupdate(t *testing.T) {
	ht, _ := New("", false)
	id := int64(1)

	tests := []struct {
		name string
		id   int64
		text []byte
		want bool
	}{
		{"empty text", id, []byte(""), false},
		{"with hashtag", id, []byte("This is a #test"), true},
		{"with mention", id, []byte("Hello @world"), true},
		{"with both", id, []byte("Hello @world and #test"), true},
		{"no tags", id, []byte("Plain text without tags"), true}, // remove tags from previous test
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ht.IDupdate(tt.id, tt.text); got != tt.want {
				t.Errorf("%q: THashTags.IDupdate() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // Test_THashTags_IDupdate()

func Test_THashTags_removeHM(t *testing.T) {
	ht, _ := New("", false)

	// Setup test data
	ht.HashAdd("#test1", 101)
	ht.HashAdd("#test2", 102)
	ht.MentionAdd("@user1", 201)
	ht.MentionAdd("@user2", 202)

	type tArgs struct {
		aDelim byte
		aName  string
		aID    int64
	}
	tests := []struct {
		name string
		args tArgs
		want bool
	}{
		{"existing hash", tArgs{MarkHash, "test1", 101}, true},
		{"non-existing hash", tArgs{MarkHash, "test3", 103}, false},
		{"existing mention", tArgs{MarkMention, "user1", 201}, true},
		{"non-existing mention", tArgs{MarkMention, "user3", 203}, false},
		{"with empty name", tArgs{MarkHash, "", 101}, false},
		{"with space name", tArgs{MarkHash, " ", 102}, false},
		{"with invalid delimiter", tArgs{'X', "test2", 103}, false},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ht.removeHM(tt.args.aDelim, tt.args.aName, tt.args.aID); got != tt.want {
				t.Errorf("THashTags.removeHM() = %v, want %v", got, tt.want)
			}
		})
	}
} // Test_THashTags_removeHM()

func Test_THashTags_SetFilename(t *testing.T) {
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
			if (nil != err) != tt.wantErr {
				t.Errorf("THashTags.SetFilename() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && ht.fn != tt.filename {
				t.Errorf("THashTags.SetFilename() filename = %v, want %v", ht.fn, tt.filename)
			}
		})
	}
} // Test_THashTags_SetFilename()

/* EoF */
