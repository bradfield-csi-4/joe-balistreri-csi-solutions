package main

import (
	"fmt"
	"strconv"

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
	// // d.PredeterminedLevels = []int{0, 1, 0, 0, 2, 1, 2, 3}
	// d.Put([]byte("hello"), []byte("goodbye"))
	// d.Put([]byte("apple"), []byte("pie"))
	// sl.Put([]byte("aaa"), []byte("cake"))
	// sl.Put([]byte("bvb"), []byte("cake"))
	// sl.Put([]byte("ababaaa"), []byte("cake"))
	// sl.Put([]byte("zebra"), []byte("cake"))
	// sl.Put([]byte("chewie"), []byte("boy"))
	// sl.Put([]byte("apple"), []byte("cry"))
	// sl.Put([]byte("apple"), []byte("cry"))

	// sl.Print()

	d, done := db.NewKVStore("example1")
	// d := db.NewSkipList()
	defer done()

	v := []byte("stringbean")

	for i := 0; i < 500; i++ {
		d.Put(db.KeyFromIterator(i), append(v, []byte(strconv.Itoa(i))...))
	}

	v, err := d.Get(db.KeyFromIterator(10))
	fmt.Println(string(v))
	fmt.Println(err)

	// d.Put([]byte("key2"), []byte("value2"))
	// d.Put([]byte("key1"), []byte("value3"))
	// d.Put([]byte("key1"), []byte("1"))
	// d.Put([]byte("jiuce"), []byte("juice"))
	// d.Put([]byte("juicy"), []byte("juice"))
	// d.Put([]byte("change"), []byte("juice"))
	// d.Put([]byte("bang"), []byte("juice"))
	// d.Put([]byte("alpha"), []byte("juice"))
	// d.Put([]byte("jiuce"), []byte("juice"))
	// d.Put([]byte("crazy"), []byte("juice"))
	// d.Put([]byte("casdfasdfrazy"), []byte("juice"))
	// d.Put([]byte("asdf"), []byte("juice"))

	// i, _ := d.RangeScan(nil, nil)
	// for i.Next() {
	// 	fmt.Println(string(i.Key()), string(i.Value()))
	// }
}
