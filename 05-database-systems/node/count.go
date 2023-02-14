package node

import "strconv"

type CountNode struct {
	Groupable
	currCounts map[string]int
	returned   bool
	underlying ExecutionNode
}

func NewCountNode(underlying ExecutionNode, groupBy *string) *CountNode {
	res := &CountNode{
		Groupable:  Groupable{groupBy: groupBy},
		underlying: underlying,
		currCounts: map[string]int{},
	}
	return res
}

func (s *CountNode) Next() map[string]string {
	if s.returned {
		return nil
	}
	for curr := s.underlying.Next(); curr != nil; curr = s.underlying.Next() {
		groupByValue := s.groupByValue(curr)
		s.currCounts[groupByValue] += 1
	}
	s.returned = true

	result := map[string]string{}

	for k, v := range s.currCounts {
		result[k] = strconv.Itoa(v)
	}

	return result
}
