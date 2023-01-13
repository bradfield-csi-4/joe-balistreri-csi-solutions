package db

import "fmt"

type SkipList struct {
	head *Node
}

type Node struct {
	key   []byte
	value []byte
	next  *Node
}

func NewSkipList() SkipList {
	return SkipList{}
}

func (s *SkipList) Get(key []byte) (value []byte, err error) {
	node := s.head
	for node != nil {
		if compareBytes(node.key, key) == 0 {
			return node.value, nil
		}
		node = node.next
	}
	return nil, nil
}

func (s *SkipList) Has(key []byte) (ret bool, err error) {
	k, err := s.Get(key)
	if err != nil {
		return false, err
	}
	return k != nil, nil
}

func (s *SkipList) Put(key, value []byte) error {
	if s.head == nil {
		s.head = &Node{key: key, value: value}
		return nil
	}
	node := s.head
	var prev *Node
	for node != nil {
		c := compareBytes(node.key, key)
		switch c {
		case -1:
			prev = node
			node = node.next
			continue
		case 1:
			newNode := &Node{key: key, value: value, next: node}
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
	prev.next = &Node{key: key, value: value}
	return nil
}

type SkipListIterator struct {
}

func (m *SkipListIterator) Next() bool {
	return false
}

func (m *SkipListIterator) Error() error {
	return nil
}

func (m *SkipListIterator) Key() []byte {
	return nil
}

func (m *SkipListIterator) Value() []byte {
	return nil
}
