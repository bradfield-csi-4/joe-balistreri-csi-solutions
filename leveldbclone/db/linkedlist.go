package db

import "fmt"

type LinkedList struct {
	head *LLNode
}

type LLNode struct {
	key   []byte
	value []byte
	next  *LLNode
}

func NewLinkedList() DB {
	return &LinkedList{}
}

func (s *LinkedList) Get(key []byte) (value []byte, err error) {
	node, err := s.getNode(key)
	if err != nil {
		return nil, err
	}
	if node == nil || node.value == nil {
		return nil, ErrNotFound
	}
	return node.value, nil
}

func (s *LinkedList) getStartNode(key []byte) (*LLNode, error) {
	node := s.head
	for node != nil && compareBytes(node.key, key) == -1 {
		node = node.next
	}
	return node, nil
}

func (s *LinkedList) getNode(key []byte) (*LLNode, error) {
	node := s.head
	for node != nil {
		if compareBytes(node.key, key) == 0 {
			return node, nil
		}
		node = node.next
	}
	return nil, nil
}

func (s *LinkedList) Has(key []byte) (ret bool, err error) {
	k, err := s.Get(key)
	if err != nil {
		if err == ErrNotFound {
			return false, nil
		}
		return false, err
	}
	return k != nil, nil
}

func (s *LinkedList) Delete(key []byte) error {
	node, err := s.getNode(key)
	if err != nil {
		return err
	}
	if node == nil {
		return nil
	}
	node.value = nil
	return nil
}

func (s *LinkedList) Put(key, value []byte) error {
	if s.head == nil {
		s.head = &LLNode{key: key, value: value}
		return nil
	}
	node := s.head
	var prev *LLNode
	for node != nil {
		c := compareBytes(node.key, key)
		switch c {
		case -1:
			prev = node
			node = node.next
			continue
		case 1:
			newNode := &LLNode{key: key, value: value, next: node}
			if prev != nil {
				prev.next = newNode
			} else {
				s.head = newNode
			}
			return nil
		case 0:
			node.value = value
			return nil
		default:
			panic(fmt.Sprintf("unexpected condition %d for %s and %s", c, node.key, key))
		}
	}
	prev.next = &LLNode{key: key, value: value}
	return nil
}

func (s *LinkedList) RangeScan(start, limit []byte) (Iterator, error) {
	node, _ := s.getStartNode(start)
	return &LinkedListIterator{node: node, limit: limit}, nil
}

type LinkedListIterator struct {
	node  *LLNode
	limit []byte
}

func (m *LinkedListIterator) Next() bool {
	// skip deleted nodes
	for m.node != nil && m.node.next != nil && m.node.next.value == nil {
		m.node = m.node.next
	}

	if m.node == nil || m.node.next == nil {
		m.node = nil
		return false
	}
	if compareBytes(m.node.next.key, m.limit) == 1 {
		m.node = nil
		return false
	}
	m.node = m.node.next
	return true
}

func (m *LinkedListIterator) Error() error {
	return nil
}

func (m *LinkedListIterator) Key() []byte {
	if m.node == nil {
		return nil
	}
	return m.node.key
}

func (m *LinkedListIterator) Value() []byte {
	if m.node == nil {
		return nil
	}
	return m.node.value
}
