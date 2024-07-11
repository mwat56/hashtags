/*
Copyright © 2023, 2024  M.Watermann, 10247 Berlin, Germany

		All rights reserved
	EMail : <support@mwat.de>
*/
package hashtags

import (
	"log"
	"os"
	"reflect"
	"testing"
)

//lint:file-ignore ST1017 - I prefer Yoda conditions

const (
	testHmStore = "testHmStore.db"
)

func hmFilename() string {
	os.Remove(testHmStore)

	return testHmStore
} // hmFilename()

func prepHashMap() *tHashMap {
	hm := make(tHashMap, 8)
	// add already sorts keys
	hm.add("#hash1", 111).
		add("#hash2", 222).
		add("#hash3", 333).
		add("@mention1", 111).
		add("@mention2", 222)

	return &hm
} // prepHashMap()

func Test_tHashMap_add(t *testing.T) {
	hm1 := prepHashMap()
	wm1 := prepHashMap().add("@mention2", 222)

	hm2 := prepHashMap().add("#hash4", 444)
	wm2 := prepHashMap().add("#hash4", 444).add("#hash4", 222)

	type tArgs struct {
		aName string
		aID   uint64
	}
	tests := []struct {
		name string
		hm   *tHashMap
		args tArgs
		want tHashMap
	}{
		{"1", hm1, tArgs{"@mention2", 222}, *wm1},
		{"2", hm2, tArgs{"#hash4", 222}, *wm2},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hm.add(tt.args.aName, tt.args.aID); !got.compareTo(tt.want) {
				t.Errorf("%q: tHashMap.add() =\n%v\n>>>> want: >>>>\n%v",
					tt.name, got, tt.want)
			}
		})
	}
} // Test_tHashMap_add()

func Test_tHashMap_clear(t *testing.T) {
	hm1 := prepHashMap()
	wm1 := tHashMap{}

	tests := []struct {
		name string
		hm   *tHashMap
		want tHashMap
	}{
		{"1", hm1, wm1},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hm.clear(); !got.compareTo(tt.want) {
				t.Errorf("%q: tHashMap.clear() =\n%v\n>>>> want >>>>\n%v",
					tt.name, got, tt.want)
			}
		})
	}
} // Test_tHashMap_clear()

func Test_tHashMap_compareTo(t *testing.T) {
	hm1 := prepHashMap()
	om1 := prepHashMap()
	om2 := prepHashMap().add("#hash4", 222)

	tests := []struct {
		name string
		hm   *tHashMap
		oMap *tHashMap
		want bool
	}{
		{"1", hm1, om1, true},
		{"2", hm1, om2, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hm.compareTo(*tt.oMap); got != tt.want {
				t.Errorf("%q: tHashMap.compareTo() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // Test_tHashMap_compareTo()

func Test_tHashMap_count(t *testing.T) {
	hm1 := prepHashMap()
	wm1 := 3
	hm2 := prepHashMap().add("#hash4", 222).sort()
	wm2 := 4
	hm3 := prepHashMap().add("@mention3", 222).sort()
	wm3 := 3

	tests := []struct {
		name    string
		hm      tHashMap
		delim   byte
		wantInt int
	}{
		{"1", *hm1, MarkHash, wm1},
		{"2", *hm2, MarkHash, wm2},
		{"3", *hm3, MarkMention, wm3},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotInt := tt.hm.count(tt.delim); gotInt != tt.wantInt {
				t.Errorf("%q: tHashMap.count() = %v, want %v",
					tt.name, gotInt, tt.wantInt)
			}
		})
	}
} // Test_tHashMap_count()

func Test_tHashMap_countedList(t *testing.T) {
	hm1 := prepHashMap()

	wc1 := TCountList{
		{1, "#hash1"},
		{1, "#hash2"},
		{1, "#hash3"},
		{1, "@mention1"},
		{1, "@mention2"},
	}

	hm2 := prepHashMap().add("@Alphons", 222)
	wc2 := TCountList{
		{1, "#hash1"},
		{1, "#hash2"},
		{1, "#hash3"},
		{1, "@alphons"},
		{1, "@mention1"},
		{1, "@mention2"},
	}

	tests := []struct {
		name string
		hm   tHashMap
		want TCountList
	}{
		{"1", *hm1, wc1},
		{"2", *hm2, wc2},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hm.countedList(); !got.compareTo(tt.want) {
				t.Errorf("%q: tHashMap.countedList() = \n%v\n>>>> want: >>>>\n%v",
					tt.name, got, tt.want)
			}
		})
	}
} // Test_tHashMap_countedList()

func Test_tHashMap_idList(t *testing.T) {
	hm := prepHashMap()

	id1 := uint64(111)
	wl1 := []string{
		"#hash1",
		"@mention1",
	}

	id2 := uint64(222)
	wl2 := []string{
		"#hash2",
		"@mention2",
	}

	tests := []struct {
		name string
		id   uint64
		want []string
	}{
		{"1", id1, wl1},
		{"2", id2, wl2},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hm.idList(tt.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%q: tHashMap.idList() =\n%v\n>>>> want: >>>>\n%v",
					tt.name, got, tt.want)
			}
		})
	}
} // Test_tHashMap_idList()

func Test_tHashMap_idxLen(t *testing.T) {
	hm := prepHashMap().
		add("#hash2", 333)

	type tArgs struct {
		aDelim byte
		aName  string
	}
	tests := []struct {
		name string
		args tArgs
		want int
	}{
		{"0", tArgs{}, -1},
		{"1", tArgs{MarkHash, "#hash2"}, 2},
		{"2", tArgs{MarkHash, "hash2"}, 2},
		{"3", tArgs{MarkMention, "@hash2"}, -1},
		{"4", tArgs{MarkMention, "mention1"}, 1},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hm.idxLen(tt.args.aDelim, tt.args.aName); got != tt.want {
				t.Errorf("%q: tHashMap.idxLen() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // Test_tHashMap_idxLen()

func Test_tHashMap_keys(t *testing.T) {
	// hm.add("#hash1", 111).
	// 	add("#hash2", 222).
	// 	add("#hash3", 333).
	// 	add("@mention1", 111).
	// 	add("@mention2", 222)
	hm0 := &tHashMap{}
	var wl0 []string
	hm1 := prepHashMap().add("#hash4", 444)
	wm1 := []string{"#hash1", "#hash2", "#hash3", "#hash4", "@mention1", "@mention2"}

	tests := []struct {
		name string
		hm   *tHashMap
		want []string
	}{
		{"0", hm0, wl0},
		{"1", hm1, wm1},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hm.keys(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%q: tHashMap.keys() =\n%v\n>>>> want: >>>>\n%v",
					tt.name, got, tt.want)
			}
		})
	}
} // Test_tHashMap_keys()

func Test_tHashMap_list(t *testing.T) {
	hm := prepHashMap().
		add("#hash3", 33)
	wl0 := tSourceList{}
	wl1 := tSourceList{
		111,
	}
	wl2 := tSourceList{
		33,
		333,
	}

	type tArgs struct {
		aDelim byte
		aName  string
	}
	tests := []struct {
		name     string
		args     tArgs
		wantList tSourceList
	}{
		{"0", tArgs{}, wl0},
		{"1", tArgs{MarkHash, "hash1"}, wl1},
		{"2", tArgs{MarkHash, "hash3"}, wl2},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotList := hm.list(tt.args.aDelim, tt.args.aName)
			if !tt.wantList.compareTo(gotList) {
				t.Errorf("%q: tHashMap.list() =\n%v\n>>>> want: >>>>\n%v",
					tt.name, gotList, tt.wantList)
			}
		})
	}
} // Test_tHashMap_list()

func Test_tHashMap_Load(t *testing.T) {
	saveBinary := UseBinaryStorage
	defer func() {
		os.Remove(testHmStore)
		UseBinaryStorage = saveBinary
	}()
	hm1 := prepHashMap().
		add("@CrashTestDummy", 1)
	wm1 := prepHashMap().
		add("@CrashTestDummy", 1)

	tests := []struct {
		name    string
		hm      *tHashMap
		binary  bool
		want    tHashMap
		wantErr bool
	}{
		{"1", hm1, true, *wm1, false},
		{"2", hm1, false, *wm1, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			UseBinaryStorage = tt.binary
			tt.hm.store(testHmStore)
			got, err := tt.hm.Load(testHmStore)
			if (nil != err) != tt.wantErr {
				t.Errorf("%q: tHashMap.Load() =\n%v\n>>>> want >>>>\n%v",
					tt.name, err, tt.wantErr)
				return
			}
			if !tt.want.compareTo(*got) {
				t.Errorf("%q: tHashMap.Load() =\n%v\n>>>> want >>>>\n%v",
					tt.name, got, tt.want)
			}
		})
	}
} // Test_tHashMap_Load()

func Test_tHashMap_remove(t *testing.T) {
	hm := prepHashMap().
		add("#nameX", 999).
		add("#hash3", 333)
	wm0 := hm

	wm1 := prepHashMap().
		add("#nameX", 999).
		add("#hash3", 333)

	wm2 := prepHashMap().
		add("#hash3", 333)

	type tArgs struct {
		aDelim byte
		aName  string
		aID    uint64
	}
	tests := []struct {
		name string
		args tArgs
		want tHashMap
	}{
		{"0", tArgs{}, *wm0},
		{"1", tArgs{MarkHash, "#nameX", 111}, *wm1},
		{"2", tArgs{MarkHash, "#nameX", 999}, *wm2},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hm.remove(tt.args.aDelim, tt.args.aName, tt.args.aID)
			if !got.compareTo(tt.want) {
				t.Errorf("%q: tHashMap.remove() =\n%v\n>>>> want >>>>\n%v",
					tt.name, got, tt.want.String())
			}
		})
	}
} // Test_tHashMap_remove()

func Test_tHashMap_removeID(t *testing.T) {
	hm0 := newHashMap()
	wm0 := newHashMap()
	hm1 := prepHashMap().
		add("#hash2", 999).
		add("#hash3", 888).
		add("#CrashTestDummy", 777)
	wm1 := prepHashMap().
		add("#CrashTestDummy", 777).
		add("#hash3", 888)
	wm2 := prepHashMap().
		add("#CrashTestDummy", 777)
	wm3 := prepHashMap()

	tests := []struct {
		name string
		hm   *tHashMap
		id   uint64
		want *tHashMap
	}{
		{"0", &hm0, 999, &wm0},
		{"1", hm1, 999, wm1},
		{"2", hm1, 888, wm2},
		{"3", hm1, 777, wm3},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hm.removeID(tt.id); !tt.want.compareTo(*got) {
				t.Errorf("%q: tHashMap.removeID() =\n%v\n>>>> want >>>>\n%v",
					tt.name, got, tt.want.String())
			}
		})
	}
} // Test_tHashMap_removeID()

func Test_tHashMap_renameID(t *testing.T) {
	id1, id2, id3 := uint64(11), uint64(22), uint64(33)

	hm1 := prepHashMap().
		add("#hash1", id1).
		add("#hash1", id3).
		add("@mention1", id3).
		sort()
	wm1 := prepHashMap().
		add("#hash1", id2).
		add("#hash1", id3).
		add("@mention1", id3).
		sort()
	if 0 == id1 {
		return
	}

	hm2 := prepHashMap().
		add("#hash3", id2).
		add("#hash4", id2).
		sort()
	wm2 := prepHashMap().
		add("#hash3", id3).
		add("#hash4", id3).
		sort()

	type tArgs struct {
		aOldID, aNewID uint64
	}
	tests := []struct {
		name string
		hm   *tHashMap
		args tArgs
		want *tHashMap
	}{
		{"0", hm1, tArgs{}, hm1}, // no change
		{"1", hm1, tArgs{id1, id2}, wm1},
		{"2", hm2, tArgs{id2, id3}, wm2},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hm.renameID(tt.args.aOldID, tt.args.aNewID); !got.compareTo(*tt.want) {
				t.Errorf("%q: tHashMap.renameID() =\n%v\n>>>> want: >>>>\n%v",
					tt.name, got, tt.want)
			}
		})
	}
} // Test_tHashMap_renameID

func Test_tHashMap_sort(t *testing.T) {
	hm1 := &tHashMap{
		"#hash1": &tSourceList{
			uint64(111),
		},
		"@mention1": &tSourceList{
			uint64(111),
		},
		"#hash2": &tSourceList{
			uint64(222),
		},
		"@mention2": &tSourceList{
			uint64(333),
			uint64(222),
		},
		"#hash3": &tSourceList{
			uint64(333),
		},
	}
	wm1 := tHashMap{
		"#hash1": &tSourceList{
			uint64(111),
		},
		"#hash2": &tSourceList{
			uint64(222),
		},
		"#hash3": &tSourceList{
			uint64(333),
		},
		"@mention1": &tSourceList{
			uint64(111),
		},
		"@mention2": &tSourceList{
			uint64(222),
			uint64(333),
		},
	}

	tests := []struct {
		name string
		hm   *tHashMap
		want tHashMap
	}{
		{"1", hm1, wm1},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hm.sort(); !got.compareTo(tt.want) {
				t.Errorf("%q: tHashMap.sort() = \n%v\n>>>> want: >>>>\n%v",
					tt.name, got, tt.want)
			}
		})
	}
} // Test_tHashMap_sort()

func Test_tHashMap_store(t *testing.T) {
	fn := hmFilename()
	hm1 := prepHashMap().
		add("@alphons", 1)

	tests := []struct {
		name      string
		hm        *tHashMap
		useBinary bool
		wantInt   int
		wantErr   bool
	}{
		{"1", hm1, false, 84, false},
		{"2", hm1, true, 109, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		fName := fn + tt.name
		UseBinaryStorage = tt.useBinary
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.hm.store(fName)
			if (nil != err) != tt.wantErr {
				t.Errorf("%q: tHashMap.store() error = %v, wantErr %v",
					tt.name, err, tt.wantErr)
				return
			}
			if got != tt.wantInt {
				t.Errorf("%q: tHashMap.store() = %v, want %v",
					tt.name, got, tt.wantInt)
			}
		})
		os.Remove(fName)
	}
} // Test_tHashMap_store()

func Test_tHashMap_String(t *testing.T) {
	sl0 := tHashMap{}
	ws0 := ""
	sl1 := prepHashMap()
	ws1 := "[#hash1]\n111\n[#hash2]\n222\n[#hash3]\n333\n[@mention1]\n111\n[@mention2]\n222\n"

	tests := []struct {
		name string
		hm   tHashMap
		want string
	}{
		{"0", sl0, ws0},
		{"1", *sl1, ws1},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hm.String(); got != tt.want {
				t.Errorf("%q: tHashMap.String() = \n%v\n>>>> want: >>>>\n%v",
					tt.name, got, tt.want)
			}
		})
	}
} // Test_tHashMap_String()

func Test_tHashMap_walk(t *testing.T) {
	type args struct {
		aFunc TWalkFunc
	}
	tests := []struct {
		name      string
		hm        *tHashMap
		args      args
		wantRBool bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRBool := tt.hm.walk(tt.args.aFunc); gotRBool != tt.wantRBool {
				t.Errorf("tHashMap.walk() = %v, want %v", gotRBool, tt.wantRBool)
			}
		})
	}
}

func Benchmark_LoadTxT(b *testing.B) {
	saveBinary := UseBinaryStorage
	defer func() {
		os.Remove(testHmStore)
		UseBinaryStorage = saveBinary
	}()
	hm := prepHashMap().
		add("@CrashTestDummy", 1)
	UseBinaryStorage = false

	hm.store(testHmStore)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		if _, err := hm.Load(testHmStore); nil != err {
			log.Printf("LoadTxt(): %v", err)
		}
	}
} // Benchmark_LoadTxt()

func Benchmark_LoadBin(b *testing.B) {
	saveBinary := UseBinaryStorage
	defer func() {
		os.Remove(testHmStore)
		UseBinaryStorage = saveBinary
	}()
	hm := prepHashMap().
		add("@CrashTestDummy", 1)
	UseBinaryStorage = true

	hm.store(testHmStore)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		if _, err := hm.Load(testHmStore); nil != err {
			log.Printf("LoadBin(): %v", err)
		}
	}
} // Benchmark_LoadBin()

func Benchmark_StoreTxt(b *testing.B) {
	saveBinary := UseBinaryStorage
	defer func() {
		os.Remove(testHmStore)
		UseBinaryStorage = saveBinary
	}()
	hm := prepHashMap().
		add("@CrashTestDummy", 1)
	UseBinaryStorage = false

	for n := 0; n < b.N; n++ {
		if _, err := hm.store(testHmStore); nil != err {
			log.Printf("StoreTxt(): %v", err)
		}
	}
} // Benchmark_StoreTxt()

func Benchmark_StoreBin(b *testing.B) {
	saveBinary := UseBinaryStorage
	defer func() {
		os.Remove(testHmStore)
		UseBinaryStorage = saveBinary
	}()
	hm := prepHashMap().
		add("@CrashTestDummy", 1)
	UseBinaryStorage = true

	for n := 0; n < b.N; n++ {
		if _, err := hm.store(testHmStore); nil != err {
			log.Printf("StoreBin(): %v", err)
		}
	}
} // Benchmark_StoreBin()
