package node

import (
	"sort"
)

type SortNode struct {
	sortBy      string
	sorted      []Row
	index       int
	underlying  ExecutionNode
	initialized bool
	desc        bool
}

func NewSortNode(sortBy string, desc bool, underlying ExecutionNode) *SortNode {
	return &SortNode{
		underlying: underlying,
		sortBy:     sortBy,
		desc:       desc,
	}
}

func (s *SortNode) Next() Row {
	if !s.initialized {
		for curr := s.underlying.Next(); curr != nil; curr = s.underlying.Next() {
			s.sorted = append(s.sorted, curr)
		}
		sortFunc := func(i, j int) bool {
			return s.sorted[i][s.sortBy] < s.sorted[j][s.sortBy]
		}
		if s.desc {
			sortFunc = func(i, j int) bool {
				return s.sorted[i][s.sortBy] > s.sorted[j][s.sortBy]
			}
		}
		sort.SliceStable(s.sorted, sortFunc)
	}

	if s.index >= len(s.sorted) {
		return nil
	}

	res := s.sorted[s.index]
	s.index++
	return res
}
