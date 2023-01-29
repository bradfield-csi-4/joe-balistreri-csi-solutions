package db

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/smartystreets/assertions/should"
)

func TestSSTable(t *testing.T) {
	t.Run("After flushing the Memtable to an SSTable, an early value is found in the SSTable but not in the Memtable", func(t *testing.T) {
		kv, done := NewKVStore("test")

		v := []byte("stringbean")

		// using 250 because it's large enough that only 1 sstable is flushed - we'll
		for i := 0; i < 250; i++ {
			kv.Put(KeyFromIterator(i), append(v, []byte(strconv.Itoa(i))...))
		}

		v, err := kv.Get(KeyFromIterator(10))
		So(t, should.BeNil(err))
		So(t, should.Resemble(string(v), "stringbean10"))

		v, err = kv.(*KVStore).memtable.Get(KeyFromIterator(10))
		So(t, should.BeNil(v))
		So(t, should.NotBeNil(err))

		for i := 0; i < 250; i++ {
			v, err := kv.Get(KeyFromIterator(i))
			So(t, should.BeNil(err))
			So(t, should.Resemble(string(v), fmt.Sprintf("stringbean%d", i)))
		}

		done() // finish using the first instance of the db

		// check that the db can load from disk and WAL is properly truncated
		kv2, done2 := NewKVStore("test")
		defer done2()
		v, err = kv2.Get(KeyFromIterator(10))
		So(t, should.BeNil(err))
		So(t, should.Resemble(string(v), "stringbean10"))

		v, err = kv2.(*KVStore).memtable.Get(KeyFromIterator(10))
		So(t, should.BeNil(v))
		So(t, should.NotBeNil(err))

		for i := 0; i < 250; i++ {
			v, err := kv2.Get(KeyFromIterator(i))
			So(t, should.BeNil(err))
			So(t, should.Resemble(string(v), fmt.Sprintf("stringbean%d", i)))
		}
	})
	t.Run("Correctly handles missing values", func(t *testing.T) {
		kv, done := NewKVStore("test")
		defer done()

		v := []byte("stringbean")

		// only write even values
		for i := 10; i < 500; i++ {
			if i%2 == 0 {
				kv.Put(KeyFromIterator(i), append(v, []byte(strconv.Itoa(i))...))
			}
		}

		for i := 0; i < 520; i++ {
			t.Logf("testing %d", i)
			v, err := kv.Get(KeyFromIterator(i))
			if i%2 == 0 && i >= 10 && i < 500 {
				So(t, should.BeNil(err))
				So(t, should.Resemble(string(v), fmt.Sprintf("stringbean%d", i)))
			} else {
				So(t, should.NotBeNil(err))
				So(t, should.HaveSameTypeAs(err, &NotFoundError{}))
				So(t, should.BeNil(v))
			}
		}
	})
	t.Run("Correctly handles overwrites", func(t *testing.T) {
		kv, done := NewKVStore("test")
		defer done()

		v := []byte("stringbean")

		// using 250 because it's large enough that only 1 sstable is flushed - we'll
		for i := 0; i < 250; i++ {
			kv.Put(KeyFromIterator(i), append(v, []byte(strconv.Itoa(i))...))
		}

		v2 := []byte("helloworld")
		for i := 20; i < 40; i++ {
			kv.Put(KeyFromIterator(i), append(v2, []byte(strconv.Itoa(i))...))
		}

		for i := 0; i < 250; i++ {
			v, err := kv.Get(KeyFromIterator(i))
			if i >= 20 && i < 40 {
				So(t, should.BeNil(err))
				So(t, should.Resemble(string(v), fmt.Sprintf("helloworld%d", i)))
			} else {
				So(t, should.BeNil(err))
				So(t, should.Resemble(string(v), fmt.Sprintf("stringbean%d", i)))
			}
		}
	})
	// TODO: implement Has for SSTable and test
	// TODO: implement RangeScan for SSTable
	// TODO: switch flush to use flush(w io.Writer, it Iterator) - then we can reuse this for compaction with a merged iterator
	// TODO: use rangescan to read all values (in the test)

	// TODO: read bigtable paper

	// TODO: do pre-work for thursday's class

	// TODO:

	// TODO: increase the number of values written in the test such that we write multiple ssTables
	// TODO: flush can write to multiple files if it exceeds the limit of a single SSTable file
	// TODO: actual limits: 2MB for SSTable file, 2 KB for each chunk in the index

	// TODO: leveled compaction - level 0 is temporary; to move down a level, compact an SSTable with all tables it intersects in the subsequent level; each level deeper should have 10x more data

	// BONUS: make flush an async operation - can have two memtables during the flush - one taking reads and one frozen
	// BONUS: how to compress the ssTable files? and still do random io?
}
