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
		return nil, ErrNotFound
	}
	bReader := readSegment(s.f, idx.Offset, idx.Length)

	_, currKey, currVal, err := readLogLine(bReader)
	for compareBytes(currKey, key) == -1 && err == nil {
		_, currKey, currVal, err = readLogLine(bReader)
	}
	if compareBytes(currKey, key) == 0 {
		return currVal, nil
	}
	if err == nil || err == io.EOF {
		return nil, ErrNotFound
	}
	return nil, err
}

func (s *SSTable) findIndexEntry(key []byte) *IndexEntry {
	idx := s.findIndexEntryIdx(key)
	if idx == -1 {
		return nil
	}
	return &s.index[idx]
}

func (s *SSTable) findIndexEntryIdx(key []byte) int {
	for i, v := range s.index {
		if compareBytes(v.Key, key) == 1 { // we've gotten to an index greater than our key
			if i == 0 {
				return -1 // if we're at the first index entry, every entry is greater than our key
			}
			return i - 1 // otherwise, return the index before the one that exceeded our key
		}
	}
	return len(s.index) - 1 // the start of the last index is less than the key, so we'll search it
}

func (s *SSTable) Has(key []byte) (ret bool, err error) {
	res, err := s.Get(key)
	if err != nil {
		return false, err
	}
	return res != nil, nil
}

func (s *SSTable) RangeScan(start, limit []byte) (Iterator, error) {
	// find the correct index entry
	indexIdx := s.findIndexEntryIdx(start)
	if indexIdx == -1 {
		indexIdx = 0
	}
	idx := s.index[indexIdx]

	f, err := os.Open(s.f.Name())
	if err != nil {
		return nil, err
	}

	return &SSTableIterator{
		limit:       limit,
		start:       start,
		currSegment: readSegment(f, idx.Offset, idx.Length),
		f:           f,
		indexes:     s.index,
		indexIdx:    indexIdx,
		moreToRead:  true,
	}, nil
}

func readSegment(f io.ReadSeeker, offset, length int) io.Reader {
	_, err := f.Seek(int64(offset), 0)
	if err != nil {
		panic(err)
	}

	// read the index into memory
	b := make([]byte, length)
	n, err := f.Read(b)
	if err != nil {
		panic(err)
	}
	if n != length {
		panic("read wrong number of bytes")
	}

	return bytes.NewReader(b)
}

type SSTableIterator struct {
	f           *os.File
	indexes     []IndexEntry
	indexIdx    int
	currSegment io.Reader
	limit       []byte
	start       []byte
	moreToRead  bool

	key   []byte
	value []byte
	err   error
}

func (m *SSTableIterator) Next() bool {
	if !m.moreToRead {
		m.err = nil
		m.key = nil
		m.value = nil
		return false
	}

	for ok, key, val, err := m.next(); ok; ok, key, val, err = m.next() {
		if err != nil {
			m.err = err
			m.moreToRead = false
			return false
		}
		if compareBytes(key, m.start) == -1 && m.start != nil {
			continue
		}
		if compareBytes(key, m.limit) == 1 && m.limit != nil {
			m.moreToRead = false
			return false
		}
		m.key = key
		m.value = val
		return true
	}

	m.moreToRead = false
	return false
}

func (m *SSTableIterator) next() (bool, []byte, []byte, error) {
	_, currKey, currVal, err := readLogLine(m.currSegment)
	if err == nil {
		return true, currKey, currVal, nil
	}
	if err != io.EOF {
		return false, nil, nil, err
	}

	m.indexIdx++
	if m.indexIdx >= len(m.indexes) {
		return false, nil, nil, nil
	}
	idx := m.indexes[m.indexIdx]
	m.currSegment = readSegment(m.f, idx.Offset, idx.Length)
	return m.next()
}

func (m *SSTableIterator) Error() error {
	return m.err
}

func (m *SSTableIterator) Key() []byte {
	return m.key
}

func (m *SSTableIterator) Value() []byte {
	return m.value
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
