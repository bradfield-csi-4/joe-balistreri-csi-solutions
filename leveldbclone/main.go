package main

import (
	"github.com/jdbalistreri/bradfield-csi-solutions/leveldbclone/db"
)

func main() {
	// fmt.Println("hello!")
	// j := db.NewMemTableV1()
	// j.Put([]byte("Hello"), []byte("World"))
	// v, err := j.Get([]byte("Hello"))
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(string(v))

	// rand.Seed(time.Now().UnixNano())

	// sl := db.NewSkipList(db.MAX_LEVEL)
	// // sl.PredeterminedLevels = []int{0, 1, 0, 0, 2, 1, 2, 3}
	// sl.Put([]byte("hello"), []byte("goodbye"))
	// sl.Put([]byte("apple"), []byte("pie"))
	// sl.Put([]byte("aaa"), []byte("cake"))
	// sl.Put([]byte("bvb"), []byte("cake"))
	// sl.Put([]byte("ababaaa"), []byte("cake"))
	// sl.Put([]byte("zebra"), []byte("cake"))
	// sl.Put([]byte("chewie"), []byte("boy"))
	// sl.Put([]byte("apple"), []byte("cry"))
	// sl.Put([]byte("apple"), []byte("cry"))
	// sl.Put([]byte("jiuce"), []byte("juice"))
	// sl.Put([]byte("juicy"), []byte("juice"))
	// sl.Put([]byte("change"), []byte("juice"))
	// sl.Put([]byte("bang"), []byte("juice"))
	// sl.Put([]byte("alpha"), []byte("juice"))
	// sl.Put([]byte("jiuce"), []byte("juice"))
	// sl.Put([]byte("crazy"), []byte("juice"))
	// sl.Put([]byte("casdfasdfrazy"), []byte("juice"))
	// sl.Put([]byte("asdf"), []byte("juice"))

	// sl.Print()

	d := db.NewKVStore("example3")
	d.Put([]byte("key2"), []byte("value2"))
	d.Put([]byte("key1"), []byte("value3"))
	d.Put([]byte("key1"), []byte("1"))
}
