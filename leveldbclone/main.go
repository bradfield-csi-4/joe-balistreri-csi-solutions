package main

import (
	"fmt"

	"github.com/jdbalistreri/bradfield-csi-solutions/leveldbclone/db"
)

func main() {
	fmt.Println("hello!")
	j := db.NewMemTableV1()
	j.Put([]byte("Hello"), []byte("World"))
	v, err := j.Get([]byte("Hello"))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(v))
}
