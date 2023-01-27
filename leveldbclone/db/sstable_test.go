package db

import (
	"strconv"
	"testing"

	"github.com/smartystreets/assertions/should"
)

func TestSSTable(t *testing.T) {
	t.Run("Has, Put, and Get work as expected", func(t *testing.T) {
		kv, done := NewKVStore("test")
		defer done()

		v := []byte("stringbean")

		for i := 0; i < 11; i++ {
			kv.Put(KeyFromIterator(i), append(v, []byte(strconv.Itoa(i))...))
		}

		v, err := kv.Get(KeyFromIterator(10))
		So(t, should.BeNil(err))
		So(t, should.Resemble(string(v), "stringbean10"))

		kv.(*KVStore).memtable.Get(KeyFromIterator(10))
		So(t, should.NotBeNil(err))
		So(t, should.BeNil(v))
	})
}
