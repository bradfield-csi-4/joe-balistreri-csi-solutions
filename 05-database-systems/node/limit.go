package node

import "strconv"

type LimitNode struct {
	currCount  int
	limit      int
	underlying ExecutionNode
}

func NewLimitNode(limit string, underlying ExecutionNode) *LimitNode {
	l, err := strconv.Atoi(limit)
	if err != nil {
		panic(err)
	}
	return &LimitNode{
		limit:      l,
		underlying: underlying,
	}
}

func (s *LimitNode) Next() Row {
	if s.currCount >= s.limit {
		return nil
	}
	curr := s.underlying.Next()
	if curr == nil {
		return nil
	}
	s.currCount++
	return curr
}
