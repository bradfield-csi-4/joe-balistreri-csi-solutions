package db

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/smartystreets/assertions/should"
)

func TestSSTable(t *testing.T) {
	t.Run("After flushing the Memtable to an SSTable, an early value is found in the SSTable but not in the Memtable", func(t *testing.T) {
		kv, done := NewKVStore("test")

		v := []byte("stringbean")

		for i := 0; i < 500; i++ {
			kv.Put(KeyFromIterator(i), append(v, []byte(strconv.Itoa(i))...))
		}

		v, err := kv.Get(KeyFromIterator(10))
		So(t, should.BeNil(err))
		So(t, should.Resemble(string(v), "stringbean10"))

		v, err = kv.(*KVStore).memtable.Get(KeyFromIterator(10))
		So(t, should.BeNil(v))
		So(t, should.NotBeNil(err))

		for i := 0; i < 500; i++ {
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

		for i := 0; i < 500; i++ {
			v, err := kv2.Get(KeyFromIterator(i))
			So(t, should.BeNil(err))
			So(t, should.Resemble(string(v), fmt.Sprintf("stringbean%d", i)))
		}
	})

	t.Run("RangeScan works for an SSTable", func(t *testing.T) {
		kv, done := NewKVStore("atest")
		defer done()

		v := []byte("stringbean")

		// using 250 because it's large enough that only 1 sstable is flushed - we'll
		for i := 0; i < 250; i++ {
			kv.Put(KeyFromIterator(i), append(v, []byte(strconv.Itoa(i))...))
		}

		// full range
		it, err := kv.(*KVStore).SSTables[0].RangeScan(nil, nil)
		So(t, should.BeNil(err))
		start, last, count := runIterator(t, it, func(_ *testing.T, _, _ []byte) {})
		So(t, should.Resemble(start, []byte{0, 0, 0, 0}))
		So(t, should.Resemble(last, []byte{0, 0, 0, 247}))
		So(t, should.Equal(count, 248))

		// restricted start range
		it, err = kv.(*KVStore).SSTables[0].RangeScan([]byte{0, 0, 0, 99}, nil)
		So(t, should.BeNil(err))
		start, last, count = runIterator(t, it, func(_ *testing.T, _, _ []byte) {})
		So(t, should.Resemble(start, []byte{0, 0, 0, 99}))
		So(t, should.Resemble(last, []byte{0, 0, 0, 247}))
		So(t, should.Equal(count, 149))

		// restricted end range
		it, err = kv.(*KVStore).SSTables[0].RangeScan(nil, []byte{0, 0, 0, 100})
		So(t, should.BeNil(err))
		start, last, count = runIterator(t, it, func(_ *testing.T, _, _ []byte) {})
		So(t, should.Resemble(start, []byte{0, 0, 0, 0}))
		So(t, should.Resemble(last, []byte{0, 0, 0, 99}))
		So(t, should.Equal(count, 100))

		// restricted start and end range
		it, err = kv.(*KVStore).SSTables[0].RangeScan([]byte{0, 0, 0, 50}, []byte{0, 0, 0, 100})
		So(t, should.BeNil(err))
		start, last, count = runIterator(t, it, func(_ *testing.T, _, _ []byte) {})
		So(t, should.Resemble(start, []byte{0, 0, 0, 50}))
		So(t, should.Resemble(last, []byte{0, 0, 0, 99}))
		So(t, should.Equal(count, 50))

		// empty range
		it, err = kv.(*KVStore).SSTables[0].RangeScan([]byte{0, 0, 0, 50}, []byte{0, 0, 0, 50})
		So(t, should.BeNil(err))
		start, last, count = runIterator(t, it, func(_ *testing.T, _, _ []byte) {})
		So(t, should.Resemble(start, []byte(nil)))
		So(t, should.Resemble(last, []byte(nil)))
		So(t, should.Equal(count, 0))

		// invalid range
		it, err = kv.(*KVStore).SSTables[0].RangeScan([]byte{0, 0, 0, 100}, []byte{0, 0, 0, 50})
		So(t, should.BeNil(err))
		start, last, count = runIterator(t, it, func(_ *testing.T, _, _ []byte) {})
		So(t, should.Resemble(start, []byte(nil)))
		So(t, should.Resemble(last, []byte(nil)))
		So(t, should.Equal(count, 0))
	})

	t.Run("RangeScan works multiple SSTables", func(t *testing.T) {
		kv, done := NewKVStore("zztest")
		defer done()

		v := []byte("stringbean")
		for i := 0; i < 1000; i++ {
			kv.Put(KeyFromIterator(i), append(v, []byte(strconv.Itoa(i))...))
		}

		v2 := []byte("overwrite")
		for i := 300; i < 600; i++ {
			kv.Put(KeyFromIterator(i), append(v2, []byte(strconv.Itoa(i))...))
		}

		for i := 500; i < 700; i++ {
			kv.Delete(KeyFromIterator(i))
		}

		for i := 0; i < 100; i++ {
			kv.Put(KeyFromIterator(i), append(v, []byte(strconv.Itoa(i))...))
		}

		testResult := func(subT *testing.T, key, val []byte) {
			subT.Helper()
			if compareBytes(KeyFromIterator(299), key) == -1 && compareBytes(KeyFromIterator(599), key) == 1 {
				if !(strings.HasPrefix(string(v), "overwrite")) {
					panic("fuck")
				}
			} else {
				if !(strings.HasPrefix(string(v), "stringbean")) {
					panic("me")
				}
			}
		}

		// full range
		it, err := kv.RangeScan(nil, nil)
		So(t, should.BeNil(err))
		start, last, count := runIterator(t, it, testResult)
		So(t, should.Resemble(start, []byte{0, 0, 0, 0}))
		So(t, should.Resemble(last, KeyFromIterator(999)))
		So(t, should.Equal(count, 800))

		// restricted start range
		it, err = kv.RangeScan([]byte{0, 0, 0, 99}, nil)
		So(t, should.BeNil(err))
		start, last, count = runIterator(t, it, testResult)
		So(t, should.Resemble(start, []byte{0, 0, 0, 99}))
		So(t, should.Resemble(last, KeyFromIterator(999)))
		So(t, should.Equal(count, 700))

		// restricted end range
		it, err = kv.RangeScan(nil, []byte{0, 0, 0, 100})
		So(t, should.BeNil(err))
		start, last, count = runIterator(t, it, testResult)
		So(t, should.Resemble(start, []byte{0, 0, 0, 0}))
		So(t, should.Resemble(last, []byte{0, 0, 0, 99}))
		So(t, should.Equal(count, 100))

		// restricted start and end range
		it, err = kv.RangeScan([]byte{0, 0, 0, 50}, []byte{0, 0, 0, 100})
		So(t, should.BeNil(err))
		start, last, count = runIterator(t, it, testResult)
		So(t, should.Resemble(start, []byte{0, 0, 0, 50}))
		So(t, should.Resemble(last, []byte{0, 0, 0, 99}))
		So(t, should.Equal(count, 50))

		// empty range
		it, err = kv.RangeScan([]byte{0, 0, 0, 50}, []byte{0, 0, 0, 50})
		So(t, should.BeNil(err))
		start, last, count = runIterator(t, it, testResult)
		So(t, should.Resemble(start, []byte(nil)))
		So(t, should.Resemble(last, []byte(nil)))
		So(t, should.Equal(count, 0))

		// invalid range
		it, err = kv.RangeScan([]byte{0, 0, 0, 100}, []byte{0, 0, 0, 50})
		So(t, should.BeNil(err))
		start, last, count = runIterator(t, it, testResult)
		So(t, should.Resemble(start, []byte(nil)))
		So(t, should.Resemble(last, []byte(nil)))
		So(t, should.Equal(count, 0))
	})

	t.Run("Correctly handles missing values", func(t *testing.T) {
		kv, done := NewKVStore("btest")
		defer done()

		v := []byte("stringbean")

		// only write even values
		for i := 10; i < 1000; i++ {
			if i%2 == 0 {
				kv.Put(KeyFromIterator(i), append(v, []byte(strconv.Itoa(i))...))
			}
		}

		for i := 1; i < 1020; i++ {
			v, err := kv.Get(KeyFromIterator(i))
			hres, hErr := kv.Has(KeyFromIterator(i))
			if i%2 == 0 && i >= 10 && i < 1000 {
				So(t, should.BeNil(err))
				So(t, should.Resemble(string(v), fmt.Sprintf("stringbean%d", i)))
				So(t, should.BeNil(hErr))
				So(t, should.BeTrue(hres))
			} else {
				So(t, should.NotBeNil(err))
				So(t, should.Equal(err, ErrNotFound))
				So(t, should.BeNil(v))
				So(t, should.Equal(hErr, ErrNotFound))
				So(t, should.BeFalse(hres))
			}
		}
	})
	t.Run("Correctly handles overwrites", func(t *testing.T) {
		kv, done := NewKVStore("ctest")
		defer done()

		v := []byte("stringbean")

		for i := 0; i < 500; i++ {
			kv.Put(KeyFromIterator(i), append(v, []byte(strconv.Itoa(i))...))
		}

		v2 := []byte("helloworld")
		for i := 100; i < 400; i++ {
			kv.Put(KeyFromIterator(i), append(v2, []byte(strconv.Itoa(i))...))
		}

		for i := 0; i < 500; i++ {
			v, err := kv.Get(KeyFromIterator(i))
			if i >= 100 && i < 400 {
				So(t, should.BeNil(err))
				So(t, should.Resemble(string(v), fmt.Sprintf("helloworld%d", i)))
			} else {
				So(t, should.BeNil(err))
				So(t, should.Resemble(string(v), fmt.Sprintf("stringbean%d", i)))
			}
		}
	})

	t.Run("Correctly handles deletes", func(t *testing.T) {
		kv, done := NewKVStore("dtest")
		defer done()

		v := []byte("stringbean")

		for i := 0; i < 500; i++ {
			kv.Put(KeyFromIterator(i), append(v, []byte(strconv.Itoa(i))...))
		}

		for i := 100; i < 400; i++ {
			kv.Delete(KeyFromIterator(i))
		}

		for i := 0; i < 500; i++ {
			if i >= 100 && i < 400 {
				v, err := kv.Get(KeyFromIterator(i))
				So(t, should.Equal(err, ErrKeyDeleted))
				So(t, should.BeNil(v))
			} else {
				v, err := kv.Get(KeyFromIterator(i))
				So(t, should.BeNil(err))
				So(t, should.Resemble(string(v), fmt.Sprintf("stringbean%d", i)))
			}
		}
	})
	// TODO: do pre-work for monday's class
	// - implement RangeScan for KVStore with multiple SSTables

	// TODO: leveled compaction - level 0 is temporary; to move down a level, compact an SSTable with all tables it intersects in the subsequent level; each level deeper should have 10x more data

	// TODO: read bigtable paper

	// TODO: increase the number of values written in the test such that we write multiple ssTables
	// TODO: flush can write to multiple files if it exceeds the limit of a single SSTable file
	// TODO: actual limits: 2MB for SSTable file, 2 KB for each chunk in the index

	// BONUS: handle nil value separately from deleted
	// BONUS: make flush an async operation - can have two memtables during the flush - one taking reads and one frozen
	// BONUS: how to compress the ssTable files? and still do random io?
	// BONUS: add the ability to do snapshots
	// BONUS: store the largest key in each index entry in order to reduce an extra read; can also add the first and last key in the table so we can avoid scanning a block where we don't have something
	// BONUS: could move the index length to the back of the file (this lets you start writing the file before you know everything about it)
	// BONUS: replace gob for encoding the index
}

func runIterator(t *testing.T, it Iterator, checkResult func(*testing.T, []byte, []byte)) ([]byte, []byte, int) {
	t.Helper()
	var start, last []byte
	var count int
	for it.Next() {
		if compareBytes(it.Key(), []byte{0, 0, 0, 247}) == 0 {
			fmt.Println("hello")
		}
		// checkResult(t, it.Key(), it.Value())
		if start == nil {
			start = it.Key()
		}
		last = it.Key()
		count++
	}
	return start, last, count
}
