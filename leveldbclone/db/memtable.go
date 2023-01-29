package db

import (
	"sort"
)

type MemTableV1 struct {
	m map[string][]byte
}

func NewMemTableV1() DB {
	return &MemTableV1{
		m: map[string][]byte{},
	}
}

// Get - O(1)
func (m *MemTableV1) Get(key []byte) (value []byte, err error) {
	v, ok := m.m[string(key)]
	if ok {
		return v, nil
	}
	return nil, ErrNotFound
}

// Has - O(1)
func (m *MemTableV1) Has(key []byte) (ret bool, err error) {
	_, ok := m.m[string(key)]
	return ok, nil
}

// Put - O(1)
func (m *MemTableV1) Put(key, value []byte) error {
	m.m[string(key)] = value
	return nil
}

// Delete - O(1)
func (m *MemTableV1) Delete(key []byte) error {
	delete(m.m, string(key))
	return nil
}

// RangeScan - O(Nlog(N)) (or worse, with the string comparison for the range handling)
func (m *MemTableV1) RangeScan(start, limit []byte) (Iterator, error) {
	var sortedPairs [][2][]byte
	for k, v := range m.m {
		sortedPairs = append(sortedPairs, [2][]byte{[]byte(k), v})
	}
	sort.SliceStable(sortedPairs, func(i, j int) bool {
		return string(sortedPairs[i][0]) < string(sortedPairs[j][0])
	})
	var keys [][]byte
	var values [][]byte
	sstring := string(start)
	lstring := string(limit)
	// TODO: could do binary search here instead
	for _, pair := range sortedPairs {
		if string(pair[0]) < sstring {
			continue
		}
		if string(pair[0]) > lstring {
			break
		}
		keys = append(keys, pair[0])
		values = append(values, pair[1])
	}
	return &MemIterator{
		keys:   keys,
		values: values,
	}, nil
}
