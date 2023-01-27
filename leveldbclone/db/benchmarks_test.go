package db

import (
	"math/rand"
	"testing"
)

const MAX_DB_SIZE = 1000

func BenchmarkFillSeqTest(b *testing.B) {
	mt, done := NewMemTable()
	defer done()
	v := []byte("World")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mt.Put(KeyFromIterator(i), v)
	}
}

func BenchmarkFillRandom(b *testing.B) {
	mt, done := NewMemTable()
	defer done()
	v := []byte("World")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mt.Put(KeyFromIterator(i), v)
	}
}

func BenchmarkOverwrite(b *testing.B) {
	mt, done := NewMemTable()
	defer done()
	v := []byte("World")
	for i := 0; i < MAX_DB_SIZE; i++ {
		mt.Put(KeyFromIterator(i), v)
	}
	v2 := []byte("Cadabra")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mt.Put(KeyFromIterator(rand.Intn(b.N)), v2)
	}
}

func BenchmarkDeleteSeq(b *testing.B) {
	mt, done := NewMemTable()
	defer done()
	v := []byte("World")
	for i := 0; i <= MAX_DB_SIZE; i++ {
		mt.Put(KeyFromIterator(i), v)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mt.Delete(KeyFromIterator(i))
	}
}

func BenchmarkDeleteRandom(b *testing.B) {
	mt, done := NewMemTable()
	defer done()
	v := []byte("World")
	for i := 0; i <= MAX_DB_SIZE; i++ {
		mt.Put(KeyFromIterator(i), v)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mt.Delete(KeyFromIterator(rand.Intn(b.N)))
	}
}

func BenchmarkReadSeq(b *testing.B) {
	mt, done := NewMemTable()
	defer done()
	v := []byte("World")
	for i := 0; i <= MAX_DB_SIZE; i++ {
		mt.Put(KeyFromIterator(i), v)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mt.Get(KeyFromIterator(i))
	}
}

func BenchmarkReadReverse(b *testing.B) {
	mt, done := NewMemTable()
	defer done()
	v := []byte("World")
	for i := 0; i <= MAX_DB_SIZE; i++ {
		mt.Put(KeyFromIterator(i), v)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mt.Get(KeyFromIterator(b.N - i))
	}
}

func BenchmarkReadRandom(b *testing.B) {
	mt, done := NewMemTable()
	defer done()
	v := []byte("World")
	for i := 0; i <= MAX_DB_SIZE; i++ {
		mt.Put(KeyFromIterator(i), v)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mt.Get(KeyFromIterator(rand.Intn(b.N)))
	}
}

func BenchmarkRangeScanNoIteration(b *testing.B) {
	mt, done := NewMemTable()
	defer done()
	v := []byte("World")
	for i := 0; i <= MAX_DB_SIZE; i++ {
		mt.Put(KeyFromIterator(i), v)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mt.RangeScan(getRange(b.N))
	}
}

func BenchmarkRangeScanWithIteration(b *testing.B) {
	mt, done := NewMemTable()
	defer done()
	v := []byte("World")
	for i := 0; i <= MAX_DB_SIZE; i++ {
		mt.Put(KeyFromIterator(i), v)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		iterator, err := mt.RangeScan(getRange(b.N))
		if err != nil {
			panic(err)
		}
		for iterator.Next() {
			iterator.Key()
			iterator.Value()
		}
	}
}

func BenchmarkRangeAndPut(b *testing.B) {
	mt, done := NewMemTable()
	defer done()
	v := []byte("World")
	b.ResetTimer()
	for i := 0; i < MAX_DB_SIZE; i++ {
		mt.Put(KeyFromIterator(i), v)
		mt.RangeScan(getRange(b.N))
	}
}

func getRange(max int) ([]byte, []byte) {
	start, finish := rand.Intn(max), rand.Intn(max)
	if finish < start {
		temp := start
		start = finish
		finish = temp
	}
	return KeyFromIterator(start), KeyFromIterator(finish)
}

// Comma-separated list of operations to run in the specified order
//   Actual benchmarks:
//      fillseq       -- write N values in sequential key order in async mode
//      fillrandom    -- write N values in random key order in async mode
//      overwrite     -- overwrite N values in random key order in async mode
//      fillsync      -- write N/100 values in random key order in sync mode
//      fill100K      -- write N/1000 100K values in random order in async mode
//      deleteseq     -- delete N keys in sequential order
//      deleterandom  -- delete N keys in random order
//      readseq       -- read N times sequentially
//      readreverse   -- read N times in reverse order
//      readrandom    -- read N times in random order
//      readmissing   -- read N missing keys in random order
//      readhot       -- read N times in random order from 1% section of DB
//      seekrandom    -- N random seeks
//      seekordered   -- N ordered seeks
//      open          -- cost of opening a DB
//      crc32c        -- repeated crc32c of 4K of data
//   Meta operations:
//      compact     -- Compact the entire DB
//      stats       -- Print DB stats
//      sstables    -- Print sstable info
//      heapprofile -- Dump a heap profile (if supported by this port)
