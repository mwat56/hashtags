package hashtags

import (
	"reflect"
	"testing"
)

func TestLoadList(t *testing.T) {
	fn := "hashlist.db"
	hash1, hash2 := "#hash1", "#hash2"
	id1, id2 := "id_c", "id_a"
	hl1 := NewList().
		HashAdd(hash1, id1).
		HashAdd(hash2, id2).
		HashAdd(hash2, id1).
		HashAdd(hash1, id2)
	hl1.Store(fn)
	hl2 := NewList()
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
		{" 2", args{"does.not.exist"}, hl2, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadList(tt.args.aFilename)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadList() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestLoadList()

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

func TestTHashList_HashAdd(t *testing.T) {
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
			if got := tt.hl.HashAdd(tt.args.aHash, tt.args.aSource).Len(); got != tt.want {
				t.Errorf("THashList.HashAdd() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestTHashList_HashAdd()

func TestTHashList_Clear(t *testing.T) {
	hash1, hash2 := "#hash1", "#hash2"
	id1, id2 := "id_c", "id_a"
	hl1 := NewList().
		HashAdd(hash1, id1).
		HashAdd(hash2, id2).
		HashAdd(hash2, id1).
		HashAdd(hash1, id2)
	tests := []struct {
		name string
		hl   *THashList
		want bool
	}{
		// TODO: Add test cases.
		{" 1", hl1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hl.Clear(); got != tt.want {
				t.Errorf("THashList.Clear() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestTHashList_Clear()

func TestTHashList_HashList(t *testing.T) {
	hash1, hash2 := "#hash1", "#hash2"
	id1, id2 := "id_c", "id_a"
	hl1 := NewList().
		HashAdd(hash1, id1).
		HashAdd(hash2, id2).
		HashAdd(hash2, id1).
		HashAdd(hash1, id2)
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

func TestTHashList_Len(t *testing.T) {
	hl1 := NewList()
	hl2 := NewList().HashAdd("#hash", "source")
	hl3 := NewList().HashAdd("#hash2", "source1")
	hl4 := NewList().HashAdd("#hash2", "source1").HashAdd("#hash3", "source2")
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
	fn := "hashlist.db"
	hash1, hash2 := "#hash1", "#hash2"
	id1, id2 := "id_c", "id_a"
	hl1 := NewList().
		HashAdd(hash1, id1).
		HashAdd(hash2, id2).
		HashAdd(hash2, id1).
		HashAdd(hash1, id2)
	hl1.Store(fn)
	hl1.Clear()
	hl2 := NewList()
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
		{" 1", hl1, args{fn}, 2, false},
		{" 2", hl2, args{".does.not.exist"}, 0, false},
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
		aDelim byte
		aID    string
		aText  []byte
	}
	tests := []struct {
		name string
		hl   *THashList
		args args
		want int
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

func TestTHashList_HashRemove(t *testing.T) {
	hash1, hash2 := "#hash1", "#hash2"
	id1, id2 := "id_c", "id_a"
	hl1 := NewList().
		HashAdd(hash1, id1).
		HashAdd(hash2, id2).
		HashAdd(hash2, id1).
		HashAdd(hash1, id2)
	type args struct {
		aHash string
		aID   string
	}
	tests := []struct {
		name string
		hl   *THashList
		args args
		want int //*THashList
	}{
		// TODO: Add test cases.
		{" 1", hl1, args{hash1, id1}, 1},
		{" 2", hl1, args{hash1, id2}, 0},
		{" 3", hl1, args{hash2, id1}, 1},
		{" 4", hl1, args{hash2, id2}, 0},
		{" 5", hl1, args{hash1, id1}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hl.HashRemove(tt.args.aHash, tt.args.aID).HashLen(tt.args.aHash); got != tt.want {
				t.Errorf("THashList.HashRemove() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestTHashList_HashRemove()

func TestTHashList_Store(t *testing.T) {
	fn := "hashlist.db"
	hash1, hash2 := "#hash1", "#hash2"
	id1, id2 := "id_c", "id_a"
	hl1 := NewList().
		HashAdd(hash1, id1).
		HashAdd(hash2, id2).
		HashAdd(hash2, id1).
		HashAdd(hash1, id2)
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
		{" 1", hl1, args{fn}, 38, false},
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
		HashAdd(hash1, id1).
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
				t.Errorf("THashList.String() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestTHashList_String()

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

func Test_tSourceList_remove(t *testing.T) {
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
		"five",
	}
	wl3 := &tSourceList{
		"two",
		"three",
	}
	wl4 := &tSourceList{
		"two",
	}
	wl5 := &tSourceList{}
	type args struct {
		aIdx int
	}
	tests := []struct {
		name string
		sl   *tSourceList
		args args
		want *tSourceList
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
} // Test_tSourceList_remove()

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
