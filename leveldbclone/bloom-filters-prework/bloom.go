package main

import (
	"encoding/binary"

	"github.com/spaolacci/murmur3"
)

type bloomFilter interface {
	add(item string)

	// `false` means the item is definitely not in the set
	// `true` means the item might be in the set
	maybeContains(item string) bool

	// Number of bytes used in any underlying storage
	memoryUsage() int
}

type trivialBloomFilter struct {
	data []uint64
}

func newTrivialBloomFilter() *trivialBloomFilter {
	return &trivialBloomFilter{
		data: make([]uint64, 1000),
	}
}

func (b *trivialBloomFilter) add(item string) {
	// Do nothing
}

func (b *trivialBloomFilter) maybeContains(item string) bool {
	// Technically, any item "might" be in the set
	return true
}

func (b *trivialBloomFilter) memoryUsage() int {
	return binary.Size(b.data)
}

type legitBloomFilter struct {
	kHashes      int
	mBits        uint64
	data         []uint64
	bitsPerIndex uint64
}

func newLegitBloomFilter(kHashes, mBits, nElementsExpected int) *legitBloomFilter {
	return &legitBloomFilter{
		kHashes:      kHashes,
		mBits:        uint64(mBits),
		data:         make([]uint64, mBits/64),
		bitsPerIndex: 64,
	}
}

var salts = []string{
	"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p",
}

func (b *legitBloomFilter) add(item string) {
	for i := 0; i < b.kHashes; i++ {
		bit := murmur3.Sum64([]byte(item+salts[i])) % b.mBits
		b.data[bit/b.bitsPerIndex] = b.data[bit/b.bitsPerIndex] | 1<<(bit%b.bitsPerIndex)
	}
}

func (b *legitBloomFilter) maybeContains(item string) bool {
	for i := 0; i < b.kHashes; i++ {
		bit := murmur3.Sum64([]byte(item+salts[i])) % b.mBits
		var mask uint64 = 1 << (bit % b.bitsPerIndex)
		if b.data[bit/b.bitsPerIndex]&mask != mask {
			return false
		}
	}
	return true
}

func (b *legitBloomFilter) memoryUsage() int {
	return binary.Size(b.data)
}
