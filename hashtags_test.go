/*
Copyright © 2019, 2025  M.Watermann, 10247 Berlin, Germany

	    All rights reserved
	EMail : <support@mwat.de>
*/
package hashtags

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

//lint:file-ignore ST1017 - I prefer Yoda conditions

func prepHT() *THashTags {
	fn := filepath.Join(os.TempDir(), ".testHtStore.db")
	os.Remove(fn) // remove remnants of previous runs

	ht, _ := New(fn)
	ht.safe = false // no locking wanted while testing

	return ht
} // prepHT()

func Test_New(t *testing.T) {
	testDir := t.TempDir()
	validFile := filepath.Join(testDir, "valid.db")

	tests := []struct {
		name     string
		filename string
		wantErr  bool
	}{
		{"valid new file", validFile, false},
		{"valid new file unsafe", validFile, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.filename)
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

func Test_THashTags_IDupdate(t *testing.T) {
	ht := prepHT()
	ht.safe = false
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

func Test_THashTags_parseID(t *testing.T) {
	hash1, hash2, hash3, hash4 := "#HÄSCH1", "#hash2", "#hash3", "#hash4"
	hyphTx1, hyphTx2, hyphTx3 := `#--------------`, `#---text ---`, `#-text-`

	id1, id2, id3, id4, id5, id6 := int64(987), int64(654), int64(321), int64(123), int64(456), int64(789)

	ht1 := prepHT()
	tx1 := []byte("1blabla " + hash1 + " blabla " + hash3 + ". Blabla")

	ht2 := prepHT()
	tx2 := []byte(`2blabla "` + hash2 + `". Blabla ` + hash3 + ` blabla`)

	ht3 := prepHT()
	tx3 := []byte("3\n> #KurzErklärt #Zensurheberrecht verhindern -\n> [Glyphosat-Gutachten selbst anfragen!](https://fragdenstaat.de/aktionen/zensurheberrecht-2019/)\n")

	ht4 := prepHT()
	tx4 := []byte("4blabla **" + hash1 + "** blabla\n\n_" + hash3 + "_")

	ht5 := prepHT()
	tx5 := []byte(`5blabla&#39; **` + hash2 + `** blabla\n<a href="page#fragment">txt</a> ` + hash4)

	ht6 := prepHT()
	tx6 := []byte(hash3 + ` blabla\n<a href="https://www.tagesspiegel.de/politik/martin-sonneborn-wirbt-fuer-moralische-integritaet-warum-ich-die-eu-kommission-ablehnen-werde/25263366.html#25263366">txt</a> ` + hash4)

	ht7 := prepHT()
	tx7 := []byte(`7 (https://www.faz.net/aktuell/politik/inland/jutta-ditfurth-zu-extinction-rebellion-irrationalismus-einer-endzeit-sekte-16422668.html?printPagedArticle=true#ageIndex_2)`)

	ht8 := prepHT()
	tx8 := []byte(`8
	> [Here's Everything You Need To Know](https://thehackernews.com/2018/12/australia-anti-encryption-bill.html#content) by <writer@example.com>
	`)

	ht9 := prepHT()
	tx9 := []byte(`9
	Bla *@Antoni_Comín* bla bla _#§219a_
	`)

	tmp := string(tx6) + "\n" + hyphTx1 + ` and ` + hyphTx2 + "\n" + hyphTx3
	tx10 := []byte(tmp)

	tx11 := []byte{}
	tx12 := []byte(" ")

	type tArgs struct {
		aID   int64
		aText []byte
	}
	tests := []struct {
		name string
		ht   *THashTags
		args tArgs
		want bool
	}{
		{"1", ht1, tArgs{id1, tx1}, true},
		{"2", ht2, tArgs{id2, tx2}, true},
		{"3", ht3, tArgs{id3, tx3}, true},
		{"4", ht4, tArgs{id4, tx4}, true},
		{"5", ht5, tArgs{id5, tx5}, true},
		{"6", ht6, tArgs{id6, tx6}, true},
		{"7", ht7, tArgs{7, tx7}, false},
		{"8", ht8, tArgs{8, tx8}, false},
		{"9", ht9, tArgs{9, tx9}, true},
		{"10", ht6, tArgs{id6, tx10}, false},
		{"11", ht1, tArgs{id1, tx11}, false},
		{"12", ht1, tArgs{id1, tx12}, false},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ht.parseID(tt.args.aID, tt.args.aText); got != tt.want {
				t.Errorf("%q: THashTags.parseID() = \n%v\n>>>> want >>>>\n%v\n{%s}",
					tt.name, got, tt.want, tt.ht)
			}
		})
	}
} // Test_THashList_parseID()

func Test_THashTags_removeHM(t *testing.T) {
	ht := prepHT()

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
	ht := prepHT()
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

func funcHashMentionRE(aText string) int {
	matches := htHashMentionRE.FindAllStringSubmatch(aText, -1)

	println(fmt.Sprintf("%s: %v", aText, matches))

	return len(matches)
} // funcHashMentionRE()

func Test_htHashMentionRE(t *testing.T) {
	t1 := `1blabla #HÄSCH1 blabla #hash2. Blabla`
	t2 := `2blabla #hash2. Blabla "#hash3" blabla`
	t3 := `\n>#KurzErklärt #Zensurheberrecht verhindern\n`
	t4 := `4blabla **#HÄSCH1** blabla\n\n_#hash3_`
	t5 := `5blabla&#39; **#hash2** blabla\n<a href="page#hash3">txt</a> #hash4`
	t6 := `#hash3 blabla\n<a href="https://www.tagesspiegel.de/politik/martin-sonneborn-wirbt-fuer-moralische-integritaet-warum-ich-die-eu-kommission-ablehnen-werde/25263366.html#25263366">txt</a> #hash4`
	t7 := `2blabla #hash2. @Dale_O'Leary "#hash3" @Dale_O’Leary blabla @Henry's`

	tests := []struct {
		name string
		text string
		want int
	}{
		{" 1", t1, 2},
		{" 2", t2, 2},
		{" 3", t3, 2},
		{" 4", t4, 2},
		{" 5", t5, 4},
		{" 6", t6, 3},
		{" 7", t7, 5},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := funcHashMentionRE(tt.text); got != tt.want {
				t.Errorf("%q: funcHashMentionRE() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // Test_htHashMentionRE()

/* EoF */
