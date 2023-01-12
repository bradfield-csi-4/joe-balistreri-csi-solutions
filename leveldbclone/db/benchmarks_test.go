package db

import (
	"math/rand"
	"testing"
)

const MAX_KEYS = 256 ^ 3 // 16 MB

func keyFromIterator(i int) []byte {
	return []byte{byte(i >> 24), byte(i >> 16), byte(i >> 8), byte(i)}
}

func BenchmarkFillSeqTest(b *testing.B) {
	mt := NewMemTable()
	v := []byte("World")
	for i := 0; i < b.N; i++ {
		mt.Put(keyFromIterator(i), v)
	}
}

func BenchmarkFillRandom(b *testing.B) {
	mt := NewMemTable()
	v := []byte("World")
	for i := 0; i < b.N; i++ {
		mt.Put(keyFromIterator(i), v)
	}
}

func BenchmarkOverwrite(b *testing.B) {
	mt := NewMemTable()
	v := []byte("World")
	for i := 0; i <= MAX_KEYS; i++ {
		mt.Put(keyFromIterator(i), v)
	}
	v2 := []byte("Cadabra")
	for i := 0; i < b.N; i++ {
		mt.Put(keyFromIterator(rand.Intn(MAX_KEYS)), v2)
	}
}

func BenchmarkDeleteSeq(b *testing.B) {
	mt := NewMemTable()
	v := []byte("World")
	for i := 0; i <= MAX_KEYS; i++ {
		mt.Put(keyFromIterator(i), v)
	}
	for i := 0; i < b.N; i++ {
		mt.Delete(keyFromIterator(i))
	}
}

func BenchmarkDeleteRandom(b *testing.B) {
	mt := NewMemTable()
	v := []byte("World")
	for i := 0; i <= MAX_KEYS; i++ {
		mt.Put(keyFromIterator(i), v)
	}
	for i := 0; i < b.N; i++ {
		mt.Delete(keyFromIterator(rand.Intn(MAX_KEYS)))
	}
}

func BenchmarkReadSeq(b *testing.B) {
	mt := NewMemTable()
	v := []byte("World")
	for i := 0; i <= MAX_KEYS; i++ {
		mt.Put(keyFromIterator(i), v)
	}
	for i := 0; i < b.N; i++ {
		mt.Get(keyFromIterator(i))
	}
}

func BenchmarkReadReverse(b *testing.B) {
	mt := NewMemTable()
	v := []byte("World")
	for i := 0; i <= MAX_KEYS; i++ {
		mt.Put(keyFromIterator(i), v)
	}
	for i := 0; i < b.N; i++ {
		mt.Get(keyFromIterator(MAX_KEYS - i))
	}
}

func BenchmarkReadRandom(b *testing.B) {
	mt := NewMemTable()
	v := []byte("World")
	for i := 0; i <= MAX_KEYS; i++ {
		mt.Put(keyFromIterator(i), v)
	}
	for i := 0; i < b.N; i++ {
		mt.Get(keyFromIterator(rand.Intn(MAX_KEYS)))
	}
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
