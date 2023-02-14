package node

import "strconv"

type CountNode struct {
	currCount  int
	returned   bool
	underlying ExecutionNode
}

func NewCountNode(underlying ExecutionNode) *CountNode {
	return &CountNode{
		underlying: underlying,
	}
}

func (s *CountNode) Next() map[string]string {
	if s.returned {
		return nil
	}
	for s.underlying.Next() != nil {
		s.currCount += 1
	}
	s.returned = true
	return map[string]string{"count": strconv.Itoa(s.currCount)}
}
