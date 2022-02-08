package db

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
	return &MemIterator{}, nil
}

type MemIterator struct {

}

// Next moves the iterator to the next key/value pair.
// It returns false if the iterator is exhausted.
func (m *MemIterator) Next() bool {
	return true
}

// Error returns any accumulated error. Exhausting all the key/value pairs
// is not considered to be an error.
func (m *MemIterator) Error() error {
	return nil
}

// Key returns the key of the current key/value pair, or nil if done.
func (m *MemIterator) Key() []byte {
	return nil
}

// Value returns the value of the current key/value pair, or nil if done.
func (m *MemIterator) Value() []byte {
	return nil
}