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
	return nil, nil
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
