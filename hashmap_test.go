/*
Copyright © 2023, 2025  M.Watermann, 10247 Berlin, Germany

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

var (
	exts = map[bool]string{
		true:  ".gob",
		false: ".lst",
	}
)

func hmFilename(useBinary bool) string {
	UseBinaryStorage = useBinary
	return "testHmStore" + exts[UseBinaryStorage]
} // hmFilename()

func prepHashMap() *tHashMap {
	hm := make(tHashMap, 8)
	// add already sorts keys
	hm.insert("#hash1", 111)
	hm.insert("#hash2", 222)
	hm.insert("#hash3", 333)
	hm.insert("@mention1", 111)
	hm.insert("@mention2", 222)

	return &hm
} // prepHashMap()

func Test_tHashMap_clear(t *testing.T) {
	hm1 := prepHashMap()
	wm1 := &tHashMap{}

	tests := []struct {
		name string
		hm   *tHashMap
		want *tHashMap
	}{
		{"1", hm1, wm1},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hm.clear(); !got.equals(*tt.want) {
				t.Errorf("%q: tHashMap.clear() =\n%v\n>>>> want >>>>\n%v",
					tt.name, got, tt.want)
			}
		})
	}
} // Test_tHashMap_clear()

func Test_tHashMap_count(t *testing.T) {
	hm1 := prepHashMap()
	hm2 := prepHashMap()
	hm2.insert("#hash4", 222)
	hm3 := prepHashMap()
	hm3.insert("@mention3", 222)

	tests := []struct {
		name    string
		hm      *tHashMap
		delim   byte
		wantInt int
	}{
		{"1", hm1, MarkHash, 3},
		{"2", hm2, MarkHash, 4},
		{"3", hm3, MarkMention, 3},
		{"4", hm1, 'x', 0},

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

	hm2 := prepHashMap()
	hm2.insert("@Alphons", 222)
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
		hm   *tHashMap
		want TCountList
	}{
		{"1", hm1, wc1},
		{"2", hm2, wc2},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hm.countedList(); !got.Equal(tt.want) {
				t.Errorf("%q: tHashMap.countedList() = \n%v\n>>>> want: >>>>\n%v",
					tt.name, got, tt.want)
			}
		})
	}
} // Test_tHashMap_countedList()

func Test_tHashMap_equals(t *testing.T) {
	hm1 := prepHashMap()
	om1 := prepHashMap()
	om2 := prepHashMap()
	om2.insert("#hash4", 222)

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
			if got := tt.hm.equals(*tt.oMap); got != tt.want {
				t.Errorf("%q: tHashMap.equals() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // Test_tHashMap_equals()

func Test_tHashMap_idList(t *testing.T) {
	hm := prepHashMap()

	id1 := int64(111)
	wl1 := []string{
		"#hash1",
		"@mention1",
	}

	id2 := int64(222)
	wl2 := []string{
		"#hash2",
		"@mention2",
	}

	tests := []struct {
		name string
		id   int64
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
	hm := prepHashMap()
	hm.insert("#hash2", 333)

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

func Test_tHashMap_insert(t *testing.T) {
	hm1 := prepHashMap()
	hm2 := prepHashMap()
	hm2.insert("#hash4", 444)

	type tArgs struct {
		aName string
		aID   int64
	}
	tests := []struct {
		name string
		hm   *tHashMap
		args tArgs
		want bool
	}{
		{"0", hm1, tArgs{"", 0}, false},            // empty hash
		{"1", hm1, tArgs{"@mention2", 222}, false}, // already added
		{"2", hm2, tArgs{"#hash4", 222}, true},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hm.insert(tt.args.aName, tt.args.aID); got != tt.want {
				t.Errorf("%q: tHashMap.insert() =\n%v\n>>>> want: >>>>\n%v",
					tt.name, got, tt.want)
			}
		})
	}
} // Test_tHashMap_insert()

func Test_tHashMap_keys(t *testing.T) {
	// hm.add("#hash1", 111).
	// 	add("#hash2", 222).
	// 	add("#hash3", 333).
	// 	add("@mention1", 111).
	// 	add("@mention2", 222)
	hm0 := &tHashMap{}
	wl0 := []string{}
	hm1 := prepHashMap()
	hm1.insert("#hash4", 444)
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
	hm := prepHashMap()
	hm.insert("#hash3", 33)
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
		{"1", tArgs{MarkHash, " "}, wl0},
		{"2", tArgs{MarkHash, "hash1"}, wl1},
		{"3", tArgs{MarkHash, "hash3"}, wl2},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotList := hm.list(tt.args.aDelim, tt.args.aName)
			if !tt.wantList.equals(gotList) {
				t.Errorf("%q: tHashMap.list() =\n%v\n>>>> want: >>>>\n%v",
					tt.name, gotList, tt.wantList)
			}
		})
	}
} // Test_tHashMap_list()

func Test_tHashMap_load(t *testing.T) {
	saveBinary := UseBinaryStorage
	defer func() {
		UseBinaryStorage = saveBinary
	}()

	hm1 := prepHashMap()
	hm1.insert("@CrashTestDummy", 1)

	wm1 := prepHashMap()
	wm1.insert("@CrashTestDummy", 1)

	tests := []struct {
		name    string
		hm      *tHashMap
		binary  bool
		want    *tHashMap
		wantErr bool
	}{
		{"1", hm1, false, wm1, false},
		{"2", hm1, true, wm1, false},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fName := hmFilename(tt.binary)
			// make sure there's actually data in the file:
			tt.hm.store(fName)

			got, err := tt.hm.load(fName)
			if (nil != err) != tt.wantErr {
				t.Errorf("%q: tHashMap.load() =\n%v\n>>>> want >>>>\n'%v'",
					tt.name, err, tt.wantErr)
				return
			}
			if !tt.want.equals(*got) {
				t.Errorf("%q: tHashMap.load() =\n%v\n>>>> want >>>>\n%v",
					tt.name, got, tt.want)
			}
			// fName = hmFilename(!tt.binary)
			// tt.hm.store(fName)
		})
	}
} // Test_tHashMap_load()

func Test_tHashMap_remove(t *testing.T) {
	hm := prepHashMap()
	// hm.add("#hash1", 111)
	// hm.add("#hash2", 222)
	// hm.add("#hash3", 333)
	// hm.add("@mention1", 111)
	// hm.add("@mention2", 222)
	hm.insert("#nameX", 999)
	hm.insert("#hash3", 333)

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
		{"0", tArgs{}, false},
		{"1", tArgs{MarkHash, "#nameX", 111}, false},
		{"2", tArgs{MarkHash, "#nameX", 999}, true},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hm.removeHM(tt.args.aDelim, tt.args.aName, tt.args.aID)
			if got != tt.want {
				t.Errorf("%q: tHashMap.remove() =\n%v\n>>>> want >>>>\n%v",
					tt.name, got, tt.want)
			}
		})
	}
} // Test_tHashMap_remove()

func Test_tHashMap_removeID(t *testing.T) {
	hm0 := newHashMap()

	hm1 := prepHashMap()
	hm1.insert("#hash2", 999)
	hm1.insert("#hash3", 888)
	hm1.insert("#CrashTestDummy", 777)

	tests := []struct {
		name string
		hm   *tHashMap
		id   int64
		want bool
	}{
		{"0", hm0, 999, false},
		{"1", hm1, 999, true},
		{"2", hm1, 888, true},
		{"3", hm1, 777, true},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hm.removeID(tt.id); got != tt.want {
				t.Errorf("%q: tHashMap.removeID() =\n%v\n>>>> want >>>>\n%v",
					tt.name, got, tt.want)
			}
		})
	}
} // Test_tHashMap_removeID()

func Test_tHashMap_renameID(t *testing.T) {
	id1, id2, id3 := int64(11), int64(22), int64(33)

	hm1 := prepHashMap()
	hm1.insert("#hash1", id1)
	hm1.insert("#hash1", id3)
	hm1.insert("@mention1", id3)

	hm3 := prepHashMap()
	hm3.insert("#hash3", id2)
	hm3.insert("#hash4", id2)

	type tArgs struct {
		aOldID, aNewID int64
	}
	tests := []struct {
		name string
		hm   *tHashMap
		args tArgs
		want bool
	}{
		{"0", hm1, tArgs{}, false}, // no change
		{"1", hm1, tArgs{id1, id2}, true},
		{"2", hm3, tArgs{id2, id3}, true},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hm.renameID(tt.args.aOldID, tt.args.aNewID); got != tt.want {
				t.Errorf("%q: tHashMap.renameID() =\n%v\n>>>> want: >>>>\n%v",
					tt.name, got, tt.want)
			}
		})
	}
} // Test_tHashMap_renameID

func Test_tHashMap_sort(t *testing.T) {
	hm1 := &tHashMap{
		"#hash1": &tSourceList{
			int64(111),
		},
		"@mention1": &tSourceList{
			int64(111),
		},
		"#hash2": &tSourceList{
			int64(222),
		},
		"@mention2": &tSourceList{
			int64(333),
			int64(222),
		},
		"#hash3": &tSourceList{
			int64(333),
		},
	}
	wm1 := &tHashMap{
		"#hash1": &tSourceList{
			int64(111),
		},
		"#hash2": &tSourceList{
			int64(222),
		},
		"#hash3": &tSourceList{
			int64(333),
		},
		"@mention1": &tSourceList{
			int64(111),
		},
		"@mention2": &tSourceList{
			int64(222),
			int64(333),
		},
	}

	tests := []struct {
		name string
		hm   *tHashMap
		want *tHashMap
	}{
		{"1", hm1, wm1},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hm.sort(); !got.equals(*tt.want) {
				t.Errorf("%q: tHashMap.sort() = \n%v\n>>>> want: >>>>\n%v",
					tt.name, got, tt.want)
			}
		})
	}
} // Test_tHashMap_sort()

func Test_tHashMap_store(t *testing.T) {
	saveBinary := UseBinaryStorage
	defer func() {
		UseBinaryStorage = saveBinary
	}()

	hm1 := prepHashMap()
	hm1.insert("@alphons", 1)

	tests := []struct {
		name    string
		hm      *tHashMap
		binary  bool
		wantInt int
		wantErr bool
	}{
		{"1", hm1, false, 164, false},
		{"2", hm1, true, 113, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		fName := hmFilename(tt.binary)
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
	sl0 := &tHashMap{}
	ws0 := ""

	sl1 := prepHashMap()
	ws1 := "[#hash1]\n000000000000006f\n[#hash2]\n00000000000000de\n[#hash3]\n000000000000014d\n[@mention1]\n000000000000006f\n[@mention2]\n00000000000000de\n"

	tests := []struct {
		name string
		hm   *tHashMap
		want string
	}{
		{"0", sl0, ws0},
		{"1", sl1, ws1},
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

// func Test_tHashMap_walk(t *testing.T) {
// 	type args struct {
// 		aFunc TWalkFunc
// 	}
// 	tests := []struct {
// 		name      string
// 		hm        *tHashMap
// 		args      args
// 		wantRBool bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if gotRBool := tt.hm.walk(tt.args.aFunc); gotRBool != tt.wantRBool {
// 				t.Errorf("tHashMap.walk() = %v, want %v", gotRBool, tt.wantRBool)
// 			}
// 		})
// 	}
// }

func Benchmark_LoadTxT(b *testing.B) {
	saveBinary := UseBinaryStorage
	defer func() {
		UseBinaryStorage = saveBinary
	}()
	fn := hmFilename(false)

	hm := prepHashMap()
	hm.insert("@CrashTestDummy", 1)

	hm.store(fn)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		if _, err := hm.load(fn); nil != err {
			log.Printf("LoadTxt(): %v", err)
		}
	}
} // Benchmark_LoadTxt()

func Benchmark_LoadBin(b *testing.B) {
	saveBinary := UseBinaryStorage
	defer func() {
		UseBinaryStorage = saveBinary
	}()
	fn := hmFilename(true)

	hm := prepHashMap()
	hm.insert("@CrashTestDummy", 1)

	hm.store(fn)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		if _, err := hm.load(fn); nil != err {
			log.Printf("LoadBin(): %v", err)
		}
	}
} // Benchmark_LoadBin()

func Benchmark_StoreTxt(b *testing.B) {
	saveBinary := UseBinaryStorage
	defer func() {
		UseBinaryStorage = saveBinary
	}()
	fn := hmFilename(false)

	hm := prepHashMap()
	hm.insert("@CrashTestDummy", 1)

	for n := 0; n < b.N; n++ {
		if _, err := hm.store(fn); nil != err {
			log.Printf("StoreTxt(): %v", err)
		}
	}
} // Benchmark_StoreTxt()

func Benchmark_StoreBin(b *testing.B) {
	saveBinary := UseBinaryStorage
	defer func() {
		UseBinaryStorage = saveBinary
	}()
	fn := hmFilename(true)

	hm := prepHashMap()
	hm.insert("@CrashTestDummy", 1)

	for n := 0; n < b.N; n++ {
		if _, err := hm.store(fn); nil != err {
			log.Printf("StoreBin(): %v", err)
		}
	}
} // Benchmark_StoreBin()
