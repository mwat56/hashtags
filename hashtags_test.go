package hashtags

import (
	"reflect"
	"testing"
)

func TestNewList(t *testing.T) {
	wl1 := make(THashList)
	tests := []struct {
		name string
		want *THashList
	}{
		// TODO: Add test cases.
		{" 1", &wl1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewList(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewList() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestNewList()

func TestTHashList_Add(t *testing.T) {
	hl1 := NewList()
	h1 := "#hash"
	s1 := "asource"
	s2 := "source2"
	h3 := "#another"
	type args struct {
		aHash   string
		aSource string
	}
	tests := []struct {
		name string
		hl   *THashList
		args args
		want int //*THashList
	}{
		// TODO: Add test cases.
		{" 0", hl1, args{h1, ""}, 0},
		{" 1", hl1, args{h1, s1}, 1},
		{" 2", hl1, args{h1, s2}, 1},
		{" 3", hl1, args{h3, s2}, 2},
		{" 4", hl1, args{"", s2}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hl.Add(tt.args.aHash, tt.args.aSource).Len(); got != tt.want {
				t.Errorf("THashList.Add() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestTHashList_Add()

func TestTHashList_Len(t *testing.T) {
	hl1 := NewList()
	hl2 := NewList().Add("#hash", "source")
	hl3 := NewList().Add("#hash2", "source1")
	hl4 := NewList().Add("#hash2", "source1").Add("#hash3", "source2")
	tests := []struct {
		name string
		hl   *THashList
		want int
	}{
		// TODO: Add test cases.
		{" 1", hl1, 0},
		{" 2", hl2, 1},
		{" 3", hl3, 1},
		{" 4", hl4, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hl.Len(); got != tt.want {
				t.Errorf("THashList.Len() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestTHashList_Len()

func TestTHashList_Load(t *testing.T) {
	hl1 := NewList()
	fn := "hashlist.db"
	type args struct {
		aFilename string
	}
	tests := []struct {
		name    string
		hl      *THashList
		args    args
		want    int //*THashList
		wantErr bool
	}{
		// TODO: Add test cases.
		{" 1", hl1, args{fn}, 2, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.hl.Load(tt.args.aFilename)
			if (err != nil) != tt.wantErr {
				t.Errorf("THashList.Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Len() != tt.want {
				t.Errorf("THashList.Load() = %v, want %v", got.Len(), tt.want)
			}
		})
	}
} // TestTHashList_Load()

func TestTHashList_parse(t *testing.T) {
	hl1 := NewList()
	s1 := []byte("What a #tag1,\n#tag2 and#tag3")
	s2 := []byte("Helle @mention!\nDid you see @other or@another?")
	type args struct {
		aDelim rune
		aID    string
		aText  []byte
	}
	tests := []struct {
		name string
		hl   *THashList
		args args
		want int // *THashList
	}{
		// TODO: Add test cases.
		{" 1", hl1, args{'#', "fn1", s1}, 2},
		{" 2", hl1, args{'@', "fn2", s2}, 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hl.parse(tt.args.aDelim, tt.args.aID, tt.args.aText).Len(); got != tt.want {
				t.Errorf("THashList.parse() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestTHashList_parse()

func TestTHashList_Store(t *testing.T) {
	hash1, hash2 := "#hash1", "#hash2"
	id1, id2 := "id_c", "id_a"
	hl1 := NewList().
		Add(hash1, id1).
		Add(hash2, id2).
		Add(hash2, id1).
		Add(hash1, id2)
	type args struct {
		aFilename string
	}
	tests := []struct {
		name    string
		hl      *THashList
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
		{" 1", hl1, args{"hashlist.db"}, 38, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.hl.Store(tt.args.aFilename)
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
	hash1, hash2 := "#hash1", "#hash2"
	id1, id2 := "id_c", "id_a"
	hl1 := NewList().
		Add(hash1, id1).
		Add(hash2, id2).
		Add(hash2, id1).
		Add(hash1, id2)
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
				t.Errorf("THashList.String() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestTHashList_String()

func TestTSourceList_remove(t *testing.T) {
	sl1 := &TSourceList{
		"one",
		"two",
		"three",
		"four",
		"five",
	}
	wl1 := &TSourceList{
		"two",
		"three",
		"four",
		"five",
	}
	wl2 := &TSourceList{
		"two",
		"three",
		"five",
	}
	wl3 := &TSourceList{
		"two",
		"three",
	}
	wl4 := &TSourceList{
		"two",
	}
	wl5 := &TSourceList{}
	type args struct {
		aIdx int
	}
	tests := []struct {
		name string
		sl   *TSourceList
		args args
		want *TSourceList
	}{
		// TODO: Add test cases.
		{" 1", sl1, args{0}, wl1},
		{" 2", sl1, args{2}, wl2},
		{" 3", sl1, args{2}, wl3},
		{" 4", sl1, args{1}, wl4},
		{" 5", sl1, args{0}, wl5},
		{" 6", sl1, args{0}, wl5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sl.remove(tt.args.aIdx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TSourceList.remove() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestTSourceList_remove()

func TestTSourceList_String(t *testing.T) {
	sl1 := &TSourceList{
		"one",
		"two",
		"three",
	}
	wl1 := "one\nthree\ntwo"
	sl2 := &TSourceList{}
	wl2 := ""
	tests := []struct {
		name string
		sl   *TSourceList
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
} // TestTSourceList_String()
