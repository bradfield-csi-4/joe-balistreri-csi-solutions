package db

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
)

const MAX_LEVEL = 24
const p = 1.0 / 3.0
const SSTABLE_INDEX_INCR_BYTES = 400 // 10% of the SSTable size

type SkipList struct {
	head                []*Node
	maxLevel            int
	level               int
	PredeterminedLevels []int
	plIndex             int
	sizeBytes           int
}

type Node struct {
	key     []byte
	value   []byte
	next    []*Node
	special byte
}

const MIN_NODE = 1
const MAX_NODE = 2

func NewSkipList() *SkipList {
	return newSkipList(MAX_LEVEL)
}

func newSkipList(maxLevel int) *SkipList {
	s := SkipList{
		maxLevel: maxLevel,
	}
	levels := make([]*Node, maxLevel+1)
	min := s.newNode(nil, nil)
	min.special = MIN_NODE
	max := s.newNode(nil, nil)
	max.special = MAX_NODE
	for i := 0; i < maxLevel; i++ {
		levels[i] = min
		min.next[i] = max
	}
	s.head = levels
	return &s
}

type indexEntry struct {
	key    string
	offset int
}

func (s *SkipList) flushSSTable(f *os.File) error {
	var fileContents []byte

	i, err := s.RangeScan(nil, nil)
	if err != nil {
		return err
	}

	var index []indexEntry
	var currBytes int
	for i.Next() {
		// append to file contents
		nextLine := toLogLine(i.Key(), i.Value())
		currBytes += len(i.Key()) + len(i.Value())
		fileContents = append(fileContents, nextLine...)

		// append to index
		if currBytes > SSTABLE_INDEX_INCR_BYTES {
			index = append(index, indexEntry{key: string(i.Key()), offset: len(fileContents)})
			currBytes = 0
		}
	}

	marshalledIndex, err := json.Marshal(index)
	if err != nil {
		return err
	}
	indexOffset := len(fileContents) + 4
	indexOffsetEncoded := make([]byte, 4)
	binary.LittleEndian.PutUint32(indexOffsetEncoded, uint32(indexOffset))
	fileContents = append(indexOffsetEncoded, fileContents...)
	fileContents = append(fileContents, marshalledIndex...)

	// write to disc and flush
	n, err := f.Write(fileContents)
	if err != nil {
		return err
	}
	if n != len(fileContents) {
		return fmt.Errorf("wrote %d bytes but expected to write %d", n, len(fileContents))
	}
	err = f.Sync()
	if err == nil {
		return err
	}
	return nil
}

func (s *SkipList) newNode(key, value []byte) *Node {
	return &Node{key: key, value: value, next: make([]*Node, s.maxLevel+1)}
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
	if err != nil {
		return nil, err
	}
	if node == nil || node.value == nil {
		return nil, &NotFoundError{}
	}
	return node.value, nil
}

func (s *SkipList) getStart(key []byte) (*Node, error) {
	level := s.level
	node := s.head[level]

	for ; level >= 0; level-- {
		for node.next[level] != nil && compareBytes(node.next[level].key, key) == -1 && node.next[level].special != MAX_NODE {
			node = node.next[level]
		}
	}
	return node.next[0], nil
}

func (s *SkipList) getNode(key []byte) (*Node, error) {
	level := s.level
	node := s.head[level]

	for ; level >= 0; level-- {
	InnerLoop:
		for node != nil && node.next[level] != nil {
			switch compareBytes(node.next[level].key, key) {
			case 0:
				return node.next[level], nil
			case -1:
				node = node.next[level]
				continue
			case 1:
				break InnerLoop
			}
		}
	}
	return nil, nil
}

func (s *SkipList) Has(key []byte) (ret bool, err error) {
	k, err := s.Get(key)
	if err != nil {
		if _, ok := err.(*NotFoundError); ok {
			return false, nil
		}
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
	s.sizeBytes -= len(node.value)
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
		s.sizeBytes += len(value) - len(node.value)
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
		s.sizeBytes += len(key) + len(value)
	}
	return nil
}

func (s *SkipList) SizeBytes() int {
	return s.sizeBytes
}

func (s *SkipList) RangeScan(start, limit []byte) (Iterator, error) {
	if start == nil && limit == nil {
		return &SkipListIterator{node: s.head[0], limit: nil}, nil
	}
	node, err := s.getStart(start)
	if err != nil {
		return nil, err
	}
	return &SkipListIterator{node: node, limit: limit}, nil
}

type SkipListIterator struct {
	node  *Node
	limit []byte
}

func (m *SkipListIterator) Next() bool {
	// skip deleted nodes
	for m.node != nil && m.node.next[0] != nil && m.node.next[0].value == nil {
		m.node = m.node.next[0]
	}

	if m.node == nil || m.node.next[0] == nil {
		m.node = nil
		return false
	}
	if m.limit != nil && compareBytes(m.node.next[0].key, m.limit) == 1 {
		m.node = nil
		return false
	}
	m.node = m.node.next[0]
	return true
}

func (m *SkipListIterator) Error() error {
	return nil
}

func (m *SkipListIterator) Key() []byte {
	if m.node == nil {
		return nil
	}
	return m.node.key
}

func (m *SkipListIterator) Value() []byte {
	if m.node == nil {
		return nil
	}
	return m.node.value
}
