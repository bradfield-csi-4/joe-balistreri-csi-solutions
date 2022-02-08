package db

import (
	"sort"
)

type MemTable struct {
	m map[string][]byte
}

func NewMemTable() DB {
	return &MemTable{
		m: map[string][]byte{},
	}
}

func (m *MemTable) Get(key []byte) (value []byte, err error) {
	v, ok := m.m[string(key)]
	if ok {
		return v, nil
	}
	return nil, &NotFoundError{}
}

func (m *MemTable) Has(key []byte) (ret bool, err error) {
	_, ok := m.m[string(key)]
	return ok, nil
}

func (m *MemTable) Put(key, value []byte) error {
	m.m[string(key)] = value
	return nil
}

func (m *MemTable) Delete(key []byte) error {
	return nil
}

func (m *MemTable) RangeScan(start, limit []byte) (Iterator, error) {
	var sortedPairs [][2][]byte
	for k, v := range m.m {
		sortedPairs = append(sortedPairs, [2][]byte{[]byte(k), v})
	}
	sort.SliceStable(sortedPairs, func(i, j int) bool {
		return string(sortedPairs[i][0]) < string(sortedPairs[j][0])
	})
	var keys [][]byte
	var values [][]byte
	for _, pair := range sortedPairs {
		keys = append(keys, pair[0])
		values = append(values, pair[1])
	}
	return &MemIterator{
		keys: keys,
		values: values,
	}, nil
}

type MemIterator struct {
	keys [][]byte
	values [][]byte
	i int
}

// Next moves the iterator to the next key/value pair.
// It returns false if the iterator is exhausted.
func (m *MemIterator) Next() bool {
	if m.i < len(m.keys) - 1 {
		m.i++
		return true
	}
	return false
}

// Error returns any accumulated error. Exhausting all the key/value pairs
// is not considered to be an error.
func (m *MemIterator) Error() error {
	return nil
}

// Key returns the key of the current key/value pair, or nil if done.
func (m *MemIterator) Key() []byte {
	return m.keys[m.i]
}

// Value returns the value of the current key/value pair, or nil if done.
func (m *MemIterator) Value() []byte {
	return m.values[m.i]
}