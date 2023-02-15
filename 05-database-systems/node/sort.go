package node

import (
	"sort"
	"strconv"
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
			a := s.sorted[i][s.sortBy]
			b := s.sorted[j][s.sortBy]
			af, aerr := strconv.ParseFloat(a, 64)
			bf, berr := strconv.ParseFloat(b, 64)
			if aerr == nil && berr == nil {
				return af < bf
			}
			return a < b
		}
		if s.desc {
			sortFunc = func(i, j int) bool {
				a := s.sorted[i][s.sortBy]
				b := s.sorted[j][s.sortBy]
				af, aerr := strconv.ParseFloat(a, 64)
				bf, berr := strconv.ParseFloat(b, 64)
				if aerr == nil && berr == nil {
					return af > bf
				}
				return a > b
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
