/*
Copyright © 2019, 2024  M.Watermann, 10247 Berlin, Germany

	    All rights reserved
	EMail : <support@mwat.de>
*/
package hashtags

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
)

//lint:file-ignore ST1017 - I prefer Yoda conditions

const (
	testHlStore = "testHlStore.db"
)

func htFilename() string {
	os.Remove(testHlStore)

	return testHlStore
} // htFilename()

func Test_newHashList(t *testing.T) {
	fn1 := htFilename() + "1"
	fn2 := htFilename() + "2"
	defer func() {
		os.Remove(fn1)
		os.Remove(fn2)
	}()

	// hash1, hash2 := "#hash1", "#hash2"
	// id1, id2 := uint64(654), uint64(321)

	hl1, _ := newHashList(fn1)
	// hl1.hashAdd(hash1, id1).
	// 	hashAdd(hash2, id2).
	// 	hashAdd(hash2, id1).
	// 	hashAdd(hash1, id2)
	// _, _ = hl1.Store()
	hl2, _ := newHashList(fn2)

	tests := []struct {
		tName   string
		fName   string
		want    *tHashList
		wantErr bool
	}{
		{" 1", fn1, hl1, false},
		{" 2", fn2, hl2, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.tName, func(t *testing.T) {
			got, err := newHashList(tt.fName)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%q: New() =\n{%v}\n>>>> want: >>>>\n{%v}",
					tt.tName, got, tt.want)
			}
		})
	}
} // Test_newHashList()

func TestTHashList_add(t *testing.T) {
	hl0 := &tHashList{
		hm: make(tHashMap, 64),
	}

	wl0 := &tHashList{
		hm: tHashMap{},
	}

	type tArgs struct {
		aDelim byte
		aName  string
		aID    uint64
	}
	tests := []struct {
		name   string
		fields *tHashList
		args   tArgs
		want   *tHashList
	}{
		{"1", hl0, tArgs{}, wl0},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hl := &tHashList{
				hm: tt.fields.hm,
			}
			got := hl.add(tt.args.aDelim, tt.args.aName, tt.args.aID)
			if !got.compareTo(tt.want) {
				t.Errorf("%q: THashList.add() =\n{%v}\n>>>> want: >>>>\n{%v}",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTHashList_add()

func TestTHashList_checksum(t *testing.T) {
	fn := htFilename()
	hash1, hash2 := "#hash1", "#hash2"
	id1, id2, id3 := uint64(987), uint64(654), uint64(321)

	hl1 := &tHashList{
		hm: tHashMap{
			hash1: &tSourceList{id2},
			hash2: &tSourceList{id3, id1},
		},
		//
	}
	h1a, _ := newHashList(fn)
	h1a.add(MarkHash, hash1, id2).
		add(MarkHash, hash2, id1).
		add(MarkHash, hash2, id3)
	w1 := h1a.hm.checksum()

	hl2 := &tHashList{
		hm: tHashMap{
			hash1: &tSourceList{id1, id2},
			hash2: &tSourceList{id2, id3},
		},
		//
	}
	h2a, _ := newHashList(fn)
	h2a.add(MarkHash, hash1, id1).
		add(MarkHash, hash1, id2).
		add(MarkHash, hash2, id2).
		add(MarkHash, hash2, id3).
		add(MarkHash, hash1, id2).
		add(MarkHash, hash2, id3)
	w2 := h2a.hm.checksum()

	tests := []struct {
		name string
		hl   *tHashList
		want uint32
	}{
		{"1", hl1, w1},
		{"2", hl2, w2},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// atomic.StoreUint32(&tt.hl.µChange, 0)
			if got := tt.hl.checksum(); got != tt.want {
				t.Errorf("%q: THashList.checksum() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTHashList_checksum()

func TestTHashList_clear(t *testing.T) {
	fn := htFilename()
	hash1, hash2 := "#hash1", "#hash2"
	id1, id2 := uint64(654), uint64(321)
	hl1, _ := newHashList(fn)
	hl1.add(MarkHash, hash1, id1).
		add(MarkHash, hash2, id2).
		add(MarkHash, hash2, id1).
		add(MarkHash, hash1, id2)
	tests := []struct {
		name string
		hl   *tHashList
		want int
	}{
		{" 1", hl1, 0},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hl.clear().Len(); got != tt.want {
				t.Errorf("%q: THashList.Clear() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTHashList_clear()

func TestTHashList_compare2(t *testing.T) {
	hl1 := &tHashList{
		hm: make(tHashMap, 64),
	}
	wl1 := &tHashList{
		hm: make(tHashMap, 64),
	}

	wl2 := &tHashList{
		hm: make(tHashMap, 64),
	}
	wl2.add(MarkHash, "hash2", 2222).add(MarkMention, "Name", 2222)

	tests := []struct {
		name string
		hl   *tHashList
		list *tHashList
		want bool
	}{
		{"1", hl1, wl1, true},
		{"2", hl1, wl2, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hl.compareTo(tt.list); got != tt.want {
				t.Errorf("THashList.compare2() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestTHashList_compare2()

func TestTHashList_Len(t *testing.T) {
	fn := htFilename()
	hl1, _ := newHashList(fn)

	hl2, _ := newHashList(fn)
	hl2.add(MarkHash, "#hash", 0)

	hl3, _ := newHashList(fn)
	hl3.add(MarkHash, "#hash2", 1).
		add(MarkHash, "#hash3", 2)

	hl4, _ := newHashList(fn)
	hl4.add(MarkHash, "#hash2", 1).
		add(MarkHash, "#hash3", 2).
		add(MarkHash, "#hash4", 3)
	tests := []struct {
		name string
		hl   *tHashList
		want int
	}{
		{" 1", hl1, 0},
		{" 2", hl2, 1},
		{" 3", hl3, 2},
		{" 4", hl4, 3},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hl.Len(); got != tt.want {
				t.Errorf("%q: THashList.Len() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTHashList_Len()

func TestTHashList_LenTotal(t *testing.T) {
	fn := htFilename()
	hash1, hash2, hash3 := "#hash1", "#hash2", "#hash3"
	id1, id2, id3 := uint64(987), uint64(654), uint64(321)

	hl1, _ := newHashList(fn)
	hl1.add(MarkHash, hash1, id1).
		add(MarkHash, hash2, id2).
		add(MarkHash, hash2, id1).
		add(MarkHash, hash1, id2)

	hl2, _ := newHashList(fn)
	hl2.add(MarkHash, hash1, id1).
		add(MarkHash, hash2, id2).
		add(MarkHash, hash2, id1).
		add(MarkHash, hash1, id2).
		add(MarkHash, hash3, id1).
		add(MarkHash, hash3, id2).
		add(MarkHash, hash3, id3)

	tests := []struct {
		name string
		hl   *tHashList
		want int
	}{
		{" 1", hl1, 6},
		{" 2", hl2, 10},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hl.LenTotal(); got != tt.want {
				t.Errorf("%q: THashList.LenTotal() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTHashList_LenTotal()

func TestTHashList_List(t *testing.T) {
	hash1, hash2, hash3 := "#hash1", "@mention1", "#another3"
	id1, id2, id3 := uint64(987), uint64(654), uint64(321)
	hl1 := &tHashList{
		hm: tHashMap{
			hash1: &tSourceList{id2},
			hash2: &tSourceList{id1, id3},
		},
	}
	wl1 := TCountList{
		{1, hash1},
		{2, hash2},
	}
	hl2 := &tHashList{
		hm: tHashMap{
			hash1: &tSourceList{id2},
			hash2: &tSourceList{id1, id3},
			hash3: &tSourceList{id1, id2, id3},
		},
	}
	wl2 := TCountList{
		{3, hash3},
		{1, hash1},
		{2, hash2},
	}

	tests := []struct {
		name string
		hl   *tHashList
		want TCountList
	}{
		{" 1", hl1, wl1},
		{" 2", hl2, wl2},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hl.List(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%q: THashList.list() = \n%v\n>>>> want: >>>>\n%v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTHashList_List()

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

func TestTHashList_parseID(t *testing.T) {
	hash1, hash2, hash3, hash4 := "#HÄSCH1", "#hash2", "#hash3", "#hash4"
	hyphTx1, hyphTx2, hyphTx3 := `#--------------`, `#---text ---`, `#-text-`
	lh1 := strings.ToLower(hash1)
	id1, id2, id3, id4, id5, id6 := uint64(987), uint64(654), uint64(321), uint64(123), uint64(456), uint64(789)
	// --------------
	hl1, _ := newHashList("")
	tx1 := []byte("1blabla " + hash1 + " blabla " + hash3 + ". Blabla")
	wl1 := &tHashList{
		hm: tHashMap{
			lh1:   &tSourceList{id1},
			hash3: &tSourceList{id1},
		},
	}
	// --------------
	hl2 := &tHashList{
		hm: tHashMap{
			hash1: &tSourceList{id3},
			hash2: &tSourceList{id3},
			hash3: &tSourceList{id3},
		},
	}
	tx2 := []byte(`2blabla "` + hash2 + `". Blabla ` + hash3 + ` blabla`)
	wl2 := &tHashList{
		hm: tHashMap{
			hash1: &tSourceList{id3},
			hash2: &tSourceList{id3, id2},
			hash3: &tSourceList{id3, id2},
		},
	}
	// --------------
	hl3, _ := newHashList("")
	tx3 := []byte("3\n> #KurzErklärt #Zensurheberrecht verhindern – \n> [Glyphosat-Gutachten selbst anfragen!](https://fragdenstaat.de/aktionen/zensurheberrecht-2019/)\n")
	wl3 := &tHashList{
		hm: tHashMap{
			"#kurzerklärt":      &tSourceList{id3},
			"#zensurheberrecht": &tSourceList{id3},
		},
	}
	// --------------
	hl4, _ := newHashList("")
	tx4 := []byte("4blabla **" + hash1 + "** blabla\n\n_" + hash3 + "_")
	wl4 := &tHashList{
		hm: tHashMap{
			lh1:   &tSourceList{id4},
			hash3: &tSourceList{id4},
		},
	}
	// --------------
	hl5, _ := newHashList("")
	tx5 := []byte(`5blabla&#39; **` + hash2 + `** blabla\n<a href="page#fragment">txt</a> ` + hash4)
	wl5 := &tHashList{
		hm: tHashMap{
			hash2: &tSourceList{id5},
			hash4: &tSourceList{id5},
		},
	}
	// --------------
	hl6, _ := newHashList("")
	tx6 := []byte(hash3 + ` blabla\n<a href="https://www.tagesspiegel.de/politik/martin-sonneborn-wirbt-fuer-moralische-integritaet-warum-ich-die-eu-kommission-ablehnen-werde/25263366.html#25263366">txt</a> ` + hash4)
	wl6 := &tHashList{
		hm: tHashMap{
			hash3: &tSourceList{id6},
			hash4: &tSourceList{id6},
		},
	}
	// --------------
	hl7, _ := newHashList("")
	tx7 := []byte(`7 (https://www.faz.net/aktuell/politik/inland/jutta-ditfurth-zu-extinction-rebellion-irrationalismus-einer-endzeit-sekte-16422668.html?printPagedArticle=true#ageIndex_2)`)
	wl7 := &tHashList{
		hm: tHashMap{},
	}
	// --------------
	hl8, _ := newHashList("")
	tx8 := []byte(`8
	> [Here's Everything You Need To Know](https://thehackernews.com/2018/12/australia-anti-encryption-bill.html#content)
	`)
	wl8 := &tHashList{
		hm: tHashMap{},
	}
	// --------------
	hl9, _ := newHashList("")
	tx9 := []byte(`9
	Bla *@Antoni_Comín* bla bla _#§219a_
	`)
	wl9 := &tHashList{
		hm: tHashMap{
			`@antoni_comín`: &tSourceList{9},
			`#§219a`:        &tSourceList{9},
		},
	}

	tmp := string(tx6) + "\n" + hyphTx1 + ` and ` + hyphTx2 + "\n" + hyphTx3
	tx10 := []byte(tmp)

	type tArgs struct {
		aID   uint64
		aText []byte
	}
	tests := []struct {
		name string
		hl   *tHashList
		args tArgs
		want *tHashList
	}{
		{"0", hl6, tArgs{id6, tx10}, wl6},
		{"1", hl1, tArgs{id1, tx1}, wl1},
		{"2", hl2, tArgs{id2, tx2}, wl2},
		{"3", hl3, tArgs{id3, tx3}, wl3},
		{"4", hl4, tArgs{id4, tx4}, wl4},
		{"5", hl5, tArgs{id5, tx5}, wl5},
		{"6", hl6, tArgs{id6, tx6}, wl6},
		{"7", hl7, tArgs{7, tx7}, wl7},
		{"8", hl8, tArgs{8, tx8}, wl8},
		{"9", hl9, tArgs{9, tx9}, wl9},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hl := &tHashList{
				hm: tt.hl.hm,
			}
			if got := hl.parseID(tt.args.aID, tt.args.aText); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%q: THashList.parseID() = \n%v\n>>>> want >>>>\n%v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTHashList_parseID()

func TestTHashList_remove(t *testing.T) {
	hash1, hash2, hash3 := "#hash1", "#hash2", "#hash3"
	id1, id2, id3 := uint64(987), uint64(654), uint64(321)

	hl1 := &tHashList{
		hm: tHashMap{
			hash1: &tSourceList{id1, id3},
			hash2: &tSourceList{id2, id3},
			hash3: &tSourceList{id1, id3},
		},
	}
	hl1.hm.sort()
	wl1 := &tHashList{
		hm: tHashMap{
			hash1: &tSourceList{id3},
			hash2: &tSourceList{id2, id3},
			hash3: &tSourceList{id1, id3},
		},
	}
	wl1.hm.sort()
	wl2 := &tHashList{
		hm: tHashMap{
			hash1: &tSourceList{id3},
			hash2: &tSourceList{id3},
			hash3: &tSourceList{id1, id3},
		},
	}
	wl2.hm.sort()
	wl3 := &tHashList{
		hm: tHashMap{
			hash1: &tSourceList{id3},
			hash2: &tSourceList{id3},
			hash3: &tSourceList{id1},
		},
	}
	wl3.hm.sort()
	wl4 := &tHashList{
		hm: tHashMap{
			hash1: &tSourceList{id3},
			hash2: &tSourceList{id3},
		},
	}
	wl5 := &tHashList{
		hm: tHashMap{
			hash2: &tSourceList{id3},
		},
	}
	wl6 := &tHashList{
		hm: tHashMap{
			hash2: &tSourceList{id3},
		},
	}
	wl7, _ := newHashList("")

	type tArgs struct {
		aDelim byte
		aName  string
		aID    uint64
	}
	tests := []struct {
		name string
		hl   *tHashList
		args tArgs
		want *tHashList
	}{
		{" 1", hl1, tArgs{MarkHash, hash1, id1}, wl1},
		{" 2", wl1, tArgs{MarkHash, hash2, id2}, wl2},
		{" 3", wl2, tArgs{MarkHash, hash3, id3}, wl3},
		{" 4", wl3, tArgs{MarkHash, hash3, id1}, wl4},
		{" 5", wl4, tArgs{MarkHash, hash1, id3}, wl5},
		{" 6", wl5, tArgs{MarkHash, hash1, id3}, wl6},
		{" 7", wl6, tArgs{MarkHash, hash2, id3}, wl7},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hl.removeHM(tt.args.aDelim, tt.args.aName, tt.args.aID); !got.compareTo(tt.want) {
				t.Errorf("%q: THashList.remove() =\n%v\n>>>> want: >>>>\n%v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTHashList_remove()

func TestTHashList_removeID(t *testing.T) {
	hash1, hash2, hash3 := "#hash1", "#hash2", "#hash3"
	id1, id2, id3 := uint64(987), uint64(654), uint64(321)
	hl1 := &tHashList{
		hm: tHashMap{
			hash1: &tSourceList{id1, id3},
			hash2: &tSourceList{id2, id3},
			hash3: &tSourceList{id1, id3},
		},
	}
	hl1.hm.sort()
	wl1 := &tHashList{
		hm: tHashMap{
			hash1: &tSourceList{id3},
			hash2: &tSourceList{id2, id3},
			hash3: &tSourceList{id3},
		},
	}
	wl1.hm.sort()
	wl2 := &tHashList{
		hm: tHashMap{
			hash1: &tSourceList{id3},
			hash2: &tSourceList{id3},
			hash3: &tSourceList{id3},
		},
	}
	wl2.hm.sort()
	wl3 := &tHashList{
		hm: tHashMap{},
	}

	tests := []struct {
		name string
		hl   *tHashList
		id   uint64
		want *tHashList
	}{
		{" 1", hl1, id1, wl1},
		{" 2", hl1, id2, wl2},
		{" 3", hl1, id3, wl3},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hl.IDremove(tt.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%q: THashList.removeID() = %v,\nwant %v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTHashList_removeID()

func TestTHashList_renameID(t *testing.T) {
	hash1, hash2, hash3 := "#hash1", "#hash2", "#hash3"
	id1, id2, id3, id4, id5, id6 := uint64(11), uint64(22), uint64(33), uint64(44), uint64(55), uint64(66)

	getHL := func() *tHashList {
		result := &tHashList{
			hm: tHashMap{
				hash1: &tSourceList{id1, id2},
				hash2: &tSourceList{id2, id3},
				hash3: &tSourceList{id1, id3},
			},
		}
		result.hm.sort()
		return result
	} // getHL()

	// hl1 := &THashList{
	// 	hm: tHashMap{
	// 		hash1: &tSourceList{id1, id2},
	// 		hash2: &tSourceList{id2, id3},
	// 		hash3: &tSourceList{id1, id3},
	// 	},
	// }
	// hl1.hm.sort()
	hl1 := getHL()
	wl1 := &tHashList{
		hm: tHashMap{
			hash1: &tSourceList{id2, id4},
			hash2: &tSourceList{id2, id3},
			hash3: &tSourceList{id3, id4},
		},
	}
	wl1.hm.sort()

	hl2 := getHL()
	wl2 := &tHashList{
		hm: tHashMap{
			hash1: &tSourceList{id1, id5},
			hash2: &tSourceList{id5, id3},
			hash3: &tSourceList{id1, id3},
		},
	}
	wl2.hm.sort()

	hl3 := getHL()
	wl3 := &tHashList{
		hm: tHashMap{
			hash1: &tSourceList{id1, id2},
			hash2: &tSourceList{id2, id6},
			hash3: &tSourceList{id1, id6},
		},
	}
	wl3.hm.sort()

	type tArgs struct {
		aOldID uint64
		aNewID uint64
	}
	tests := []struct {
		name string
		hl   *tHashList
		args tArgs
		want *tHashList
	}{
		{"1", hl1, tArgs{id1, id4}, wl1},
		{"2", hl2, tArgs{id2, id5}, wl2},
		{"3", hl3, tArgs{id3, id6}, wl3},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hl.IDrename(tt.args.aOldID, tt.args.aNewID); !got.compareTo(tt.want) {
				t.Errorf("%q: THashList.renameID() =\n%v\n>>>> want: >>>>\n%v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTHashList_renameID()

// func TestTHashList_SetFilename(t *testing.T) {
// 	hash1, hash2, hash3 := "#hash1", "#hash2", "#hash3"
// 	id1, id2, id3 := uint64(987), uint64(654), uint64(321)
// 	hl1 := &tHashList{
// 		hm: tHashMap{
// 			hash1: &tSourceList{id1, id3},
// 			hash2: &tSourceList{id2, id3},
// 			hash3: &tSourceList{id1, id3},
// 		},
// 	}

// 	tests := []struct {
// 		name  string
// 		hl    *tHashList
// 		fName string
// 		want  *tHashList
// 	}{
// 		{" 1", hl1, `fn1.db`, hl1},
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got := tt.hl.SetFilename(tt.fName)
// 			if (nil == got) || (got.fn != tt.fName) {
// 				t.Errorf("%q: THashList.SetFilename() = %v, want %v",
// 					tt.name, got, tt.want)
// 			}
// 		})
// 	}
// } // TestTHashList_SetFilename()

func TestTHashList_updateID(t *testing.T) {
	hash1, hash2, hash3 := "#hash1", "#hash2", "#hash3"
	id1, id2, id3 := uint64(987), uint64(654), uint64(321)

	hl1, _ := newHashList("")
	hl1.add(MarkHash, hash1, id1).
		add(MarkHash, hash1, id2).
		add(MarkHash, hash1, id3).
		add(MarkHash, hash2, id1).
		add(MarkHash, hash2, id2)
	tx1 := []byte("blabla " + hash1 + " blabla " + hash3 + " blabla")
	wl1 := &tHashList{
		hm: tHashMap{
			hash1: &tSourceList{id3, id2, id1},
			hash2: &tSourceList{id2},
			hash3: &tSourceList{id1},
		},
	}

	tx2 := []byte("blabla blabla blabla")
	wl2 := &tHashList{
		hm: tHashMap{
			hash1: &tSourceList{id3, id1},
			hash3: &tSourceList{id1},
		},
	}

	type tArgs struct {
		aID   uint64
		aText []byte
	}
	tests := []struct {
		name string
		hl   *tHashList
		args tArgs
		want *tHashList
	}{
		{" 1", hl1, tArgs{id1, tx1}, wl1},
		{" 2", hl1, tArgs{id2, tx2}, wl2},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.hl.IDupdate(tt.args.aID, tt.args.aText)
			if !got.compareTo(tt.want) {
				t.Errorf("%q: THashList.idUpdate() =\n%v\n>>>> want: >>>>\n%v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTHashList_updateID()

/* EoF */
