/*
Copyright © 2019, 2025  M.Watermann, 10247 Berlin, Germany

	    All rights reserved
	EMail : <support@mwat.de>
*/
package hashtags

//lint:file-ignore ST1017 - I prefer Yoda conditions

// func htFilename() string {
// 	fn := filepath.Join(os.TempDir(), "testHlStore.db")
// 	os.Remove(fn) // remove remnants of previous runs
// 	return fn
// } // htFilename()

// func Test_newHashList(t *testing.T) {
// 	fn1 := htFilename() + "1"
// 	fn2 := htFilename() + "2"
// 	defer func() {
// 		os.Remove(fn1)
// 		os.Remove(fn2)
// 	}()
// 	hl1, _ := newHashList(fn1)
// 	hl2, _ := newHashList(fn2)
// 	tests := []struct {
// 		tName   string
// 		fName   string
// 		want    *tHashList
// 		wantErr bool
// 	}{
// 		{" 1", fn1, hl1, false},
// 		{" 2", fn2, hl2, false},
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.tName, func(t *testing.T) {
// 			got, err := newHashList(tt.fName)
// 			if (nil != err) != tt.wantErr {
// 				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("%q: New() =\n{%v}\n>>>> want: >>>>\n{%v}",
// 					tt.tName, got, tt.want)
// 			}
// 		})
// 	}
// } // Test_newHashList()

// func Test_THashList_equals(t *testing.T) {
// 	hl1 := &tHashList{
// 		hm: make(tHashMap, 64),
// 	}
// 	wl1 := &tHashList{
// 		hm: make(tHashMap, 64),
// 	}
// 	wl2 := &tHashList{
// 		hm: make(tHashMap, 64),
// 	}
// 	wl2.insert(MarkHash, "hash2", 2222)
// 	wl2.insert(MarkMention, "Name", 2222)
// 	tests := []struct {
// 		name string
// 		hl   *tHashList
// 		list *tHashList
// 		want bool
// 	}{
// 		{"1", hl1, wl1, true},
// 		{"2", hl1, wl2, false},
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := tt.hl.hm.equals(tt.list.hm); got != tt.want {
// 				t.Errorf("THashList.equals() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// } // Test_THashList_equals()

// func Test_THashList_insert(t *testing.T) {
// 	hl0 := &tHashList{
// 		hm: make(tHashMap, 64),
// 	}
// 	type tArgs struct {
// 		aDelim byte
// 		aName  string
// 		aID    int64
// 	}
// 	tests := []struct {
// 		name string
// 		hl   *tHashList
// 		args tArgs
// 		want bool
// 	}{
// 		{"0", hl0, tArgs{}, false},
// 		{"1", hl0, tArgs{MarkHash, "hash1", 1}, true},
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		got := tt.hl.insert(tt.args.aDelim, tt.args.aName, tt.args.aID)
// 		if got != tt.want {
// 			t.Errorf("%q: THashList.insert() =\n{%v}\n>>>> want: >>>>\n{%v}",
// 				tt.name, got, tt.want)
// 		}
// 	}
// } // Test_THashList_insert()

// func funcHashMentionRE(aText string) int {
// 	matches := htHashMentionRE.FindAllStringSubmatch(aText, -1)
// 	println(fmt.Sprintf("%s: %v", aText, matches))
// 	return len(matches)
// } // funcHashMentionRE()

// func Test_htHashMentionRE(t *testing.T) {
// 	t1 := `1blabla #HÄSCH1 blabla #hash2. Blabla`
// 	t2 := `2blabla #hash2. Blabla "#hash3" blabla`
// 	t3 := `\n>#KurzErklärt #Zensurheberrecht verhindern\n`
// 	t4 := `4blabla **#HÄSCH1** blabla\n\n_#hash3_`
// 	t5 := `5blabla&#39; **#hash2** blabla\n<a href="page#hash3">txt</a> #hash4`
// 	t6 := `#hash3 blabla\n<a href="https://www.tagesspiegel.de/politik/martin-sonneborn-wirbt-fuer-moralische-integritaet-warum-ich-die-eu-kommission-ablehnen-werde/25263366.html#25263366">txt</a> #hash4`
// 	t7 := `2blabla #hash2. @Dale_O'Leary "#hash3" @Dale_O’Leary blabla @Henry's`
// 	tests := []struct {
// 		name string
// 		text string
// 		want int
// 	}{
// 		{" 1", t1, 2},
// 		{" 2", t2, 2},
// 		{" 3", t3, 2},
// 		{" 4", t4, 2},
// 		{" 5", t5, 4},
// 		{" 6", t6, 3},
// 		{" 7", t7, 5},
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := funcHashMentionRE(tt.text); got != tt.want {
// 				t.Errorf("%q: funcHashMentionRE() = %v, want %v",
// 					tt.name, got, tt.want)
// 			}
// 		})
// 	}
// } // Test_htHashMentionRE()

// func Test_THashList_parseID(t *testing.T) {
// 	hash1, hash2, hash3, hash4 := "#HÄSCH1", "#hash2", "#hash3", "#hash4"
// 	hyphTx1, hyphTx2, hyphTx3 := `#--------------`, `#---text ---`, `#-text-`
// 	id1, id2, id3, id4, id5, id6 := int64(987), int64(654), int64(321), int64(123), int64(456), int64(789)
// 	hl1, _ := newHashList("")
// 	tx1 := []byte("1blabla " + hash1 + " blabla " + hash3 + ". Blabla")
// 	hl2 := &tHashList{
// 		hm: tHashMap{
// 			hash1: &tSourceList{id3},
// 			hash2: &tSourceList{id3},
// 			hash3: &tSourceList{id3},
// 		},
// 	}
// 	tx2 := []byte(`2blabla "` + hash2 + `". Blabla ` + hash3 + ` blabla`)
// 	hl3, _ := newHashList("")
// 	tx3 := []byte("3\n> #KurzErklärt #Zensurheberrecht verhindern -\n> [Glyphosat-Gutachten selbst anfragen!](https://fragdenstaat.de/aktionen/zensurheberrecht-2019/)\n")
// 	hl4, _ := newHashList("")
// 	tx4 := []byte("4blabla **" + hash1 + "** blabla\n\n_" + hash3 + "_")
// 	hl5, _ := newHashList("")
// 	tx5 := []byte(`5blabla&#39; **` + hash2 + `** blabla\n<a href="page#fragment">txt</a> ` + hash4)
// 	hl6, _ := newHashList("")
// 	tx6 := []byte(hash3 + ` blabla\n<a href="https://www.tagesspiegel.de/politik/martin-sonneborn-wirbt-fuer-moralische-integritaet-warum-ich-die-eu-kommission-ablehnen-werde/25263366.html#25263366">txt</a> ` + hash4)
// 	hl7, _ := newHashList("")
// 	tx7 := []byte(`7 (https://www.faz.net/aktuell/politik/inland/jutta-ditfurth-zu-extinction-rebellion-irrationalismus-einer-endzeit-sekte-16422668.html?printPagedArticle=true#ageIndex_2)`)
// hl8, _ := newHashList("")
// tx8 := []byte(`8
// > [Here's Everything You Need To Know](https://thehackernews.com/2018/12/australia-anti-encryption-bill.html#content) by <writer@example.com>
// `)
// hl9, _ := newHashList("")
// tx9 := []byte(`9
// Bla *@Antoni_Comín* bla bla _#§219a_
// `)
// tmp := string(tx6) + "\n" + hyphTx1 + ` and ` + hyphTx2 + "\n" + hyphTx3
// tx10 := []byte(tmp)
// tx11 := []byte{}
// tx12 := []byte(" ")
// type tArgs struct {
// 	aID   int64
// 	aText []byte
// }
// tests := []struct {
// 	name string
// 	hl   *tHashList
// 	args tArgs
// 	want bool
// 	}{
// 		{"1", hl1, tArgs{id1, tx1}, true},
// 		{"2", hl2, tArgs{id2, tx2}, true},
// 		{"3", hl3, tArgs{id3, tx3}, true},
// 		{"4", hl4, tArgs{id4, tx4}, true},
// 		{"5", hl5, tArgs{id5, tx5}, true},
// 		{"6", hl6, tArgs{id6, tx6}, true},
// 		{"7", hl7, tArgs{7, tx7}, false},
// 		{"8", hl8, tArgs{8, tx8}, false},
// 		{"9", hl9, tArgs{9, tx9}, true},
// 		{"10", hl6, tArgs{id6, tx10}, false},
// 		{"11", hl1, tArgs{id1, tx11}, false},
// 		{"12", hl1, tArgs{id1, tx12}, false},
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			hl := &tHashList{
// 				hm: tt.hl.hm,
// 			}
// 			if got := hl.parseID(tt.args.aID, tt.args.aText); got != tt.want {
// 				t.Errorf("%q: THashList.parseID() = \n%v\n>>>> want >>>>\n%v\n{%s}",
// 					tt.name, got, tt.want, hl)
// 			}
// 		})
// 	}
// } // Test_THashList_parseID()

// func Test_THashList_updateID(t *testing.T) {
// 	hash1, hash2, hash3 := "hash1", "hash2", "hash3"
// 	id1, id2, id3 := int64(987), int64(654), int64(321)
// 	hl, _ := newHashList("")
// 	hl.insert(MarkHash, hash1, id1)
// 	hl.insert(MarkHash, hash2, id1)
// 	hl.insert(MarkHash, hash3, id3)
// 	hl.insert(MarkMention, hash1, id3)
// 	hl.insert(MarkMention, hash2, id2)
// 	hl.insert(MarkMention, hash3, id2)
// 	tx0 := []byte("not recognised: support@mwat.de; accepted: <support@dfg> doesn't")
// 	tx1 := []byte("blabla #" + hash1 + " blabla @" + hash3 + " blabla")
// 	tx2 := []byte("blabla @" + hash1 + " blabla #" + hash3 + " blabla")
// 	type tArgs struct {
// 		aID   int64
// 		aText []byte
// 	}
// 	tests := []struct {
// 		name string
// 		args tArgs
// 		want bool
// 	}{
// 		{"0", tArgs{0, tx0}, true},
// 		{"1", tArgs{id1, tx1}, true},
// 		{"2", tArgs{id2, tx2}, true},
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got := hl.updateID(tt.args.aID, tt.args.aText)
// 			if got != tt.want {
// 				t.Errorf("%q: THashList.IDupdate() =\n%v\n>>>> want: >>>>\n%v\n%v",
// 					tt.name, got, tt.want, hl)
// 			}
// 		})
// 	}
// } // Test_THashList_updateID()

/* EoF */
