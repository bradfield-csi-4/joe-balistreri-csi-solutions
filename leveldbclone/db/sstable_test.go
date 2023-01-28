package db

import (
	"strconv"
	"testing"

	"github.com/smartystreets/assertions/should"
)

func TestSSTable(t *testing.T) {
	t.Run("After flushing the Memtable to an SSTable, an early value is found in the SSTable but not in the Memtable", func(t *testing.T) {
		kv, done := NewKVStore("test")
		defer done()

		v := []byte("stringbean")

		for i := 0; i < 250; i++ {
			kv.Put(KeyFromIterator(i), append(v, []byte(strconv.Itoa(i))...))
		}

		v, err := kv.Get(KeyFromIterator(10))
		So(t, should.BeNil(err))
		So(t, should.Resemble(string(v), "stringbean10"))

		v, err = kv.(*KVStore).memtable.Get(KeyFromIterator(10))
		So(t, should.BeNil(v))
		So(t, should.NotBeNil(err))
	})
}
