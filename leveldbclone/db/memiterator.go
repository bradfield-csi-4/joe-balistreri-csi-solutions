package db

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
	if len(m.keys) == 0 {
		return nil
	}
	return m.keys[m.i]
}

// Value returns the value of the current key/value pair, or nil if done.
func (m *MemIterator) Value() []byte {
	if len(m.values) == 0 {
		return nil
	}
	return m.values[m.i]
}
