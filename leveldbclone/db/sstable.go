package db

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"io"
	"os"
)

func LoadSSTable(f *os.File) *SSTable {
	// read the index offset
	indexOffsetBytes := make([]byte, 4)
	n, err := f.Read(indexOffsetBytes)
	if err != nil {
		panic(err)
	}
	if n != 4 {
		panic("read wrong number of bytes")
	}
	indexOffset := int(binary.LittleEndian.Uint32(indexOffsetBytes))

	f.Seek(int64(indexOffset), 0)

	b, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	var index []IndexEntry
	gob.NewDecoder(bytes.NewReader(b)).Decode(&index)
	if err != nil {
		panic(err)
	}

	return &SSTable{
		f:           f,
		indexOffset: indexOffset,
		index:       index,
	}
}

func (s *SSTable) Get(key []byte) (value []byte, err error) {
	// find the correct index entry
	idx := s.findIndexEntry(key)
	if idx == nil {
		return nil, &NotFoundError{}
	}
	// seek to its start position
	_, err = s.f.Seek(int64(idx.Offset), 0)
	if err != nil {
		panic(err)
	}

	// read the index into memory
	b := make([]byte, idx.Length)
	n, err := s.f.Read(b)
	if err != nil {
		panic(err)
	}
	if n != idx.Length {
		panic("read wrong number of bytes")
	}

	// decode the bytes and see if we can find the value
	_, currKey, currVal, err := readLogLine(bytes.NewReader(b))
	for compareBytes(currKey, key) == -1 && err == nil {
		_, currKey, currVal, err = readLogLine(bytes.NewReader(b))
	}
	if compareBytes(currKey, key) == 0 {
		return currVal, nil
	}
	if err != nil {
		return nil, err
	}
	return nil, &NotFoundError{}
}

func (s *SSTable) findIndexEntry(key []byte) *IndexEntry {
	for i, v := range s.index {
		if compareBytes(v.Key, key) > 1 { // we've gotten to an index greater than our key
			if i == 0 {
				return nil // if we're at the first index entry, every entry is greater than our key
			}
			return &s.index[i-1] // otherwise, return the index before the one that exceeded our key
		}
	}
	return &s.index[len(s.index)-1] // the start of the last index is less than the key, so we'll search it
}

func (s *SSTable) Has(key []byte) (ret bool, err error) {
	return false, nil
}

func (s *SSTable) RangeScan(start, limit []byte) (Iterator, error) {
	return nil, nil
}

func (s *SSTable) Close() {
	if s.f != nil {
		s.f.Close()
	}
}

type SSTable struct {
	f           *os.File
	indexOffset int
	index       []IndexEntry
}

type IndexEntry struct {
	Key    []byte
	Offset int
	Length int
}
