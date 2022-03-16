/*
   Copyright © 2019, 2022 M.Watermann, 10247 Berlin, Germany
                  All rights reserved
              EMail : <support@mwat.de>
*/
package hashtags

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
)

func delDB(aFilename string) string {
	os.Remove(aFilename)

	return aFilename
} // delDB()

func Test_tSourceList_indexOf(t *testing.T) {
	sl1 := &tSourceList{
		"one",
		"two",
		"three",
		"four",
		"five",
	}
	type args struct {
		aID string
	}
	tests := []struct {
		name string
		sl   *tSourceList
		args args
		want int
	}{
		// TODO: Add test cases.
		{" 1", sl1, args{"one"}, 0},
		{" 2", sl1, args{"two"}, 1},
		{" 3", sl1, args{"three"}, 2},
		{" 4", sl1, args{"four"}, 3},
		{" 5", sl1, args{"five"}, 4},
		{" 6", sl1, args{"six"}, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sl.indexOf(tt.args.aID); got != tt.want {
				t.Errorf("tSourceList.indexOf() = %v, want %v", got, tt.want)
			}
		})
	}
} // Test_tSourceList_indexOf()

func Test_tSourceList_removeID(t *testing.T) {
	sl1 := &tSourceList{
		"one",
		"two",
		"three",
		"four",
		"five",
	}
	wl1 := &tSourceList{
		"two",
		"three",
		"four",
		"five",
	}
	wl2 := &tSourceList{
		"two",
		"three",
		"four",
	}
	wl3 := &tSourceList{
		"two",
		"four",
	}
	type args struct {
		aID string
	}
	tests := []struct {
		name string
		sl   *tSourceList
		args args
		want *tSourceList
	}{
		// TODO: Add test cases.
		{" 1", sl1, args{"one"}, wl1},
		{" 2", sl1, args{"five"}, wl2},
		{" 3", sl1, args{"three"}, wl3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sl.removeID(tt.args.aID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tSourceList.IDremove() = %v, want %v", got, tt.want)
			}
		})
	}
} // Test_tSourceList_removeID()

func Test_tSourceList_renameID(t *testing.T) {
	sl1 := &tSourceList{
		"one",
		"two",
		"three",
	}
	wl1 := &tSourceList{
		"four",
		"one",
		"two",
	}
	wl2 := &tSourceList{
		"one",
		"six",
		"two",
	}
	type args struct {
		aOldID string
		aNewID string
	}
	tests := []struct {
		name string
		sl   *tSourceList
		args args
		want *tSourceList
	}{
		// TODO: Add test cases.
		{" 1", sl1, args{"three", "four"}, wl1},
		{" 2", sl1, args{"four", "six"}, wl2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sl.renameID(tt.args.aOldID, tt.args.aNewID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tSourceList.renameID() = %v, want %v", got, tt.want)
			}
		})
	}
} // Test_tSourceList_renameID()

func Test_tSourceList_String(t *testing.T) {
	sl1 := &tSourceList{
		"one",
		"two",
		"three",
	}
	wl1 := "one\nthree\ntwo"
	sl2 := &tSourceList{}
	wl2 := ""
	tests := []struct {
		name string
		sl   *tSourceList
		want string
	}{
		// TODO: Add test cases.
		{" 1", sl1, wl1},
		{" 2", sl2, wl2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sl.String(); got != tt.want {
				t.Errorf("TSourceList.String() = %v, want %v", got, tt.want)
			}
		})
	}
} // Test_tSourceList_String()

func TestNew(t *testing.T) {
	fn := delDB("hashlist.db")
	fn2 := delDB("does.not.exist")
	hash1, hash2 := "#hash1", "#hash2"
	id1, id2 := "id_c", "id_a"
	hl1, _ := New(fn)
	hl1.HashAdd(hash1, id1).
		HashAdd(hash2, id2).
		HashAdd(hash2, id1).
		HashAdd(hash1, id2)
	_, _ = hl1.Store()
	hl2, _ := New(fn2)
	type args struct {
		aFilename string
	}
	tests := []struct {
		name    string
		args    args
		want    *THashList
		wantErr bool
	}{
		// TODO: Add test cases.
		{" 1", args{fn}, hl1, false},
		{" 2", args{fn2}, hl2, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.aFilename)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestNew()

func TestTHashList_Checksum(t *testing.T) {
	fn := delDB("hashlist.db")
	hash1, hash2 := "#hash1", "#hash2"
	id1, id2, id3 := "id_c", "id_a", "id_b"
	hl1 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id2},
			hash2: &tSourceList{id3, id1},
		},
		mtx: new(sync.RWMutex),
	}
	h1a, _ := New(fn)
	h1a.HashAdd(hash1, id2).
		HashAdd(hash2, id1).
		HashAdd(hash2, id3)
	w1 := h1a.Checksum()
	hl2 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id1, id2},
			hash2: &tSourceList{id2, id3},
		},
		mtx: new(sync.RWMutex),
	}
	h2a, _ := New(fn)
	h2a.HashAdd(hash1, id1).
		HashAdd(hash1, id2).
		HashAdd(hash2, id2).
		HashAdd(hash2, id3).
		HashAdd(hash1, id2).
		HashAdd(hash2, id3)
	w2 := h2a.Checksum()
	tests := []struct {
		name     string
		hl       *THashList
		wantRSum uint32
	}{
		// TODO: Add test cases.
		{" 1", hl1, w1},
		{" 2", hl2, w2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			atomic.StoreUint32(&tt.hl.µChange, 0)
			if gotRSum := tt.hl.Checksum(); gotRSum != tt.wantRSum {
				t.Errorf("THashList.Checksum() = %v, want %v", gotRSum, tt.wantRSum)
			}
		})
	}
} // TestTHashList_Checksum()

func TestTHashList_Clear(t *testing.T) {
	fn := delDB("hashlist.db")
	hash1, hash2 := "#hash1", "#hash2"
	id1, id2 := "id_c", "id_a"
	hl1, _ := New(fn)
	hl1.HashAdd(hash1, id1).
		HashAdd(hash2, id2).
		HashAdd(hash2, id1).
		HashAdd(hash1, id2)
	tests := []struct {
		name string
		hl   *THashList
		want int
	}{
		// TODO: Add test cases.
		{" 1", hl1, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hl.Clear().Len(); got != tt.want {
				t.Errorf("THashList.Clear() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestTHashList_Clear()

func TestTHashList_count(t *testing.T) {
	hash1, hash2, hash3, hash4 := "#hash1", "@mention2", "#another3", "@mention4"
	id1, id2, id3 := "id_c", "id_a", "id_b"
	hl1 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id2},
			hash3: &tSourceList{id1, id2, id3},
		},
		mtx: new(sync.RWMutex),
	}
	hl2 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id2},
			hash2: &tSourceList{id1, id3},
			hash4: &tSourceList{id1, id3},
		},
		mtx: new(sync.RWMutex),
	}
	type args struct {
		aDelim byte
	}
	tests := []struct {
		name     string
		hl       *THashList
		args     args
		wantRLen int
	}{
		// TODO: Add test cases.
		{" 1", hl1, args{'#'}, 2},
		{" 2", hl1, args{'@'}, 0},
		{" 3", hl2, args{'#'}, 1},
		{" 4", hl2, args{'@'}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRLen := tt.hl.count(tt.args.aDelim); gotRLen != tt.wantRLen {
				t.Errorf("THashList.count() = %v, want %v", gotRLen, tt.wantRLen)
			}
		})
	}
} // TestTHashList_count()

func TestTHashList_CountedList(t *testing.T) {
	hash1, hash2, hash3 := "#hash1", "@mention1", "#another3"
	id1, id2, id3 := "id_c", "id_a", "id_b"
	hl1 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id2},
			hash2: &tSourceList{id1, id3},
		},
		mtx: new(sync.RWMutex),
	}
	wl1 := []TCountItem{
		{1, hash1},
		{2, hash2},
	}
	hl2 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id2},
			hash2: &tSourceList{id1, id3},
			hash3: &tSourceList{id1, id2, id3},
		},
		mtx: new(sync.RWMutex),
	}
	wl2 := []TCountItem{
		{3, hash3},
		{1, hash1},
		{2, hash2},
	}
	tests := []struct {
		name      string
		hl        *THashList
		wantRList []TCountItem
	}{
		// TODO: Add test cases.
		{" 1", hl1, wl1},
		{" 2", hl2, wl2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRList := tt.hl.CountedList(); !reflect.DeepEqual(gotRList, tt.wantRList) {
				t.Errorf("THashList.CountedList() = %v, want %v", gotRList, tt.wantRList)
			}
		})
	}
} // TestTHashList_CountedList()

func TestTHashList_HashAdd(t *testing.T) {
	hash1, hash2 := "#hash2", "#hash1"
	id1, id2, id3 := "id_c", "id_a", "id_b"
	hl1 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id2},
			hash2: &tSourceList{id1},
		},
		mtx: new(sync.RWMutex),
	}
	wl1 := &THashList{
		hl: tHashMap{
			hash2: &tSourceList{id1},
			hash1: &tSourceList{id2, id1},
		},
		mtx: new(sync.RWMutex),
	}
	wl2 := &THashList{
		hl: tHashMap{
			hash2: &tSourceList{id2, id1},
			hash1: &tSourceList{id2, id1},
		},
		mtx: new(sync.RWMutex),
	}
	wl3 := &THashList{
		hl: tHashMap{
			hash2: &tSourceList{id2, id1},
			hash1: &tSourceList{id2, id3, id1},
		},
		mtx: new(sync.RWMutex),
	}
	type args struct {
		aHash string
		aID   string
	}
	tests := []struct {
		name string
		hl   *THashList
		args args
		want *THashList
	}{
		// TODO: Add test cases.
		{" 1", hl1, args{hash1, id1}, wl1},
		{" 2", wl1, args{hash2, id2}, wl2},
		{" 3", wl2, args{hash1, id3}, wl3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hl.HashAdd(tt.args.aHash, tt.args.aID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("THashList.HashAdd() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestTHashList_HashAdd()

func TestTHashList_HashCount(t *testing.T) {
	hash1, hash2, hash3 := "#hash1", "@mention2", "#another3"
	id1, id2, id3 := "id_c", "id_a", "id_b"
	hl1 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id2},
			hash2: &tSourceList{id1, id2, id3},
		},
		mtx: new(sync.RWMutex),
	}
	hl2 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id2},
			hash2: &tSourceList{id1, id2, id3},
			hash3: &tSourceList{id1, id3},
		},
		mtx: new(sync.RWMutex),
	}
	tests := []struct {
		name string
		hl   *THashList
		want int
	}{
		// TODO: Add test cases.
		{" 1", hl1, 1},
		{" 2", hl2, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hl.HashCount(); got != tt.want {
				t.Errorf("THashList.HashCount() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestTHashList_HashCount()

func TestTHashList_HashLen(t *testing.T) {
	hash1, hash2 := "#hash1", "#hash2"
	id1, id2, id3 := "id_c", "id_a", "id_b"
	hl1 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id2},
			hash2: &tSourceList{id3, id1},
		},
		mtx: new(sync.RWMutex),
	}
	type args struct {
		aHash string
	}
	tests := []struct {
		name string
		hl   *THashList
		args args
		want int
	}{
		// TODO: Add test cases.
		{" 1", hl1, args{hash1}, 1},
		{" 2", hl1, args{hash2}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hl.HashLen(tt.args.aHash); got != tt.want {
				t.Errorf("THashList.HashLen() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestTHashList_HashLen()

func TestTHashList_HashList(t *testing.T) {
	hash1, hash2 := "#hash1", "#hash2"
	id1, id2 := "id_c", "id_a"
	hl1 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id2, id1},
			hash2: &tSourceList{id2, id1},
		},
		mtx: new(sync.RWMutex),
	}
	wl1 := []string{
		id2,
		id1,
	}
	type args struct {
		aHash string
	}
	tests := []struct {
		name      string
		hl        *THashList
		args      args
		wantRList []string
	}{
		// TODO: Add test cases.
		{" 1", hl1, args{hash1}, wl1},
		{" 2", hl1, args{hash2}, wl1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRList := tt.hl.HashList(tt.args.aHash); !reflect.DeepEqual(gotRList, tt.wantRList) {
				t.Errorf("THashList.HashList() = %v, want %v", gotRList, tt.wantRList)
			}
		})
	}
} // TestTHashList_HashList()

func TestTHashList_HashRemove(t *testing.T) {
	fn := delDB("hashlist.db")
	hash1, hash2 := "#hash1", "#hash2"
	id1, id2 := "id_c", "id_a"
	hl1 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id1, id2},
			hash2: &tSourceList{id1, id2},
		},
		fn:  fn,
		mtx: new(sync.RWMutex),
	}
	wl1 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id2},
			hash2: &tSourceList{id1, id2},
		},
		fn:  fn,
		mtx: new(sync.RWMutex),
	}
	wl2 := &THashList{
		hl: tHashMap{
			hash2: &tSourceList{id1, id2},
		},
		fn:  fn,
		mtx: new(sync.RWMutex),
	}
	wl3 := &THashList{
		hl: tHashMap{
			hash2: &tSourceList{id2},
		},
		fn:  fn,
		mtx: new(sync.RWMutex),
	}
	wl4, _ := New(fn)
	type args struct {
		aHash string
		aID   string
	}
	tests := []struct {
		name string
		hl   *THashList
		args args
		want *THashList
	}{
		// TODO: Add test cases.
		{" 1", hl1, args{hash1, id1}, wl1},
		{" 2", hl1, args{hash1, id2}, wl2},
		{" 3", hl1, args{hash2, id1}, wl3},
		{" 4", hl1, args{hash2, id2}, wl4},
		{" 5", hl1, args{hash1, id1}, wl4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hl.HashRemove(tt.args.aHash, tt.args.aID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("THashList.HashRemove() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestTHashList_HashRemove()

func TestTHashList_IDlist(t *testing.T) {
	hash1, hash2, hash3 := "#hash1", "#hash2", "#hash3"
	id1, id2, id3 := "id_c", "id_a", "id_b"
	hl1 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id1, id2},
			hash2: &tSourceList{id2, id3},
			hash3: &tSourceList{id1, id3},
		},
		mtx: new(sync.RWMutex),
	}
	var wl0 []string
	wl1 := []string{hash1, hash3}
	wl2 := []string{hash1, hash2}
	wl3 := []string{hash2, hash3}
	type args struct {
		aID string
	}
	tests := []struct {
		name      string
		hl        *THashList
		args      args
		wantRList []string
	}{
		// TODO: Add test cases.
		{" 0", hl1, args{"@does.not.exist"}, wl0},
		{" 1", hl1, args{id1}, wl1},
		{" 2", hl1, args{id2}, wl2},
		{" 3", hl1, args{id3}, wl3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRList := tt.hl.IDlist(tt.args.aID); !reflect.DeepEqual(gotRList, tt.wantRList) {
				t.Errorf("THashList.IDlist() = %v, want %v", gotRList, tt.wantRList)
			}
		})
	}
} // TestTHashList_IDlist()

func TestTHashList_IDremove(t *testing.T) {
	// fn := delDB("hashlist.db")
	hash1, hash2, hash3 := "#hash1", "#hash2", "#hash3"
	id1, id2, id3 := "id_c", "id_a", "id_b"
	hl1 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id1, id3},
			hash2: &tSourceList{id2, id3},
			hash3: &tSourceList{id1, id3},
		},
		mtx: new(sync.RWMutex),
	}
	wl1 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id3},
			hash2: &tSourceList{id2, id3},
			hash3: &tSourceList{id3},
		},
		mtx: new(sync.RWMutex),
	}
	wl2 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id3},
			hash2: &tSourceList{id3},
			hash3: &tSourceList{id3},
		},
		mtx: new(sync.RWMutex),
	}
	wl3, _ := New("")
	type args struct {
		aID string
	}
	tests := []struct {
		name string
		hl   *THashList
		args args
		want *THashList
	}{
		// TODO: Add test cases.
		{" 1", hl1, args{id1}, wl1},
		{" 2", wl1, args{id2}, wl2},
		{" 3", wl2, args{id3}, wl3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hl.IDremove(tt.args.aID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("THashList.IDremove() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestTHashList_IDremove()

func TestTHashList_IDrename(t *testing.T) {
	hash1, hash2, hash3 := "#hash1", "#hash2", "#hash3"
	id1, id2, id3, id4, id5, id6 := "id_e", "id_a", "id_c", "id_g", "id_j", "id_k"
	hl1 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id3, id1},
			hash2: &tSourceList{id2, id3},
			hash3: &tSourceList{id3, id1},
		},
		mtx: new(sync.RWMutex),
	}
	wl1 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id3, id4},
			hash2: &tSourceList{id2, id3},
			hash3: &tSourceList{id3, id4},
		},
		mtx: new(sync.RWMutex),
	}
	wl2 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id3, id4},
			hash2: &tSourceList{id3, id5},
			hash3: &tSourceList{id3, id4},
		},
		mtx: new(sync.RWMutex),
	}
	wl3 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id4, id6},
			hash2: &tSourceList{id5, id6},
			hash3: &tSourceList{id4, id6},
		},
		mtx: new(sync.RWMutex),
	}
	type args struct {
		aOldID string
		aNewID string
	}
	tests := []struct {
		name string
		hl   *THashList
		args args
		want *THashList
	}{
		// TODO: Add test cases.
		{" 1", hl1, args{id1, id4}, wl1},
		{" 2", wl1, args{id2, id5}, wl2},
		{" 3", wl2, args{id3, id6}, wl3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hl.IDrename(tt.args.aOldID, tt.args.aNewID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("THashList.IDrename() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestTHashList_IDrename()

func TestTHashList_IDupdate(t *testing.T) {
	hash1, hash2, hash3 := "#hash1", "#hash2", "#hash3"
	id1, id2, id3 := "id_c", "id_a", "id_b"
	hl1 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id1, id2, id3},
			hash2: &tSourceList{id1, id2},
		},
		mtx: new(sync.RWMutex),
	}
	tx1 := []byte("blabla " + hash1 + " blabla " + hash3 + " blabla")
	wl1 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id1, id2, id3},
			hash2: &tSourceList{id2},
			hash3: &tSourceList{id1},
		},
		mtx: new(sync.RWMutex),
	}
	tx2 := []byte("blabla blabla blabla")
	wl2 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id1, id3},
			hash3: &tSourceList{id1},
		},
		mtx: new(sync.RWMutex),
	}
	type args struct {
		aID   string
		aText []byte
	}
	tests := []struct {
		name string
		hl   *THashList
		args args
		want *THashList
	}{
		// TODO: Add test cases.
		{" 1", hl1, args{id1, tx1}, wl1},
		{" 2", hl1, args{id2, tx2}, wl2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hl.IDupdate(tt.args.aID, tt.args.aText); got.String() != tt.want.String() {
				t.Errorf("THashList.Update() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestTHashList_IDupdate()

func TestTHashList_Len(t *testing.T) {
	fn := delDB("hashlist.db")
	hl1, _ := New(fn)
	hl2, _ := New(fn)
	hl2.HashAdd("#hash", "source")
	hl3, _ := New(fn)
	hl3.HashAdd("#hash2", "source1").
		HashAdd("#hash3", "source2")
	hl4, _ := New(fn)
	hl4.HashAdd("#hash2", "source1").
		HashAdd("#hash3", "source2").
		HashAdd("#hash4", "source3")
	tests := []struct {
		name string
		hl   *THashList
		want int
	}{
		// TODO: Add test cases.
		{" 1", hl1, 0},
		{" 2", hl2, 1},
		{" 3", hl3, 2},
		{" 4", hl4, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hl.Len(); got != tt.want {
				t.Errorf("THashList.Len() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestTHashList_Len()

func TestTHashList_LenTotal(t *testing.T) {
	fn := delDB("hashlist.db")
	hash1, hash2, hash3 := "#hash1", "#hash2", "#hash3"
	id1, id2, id3 := "id_c", "id_a", "id_b"
	hl1, _ := New(fn)
	hl1.HashAdd(hash1, id1).
		HashAdd(hash2, id2).
		HashAdd(hash2, id1).
		HashAdd(hash1, id2)
	hl2, _ := New(fn)
	hl2.HashAdd(hash1, id1).
		HashAdd(hash2, id2).
		HashAdd(hash2, id1).
		HashAdd(hash1, id2).
		HashAdd(hash3, id1).
		HashAdd(hash3, id2).
		HashAdd(hash3, id3)
	tests := []struct {
		name       string
		hl         *THashList
		wantRCount int
	}{
		// TODO: Add test cases.
		{" 1", hl1, 6},
		{" 2", hl2, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRCount := tt.hl.LenTotal(); gotRCount != tt.wantRCount {
				t.Errorf("THashList.LenTotal() = %v, want %v", gotRCount, tt.wantRCount)
			}
		})
	}
} // TestTHashList_LenTotal()

func TestTHashList_MentionAdd(t *testing.T) {
	hash1, hash2 := "@mention2", "@mention1"
	id1, id2, id3 := "id_c", "id_a", "id_b"
	hl1 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id2},
			hash2: &tSourceList{id1},
		},
		mtx: new(sync.RWMutex),
	}
	wl1 := &THashList{
		hl: tHashMap{
			hash2: &tSourceList{id1},
			hash1: &tSourceList{id2, id1},
		},
		mtx: new(sync.RWMutex),
	}
	wl2 := &THashList{
		hl: tHashMap{
			hash2: &tSourceList{id2, id1},
			hash1: &tSourceList{id2, id1},
		},
		mtx: new(sync.RWMutex),
	}
	wl3 := &THashList{
		hl: tHashMap{
			hash2: &tSourceList{id2, id1},
			hash1: &tSourceList{id2, id3, id1},
		},
		mtx: new(sync.RWMutex),
	}
	type args struct {
		aHash string
		aID   string
	}
	tests := []struct {
		name string
		hl   *THashList
		args args
		want *THashList
	}{
		// TODO: Add test cases.
		{" 1", hl1, args{hash1, id1}, wl1},
		{" 2", wl1, args{hash2, id2}, wl2},
		{" 3", wl2, args{hash1, id3}, wl3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hl.MentionAdd(tt.args.aHash, tt.args.aID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("THashList.MentionAdd() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestTHashList_MentionAdd()

func TestTHashList_MentionCount(t *testing.T) {
	hash1, hash2, hash3 := "#hash1", "@mention2", "@another3"
	id1, id2, id3 := "id_c", "id_a", "id_b"
	hl1 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id2},
			hash2: &tSourceList{id1, id2, id3},
		},
		mtx: new(sync.RWMutex),
	}
	hl2 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id2},
			hash2: &tSourceList{id1, id2, id3},
			hash3: &tSourceList{id1, id3},
		},
		mtx: new(sync.RWMutex),
	}
	tests := []struct {
		name string
		hl   *THashList
		want int
	}{
		// TODO: Add test cases.
		{" 1", hl1, 1},
		{" 2", hl2, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hl.MentionCount(); got != tt.want {
				t.Errorf("THashList.MentionCount() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestTHashList_MentionCount()

func TestTHashList_MentionLen(t *testing.T) {
	hash1, hash2 := "@mention1", "@mention2"
	id1, id2, id3 := "id_c", "id_a", "id_b"
	hl1 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id2},
			hash2: &tSourceList{id3, id1},
		},
		mtx: new(sync.RWMutex),
	}
	type args struct {
		aHash string
	}
	tests := []struct {
		name string
		hl   *THashList
		args args
		want int
	}{
		// TODO: Add test cases.
		{" 1", hl1, args{hash1}, 1},
		{" 2", hl1, args{hash2}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hl.MentionLen(tt.args.aHash); got != tt.want {
				t.Errorf("THashList.MentionLen() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestTHashList_MentionLen()

func TestTHashList_MentionList(t *testing.T) {
	hash1, hash2 := "@mention1", "@mention2"
	id1, id2 := "id_c", "id_a"
	hl1 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id2, id1},
			hash2: &tSourceList{id2, id1},
		},
		mtx: new(sync.RWMutex),
	}
	wl1 := []string{
		id2,
		id1,
	}
	type args struct {
		aHash string
	}
	tests := []struct {
		name      string
		hl        *THashList
		args      args
		wantRList []string
	}{
		// TODO: Add test cases.
		{" 1", hl1, args{hash1}, wl1},
		{" 2", hl1, args{hash2}, wl1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRList := tt.hl.MentionList(tt.args.aHash); !reflect.DeepEqual(gotRList, tt.wantRList) {
				t.Errorf("THashList.MentionList() = %v, want %v", gotRList, tt.wantRList)
			}
		})
	}
} // TestTHashList_MentionList()

func TestTHashList_MentionRemove(t *testing.T) {
	fn := delDB("hashlist.db")
	hash1, hash2 := "@mention1", "@mention2"
	id1, id2 := "id_c", "id_a"
	hl1 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id1, id2},
			hash2: &tSourceList{id1, id2},
		},
		fn:  fn,
		mtx: new(sync.RWMutex),
	}
	wl1 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id2},
			hash2: &tSourceList{id1, id2},
		},
		fn:  fn,
		mtx: new(sync.RWMutex),
	}
	wl2 := &THashList{
		hl: tHashMap{
			hash2: &tSourceList{id1, id2},
		},
		fn:  fn,
		mtx: new(sync.RWMutex),
	}
	wl3 := &THashList{
		hl: tHashMap{
			hash2: &tSourceList{id2},
		},
		fn:  fn,
		mtx: new(sync.RWMutex),
	}
	wl4, _ := New(fn)
	type args struct {
		aHash string
		aID   string
	}
	tests := []struct {
		name string
		hl   *THashList
		args args
		want *THashList
	}{
		// TODO: Add test cases.
		{" 1", hl1, args{hash1, id1}, wl1},
		{" 2", hl1, args{hash1, id2}, wl2},
		{" 3", hl1, args{hash2, id1}, wl3},
		{" 4", hl1, args{hash2, id2}, wl4},
		{" 5", hl1, args{hash1, id1}, wl4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hl.MentionRemove(tt.args.aHash, tt.args.aID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("THashList.MentionRemove() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestTHashList_MentionRemove()

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
	type args struct {
		aText string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
		{" 1", args{t1}, 2},
		{" 2", args{t2}, 2},
		{" 3", args{t3}, 2},
		{" 4", args{t4}, 2},
		{" 5", args{t5}, 4},
		{" 6", args{t6}, 3},
		{" 7", args{t7}, 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := funcHashMentionRE(tt.args.aText); got != tt.want {
				t.Errorf("funcHashMentionRE() = %v, want %v", got, tt.want)
			}
		})
	}
} // Test_htHashMentionRE()

func TestTHashList_parseID(t *testing.T) {
	hash1, hash2, hash3, hash4 := "#HÄSCH1", "#hash2", "#hash3", "#hash4"
	hyphTx1, hyphTx2, hyphTx3 := `#--------------`, `#---text ---`, `#-text-`
	lh1 := strings.ToLower(hash1)
	id1, id2, id3, id4, id5, id6 := "id_c", "id_a", "id_b", "id_d", "id_e", "id_f"
	// --------------
	hl1, _ := New("")
	tx1 := []byte("1blabla " + hash1 + " blabla " + hash3 + ". Blabla")
	wl1 := &THashList{
		hl: tHashMap{
			lh1:   &tSourceList{id1},
			hash3: &tSourceList{id1},
		},
		mtx: new(sync.RWMutex),
	}
	// --------------
	hl2 := &THashList{
		hl: tHashMap{
			lh1:   &tSourceList{id3},
			hash2: &tSourceList{id3},
			hash3: &tSourceList{id3},
		},
		mtx: new(sync.RWMutex),
	}
	tx2 := []byte(`2blabla "` + hash2 + `". Blabla ` + hash3 + ` blabla`)
	wl2 := &THashList{
		hl: tHashMap{
			lh1:   &tSourceList{id3},
			hash2: &tSourceList{id2, id3},
			hash3: &tSourceList{id2, id3},
		},
		mtx: new(sync.RWMutex),
	}
	// --------------
	hl3, _ := New("")
	tx3 := []byte("3\n> #KurzErklärt #Zensurheberrecht verhindern – \n> [Glyphosat-Gutachten selbst anfragen!](https://fragdenstaat.de/aktionen/zensurheberrecht-2019/)\n")
	wl3 := &THashList{
		hl: tHashMap{
			"#kurzerklärt":      &tSourceList{id3},
			"#zensurheberrecht": &tSourceList{id3},
		},
		mtx: new(sync.RWMutex),
	}
	// --------------
	hl4, _ := New("")
	tx4 := []byte("4blabla **" + hash1 + "** blabla\n\n_" + hash3 + "_")
	wl4 := &THashList{
		hl: tHashMap{
			lh1:   &tSourceList{id4},
			hash3: &tSourceList{id4},
		},
		mtx: new(sync.RWMutex),
	}
	// --------------
	hl5, _ := New("")
	tx5 := []byte(`5blabla&#39; **` + hash2 + `** blabla\n<a href="page#fragment">txt</a> ` + hash4)
	wl5 := &THashList{
		hl: tHashMap{
			hash2: &tSourceList{id5},
			hash4: &tSourceList{id5},
		},
		mtx: new(sync.RWMutex),
	}
	// --------------
	hl6, _ := New("")
	tx6 := []byte(hash3 + ` blabla\n<a href="https://www.tagesspiegel.de/politik/martin-sonneborn-wirbt-fuer-moralische-integritaet-warum-ich-die-eu-kommission-ablehnen-werde/25263366.html#25263366">txt</a> ` + hash4)
	wl6 := &THashList{
		hl: tHashMap{
			hash3: &tSourceList{id6},
			hash4: &tSourceList{id6},
		},
		mtx: new(sync.RWMutex),
	}
	// --------------
	hl7, _ := New("")
	tx7 := []byte(`7 (https://www.faz.net/aktuell/politik/inland/jutta-ditfurth-zu-extinction-rebellion-irrationalismus-einer-endzeit-sekte-16422668.html?printPagedArticle=true#ageIndex_2)`)
	wl7 := &THashList{
		hl:  tHashMap{},
		mtx: new(sync.RWMutex),
	}
	// --------------
	hl8, _ := New("")
	tx8 := []byte(`8
> [Here's Everything You Need To Know](https://thehackernews.com/2018/12/australia-anti-encryption-bill.html#content)
`)
	wl8 := &THashList{
		hl:  tHashMap{},
		mtx: new(sync.RWMutex),
	}
	// --------------
	hl9, _ := New("")
	tx9 := []byte(`9
Bla *@Antoni_Comín* bla bla _#§219a_
`)
	wl9 := &THashList{
		hl: tHashMap{
			`@antoni_comín`: &tSourceList{`id9`},
			`#§219a`:        &tSourceList{`id9`},
		},
		mtx: new(sync.RWMutex),
	}

	tmp := string(tx6) + "\n" + hyphTx1 + ` and ` + hyphTx2 + "\n" + hyphTx3
	tx10 := []byte(tmp)

	type args struct {
		aID   string
		aText []byte
	}
	tests := []struct {
		name   string
		fields *THashList //fields
		args   args
		want   *THashList
	}{
		// TODO: Add test cases.
		{"10", hl6, args{id6, tx10}, wl6},
		{" 1", hl1, args{id1, tx1}, wl1},
		{" 2", hl2, args{id2, tx2}, wl2},
		{" 3", hl3, args{id3, tx3}, wl3},
		{" 4", hl4, args{id4, tx4}, wl4},
		{" 5", hl5, args{id5, tx5}, wl5},
		{" 6", hl6, args{id6, tx6}, wl6},
		{" 7", hl7, args{`id7`, tx7}, wl7},
		{" 8", hl8, args{`id8`, tx8}, wl8},
		{" 9", hl9, args{`id9`, tx9}, wl9},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hl := &THashList{
				fn:      tt.fields.fn,
				hl:      tt.fields.hl,
				mtx:     tt.fields.mtx,
				µChange: tt.fields.µChange,
				µCC:     tt.fields.µCC,
			}
			if got := hl.parseID(tt.args.aID, tt.args.aText); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("THashList.parseID() = '%v',\nwant '%v'", got, tt.want)
			}
		})
	}
} // TestTHashList_parseID()

func TestTHashList_remove(t *testing.T) {
	// fn := delDB("hashlist.db")
	hash1, hash2, hash3 := "#hash1", "#hash2", "#hash3"
	id1, id2, id3 := "id_3", "id_1", "id_2"
	hl1 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id1, id3},
			hash2: &tSourceList{id2, id3},
			hash3: &tSourceList{id1, id3},
		},
		mtx: new(sync.RWMutex),
	}
	wl1 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id3},
			hash2: &tSourceList{id2, id3},
			hash3: &tSourceList{id1, id3},
		},
		mtx: new(sync.RWMutex),
	}
	wl2 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id3},
			hash2: &tSourceList{id3},
			hash3: &tSourceList{id1, id3},
		},
		mtx: new(sync.RWMutex),
	}
	wl3 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id3},
			hash2: &tSourceList{id3},
			hash3: &tSourceList{id1},
		},
		mtx: new(sync.RWMutex),
	}
	wl4 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id3},
			hash2: &tSourceList{id3},
		},
		mtx: new(sync.RWMutex),
	}
	wl5 := &THashList{
		hl: tHashMap{
			hash2: &tSourceList{id3},
		},
		mtx: new(sync.RWMutex),
	}
	wl6 := &THashList{
		hl: tHashMap{
			hash2: &tSourceList{id3},
		},
		mtx: new(sync.RWMutex),
	}
	wl7, _ := New("")
	type args struct {
		aDelim  byte
		aMapIdx string
		aID     string
	}
	tests := []struct {
		name string
		hl   *THashList
		args args
		want *THashList
	}{
		// TODO: Add test cases.
		{" 1", hl1, args{'#', hash1, id1}, wl1},
		{" 2", wl1, args{'#', hash2, id2}, wl2},
		{" 3", wl2, args{'#', hash3, id3}, wl3},
		{" 4", wl3, args{'#', hash3, id1}, wl4},
		{" 5", wl4, args{'#', hash1, id3}, wl5},
		{" 6", wl5, args{'#', hash1, id3}, wl6},
		{" 7", wl6, args{'#', hash2, id3}, wl7},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hl.remove(tt.args.aDelim, tt.args.aMapIdx, tt.args.aID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("THashList.remove() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestTHashList_remove()

func TestTHashList_removeID(t *testing.T) {
	hash1, hash2, hash3 := "#hash1", "#hash2", "#hash3"
	id1, id2, id3 := "id_3", "id_1", "id_2"
	hl1 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id1, id3},
			hash2: &tSourceList{id2, id3},
			hash3: &tSourceList{id1, id3},
		},
		mtx: new(sync.RWMutex),
	}
	wl1 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id3},
			hash2: &tSourceList{id2, id3},
			hash3: &tSourceList{id3},
		},
		mtx: new(sync.RWMutex),
	}
	wl2 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id3},
			hash2: &tSourceList{id3},
			hash3: &tSourceList{id3},
		},
		mtx: new(sync.RWMutex),
	}
	wl3 := &THashList{
		hl:  tHashMap{},
		mtx: new(sync.RWMutex),
	}
	type args struct {
		aID string
	}
	tests := []struct {
		name string
		hl   *THashList
		args args
		want *THashList
	}{
		// TODO: Add test cases.
		{" 1", hl1, args{id1}, wl1},
		{" 2", hl1, args{id2}, wl2},
		{" 3", hl1, args{id3}, wl3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hl.removeID(tt.args.aID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("THashList.removeID() = %v,\nwant %v", got, tt.want)
			}
		})
	}
} // TestTHashList_removeID()

func TestTHashList_SetFilename(t *testing.T) {
	hash1, hash2, hash3 := "#hash1", "#hash2", "#hash3"
	id1, id2, id3 := "id_3", "id_1", "id_2"
	hl1 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id1, id3},
			hash2: &tSourceList{id2, id3},
			hash3: &tSourceList{id1, id3},
		},
		mtx: new(sync.RWMutex),
	}
	type args struct {
		aFilename string
	}
	tests := []struct {
		name string
		hl   *THashList
		args args
		want *THashList
	}{
		// TODO: Add test cases.
		{" 1", hl1, args{`fn1.db`}, hl1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.hl.SetFilename(tt.args.aFilename)
			if (nil == got) || (got.fn != tt.args.aFilename) {
				t.Errorf("THashList.SetFilename() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestTHashList_SetFilename()

func TestTHashList_store(t *testing.T) {
	fn := delDB("hashlist.db")
	hash1, hash2 := "#hash1", "#Zensurheberrecht"
	id1, id2 := "id_c", "id_a"
	hl1, _ := New(fn)
	hl1.HashAdd(hash1, id1).
		HashAdd(hash2, id2).
		HashAdd(hash2, id1).
		HashAdd(hash1, id2)
	hl2, _ := New("")
	hl2.HashAdd(hash1, id1).
		HashAdd(hash2, id2).
		HashAdd(hash2, id1).
		HashAdd(hash1, id2)
	tests := []struct {
		name    string
		hl      *THashList
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
		{" 1", hl1, 91, false},
		{" 2", hl2, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.hl.store()
			if (err != nil) != tt.wantErr {
				t.Errorf("THashList.Store() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("THashList.Store() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestTHashList_Store()

func TestTHashList_String(t *testing.T) {
	fn := delDB("hashlist.db")
	hash1, hash2 := "#hash1", "#hash2"
	id1, id2 := "id_c", "id_a"
	hl1, _ := New(fn)
	hl1.HashAdd(hash1, id1).
		HashAdd(hash2, id2).
		HashAdd(hash2, id1).
		HashAdd(hash1, id2)
	wl1 := "[" + hash1 + "]\n" + id2 + "\n" + id1 +
		"\n[" + hash2 + "]\n" + id2 + "\n" + id1 + "\n"
	tests := []struct {
		name string
		hl   *THashList
		want string
	}{
		// TODO: Add test cases.
		{" 1", hl1, wl1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hl.String(); got != tt.want {
				t.Errorf("THashList.String() = {%v},\nwant {%v}", got, tt.want)
			}
		})
	}
} // TestTHashList_String()

func Benchmark_LoadTxT(b *testing.B) {
	hl, _ := New("")
	hl.SetFilename("load.txt")
	UseBinaryStorage = false
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		if _, err := hl.Load(); nil != err {
			log.Printf("LoadTxt(): %v", err)
		}
	}
} // Benchmark_LoadTxt()

func Benchmark_LoadBin(b *testing.B) {
	hl, _ := New("")
	hl.SetFilename("load.db")
	UseBinaryStorage = true
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		if _, err := hl.Load(); nil != err {
			log.Printf("LoadBin(): %v", err)
		}
	}
} // Benchmark_LoadBin()

func Benchmark_StoreTxt(b *testing.B) {
	UseBinaryStorage = false
	hl, _ := New("load.txt")
	_, _ = hl.Load()

	for n := 0; n < b.N; n++ {
		if _, err := hl.Store(); nil != err {
			log.Printf("StoreTxt(): %v", err)
		}
	}
} // Benchmark_StoreTxt()

func Benchmark_StoreBin(b *testing.B) {
	UseBinaryStorage = true
	hl, _ := New("load.db")
	_, _ = hl.Load()

	for n := 0; n < b.N; n++ {
		if _, err := hl.Store(); nil != err {
			log.Printf("StoreBin(): %v", err)
		}
	}
} // Benchmark_StoreBin()
