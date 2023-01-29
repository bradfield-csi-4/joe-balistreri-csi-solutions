package main

import (
	"fmt"

	"github.com/jdbalistreri/bradfield-csi-solutions/leveldbclone/db"
)

func main() {
	d, done := db.NewKVStore("example")
	defer done()

	// v := []byte("stringbean")

	// for i := 0; i < 500; i++ {
	// 	err := d.Put(db.KeyFromIterator(i), append(v, []byte(strconv.Itoa(i))...))
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }

	for i := 0; i < 500; i++ {
		v, _ := d.Get(db.KeyFromIterator(i))
		fmt.Printf("%d: %s\n", i, string(v))
	}

	// it, err := d.RangeScan(db.KeyFromIterator(0), db.KeyFromIterator(500))
	// fmt.Println(err)
	// fmt.Println(it)
	// it.Next()
	// for i := 0; it.Next() && i < 700; i++ {
	// 	fmt.Printf("%v: %s\n", it.Key(), strconv.Quote(string(it.Value())))
	// }
}
