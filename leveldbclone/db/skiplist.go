package db

import (
	"fmt"
	"math/rand"
)

const MAX_LEVEL = 16
const p = 1.0 / 3.0

type SkipList struct {
	head                []*Node
	maxLevel            int
	level               int
	PredeterminedLevels []int
	plIndex             int
}

type Node struct {
	key     []byte
	value   []byte
	next    []*Node
	special byte
}

const MIN_NODE = 1
const MAX_NODE = 2

func NewSkipList(maxLevel int) SkipList {
	s := SkipList{
		maxLevel: maxLevel,
	}
	levels := make([]*Node, maxLevel)
	min := s.newNode(nil, nil)
	min.special = MIN_NODE
	max := s.newNode(nil, nil)
	max.special = MAX_NODE
	for i := 0; i < maxLevel; i++ {
		levels[i] = min
		min.next[i] = max
	}
	s.head = levels
	return s
}

func (s *SkipList) newNode(key, value []byte) *Node {
	return &Node{key: key, value: value, next: make([]*Node, s.maxLevel)}
}

func (s *SkipList) randomLevel() int {
	if len(s.PredeterminedLevels) > s.plIndex {
		r := s.PredeterminedLevels[s.plIndex]
		s.plIndex++
		return r
	}
	var level int
	for rand.Float64() < p && level < s.maxLevel {
		level++
	}
	return level
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
				node = node.next[level]
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

func (s *SkipList) Print() {
	for i := s.level; i >= 0; i-- {
		if s.head[i] == nil {
			fmt.Println("nil")
			continue
		}
		node := s.head[i]
		for node != nil {
			fmt.Printf("(%s: %s) -> ", string(node.key), string(node.value))
			node = node.next[i]
		}
		fmt.Println()
	}
}

func (s *SkipList) Put(key, value []byte) error {
	level := s.level
	node := s.head[level]
	updates := make([]*Node, s.maxLevel)

	for ; level >= 0; level-- {
		for node.next[level] != nil && compareBytes(node.next[level].key, key) == -1 && node.next[level].special != MAX_NODE {
			node = node.next[level]
		}
		updates[level] = node
	}
	node = node.next[0]

	if compareBytes(node.key, key) == 0 {
		node.value = value
	} else {
		newLevel := s.randomLevel()
		if newLevel > s.level {
			for i := s.level + 1; i <= newLevel; i++ {
				updates[i] = s.head[i]
			}
			s.level = newLevel
		}
		newNode := s.newNode(key, value)
		for i := 0; i <= newLevel; i++ {
			newNode.next[i] =
				updates[i].next[i]
			updates[i].next[i] = newNode
		}
	}
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
