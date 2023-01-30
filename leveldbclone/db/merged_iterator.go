package db

type MergedIterator struct {
	its          []Iterator
	currKey      []byte
	currValue    []byte
	currIterator Iterator
	exhausted    bool
}

func (m *MergedIterator) Next() bool {
	// TODO: switch this to using a heap instead of a linear scan

	if len(m.its) == 0 || m.exhausted {
		return false
	}

	// loop through all iterators and make sure they're past the current key
	if m.currKey != nil {
		for _, it := range m.its {
			for it.Key() != nil && compareBytes(it.Key(), m.currKey) <= 0 && it.Next() {
			}
		}
	}

	// loop through all iterators to find the next min key and iterator
	var minKey []byte
	var minIterator Iterator
	for _, it := range m.its {
		if it.Key() != nil && (minKey == nil || compareBytes(it.Key(), minKey) == -1) {
			minKey = it.Key()
			minIterator = it
		}
	}

	// if we haven't found a new key, mark the iterator as exhausted
	if minKey == nil || compareBytes(minKey, m.currKey) == 0 {
		m.exhausted = true
		return false
	}

	m.currIterator = minIterator
	m.currKey = minKey
	m.currValue = minIterator.Value()
	return true
}

func (m *MergedIterator) Error() error {
	return nil
}

func (m *MergedIterator) Key() []byte {
	return m.currKey
}

func (m *MergedIterator) Value() []byte {
	return m.currValue
}
