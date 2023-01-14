package db

import "fmt"

const MAX_LEVEL = 16

type SkipList struct {
	head     []*Node
	maxLevel int
	level    int
}

type Node struct {
	key   []byte
	value []byte
	next  *Node
}

func NewSkipList(maxLevel int) SkipList {
	return SkipList{
		maxLevel: maxLevel,
		head:     make([]*Node, maxLevel),
	}
}

func (s *SkipList) newNode(key, value []byte, next *Node) *Node {
	return &Node{key: key, value: value, next: next}
}

func (s *SkipList) Get(key []byte) (value []byte, err error) {
	node, err := s.getNode(key)
	if err != nil || node == nil {
		return nil, err
	}
	return node.value, nil
}

func (s *SkipList) getNode(key []byte) (*Node, error) {
	level := s.level
	node := s.head[level]

	for ; level >= 0; level-- {
	InnerLoop:
		for node != nil {
			switch compareBytes(node.key, key) {
			case 0:
				return node, nil
			case -1:
				node = node.next
				continue
			case 1:
				level--
				break InnerLoop
			}
		}
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

func (s *SkipList) Delete(key []byte) error {
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

func (s *SkipList) Put(key, value []byte) error {
	if s.head[0] == nil {
		s.head[0] = s.newNode(key, value, nil)
		return nil
	}
	node := s.head[0]
	var prev *Node
	for node != nil {
		c := compareBytes(node.key, key)
		switch c {
		case -1:
			prev = node
			node = node.next
			continue
		case 1:
			newNode := s.newNode(key, value, node)
			if prev != nil {
				prev.next = newNode
			} else {
				s.head[0] = newNode
			}
			return nil
		case 0:
			node.value = value
			return nil
		default:
			panic(fmt.Sprintf("unexpected condition %d for %s and %s", c, node.key, key))
		}
	}
	prev.next = s.newNode(key, value, nil)
	return nil
}

func (s *SkipList) RangeScan(start, limit []byte) (Iterator, error) {
	// skipping for now
	return nil, nil
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
