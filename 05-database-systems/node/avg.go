package node

import (
	"fmt"
	"strconv"
)

type AvgNode struct {
	Groupable
	field      string
	currCounts map[string]int
	currTotals map[string]float64
	returned   bool
	underlying ExecutionNode
}

func NewAvgNode(underlying ExecutionNode, field string, groupBy *string) *AvgNode {
	res := &AvgNode{
		field:      field,
		Groupable:  Groupable{groupBy: groupBy},
		underlying: underlying,
		currCounts: map[string]int{},
		currTotals: map[string]float64{},
	}
	return res
}

func (s *AvgNode) Next() map[string]string {
	if s.returned {
		return nil
	}
	for curr := s.underlying.Next(); curr != nil; curr = s.underlying.Next() {
		fieldValue, err := strconv.ParseFloat(curr[s.field], 64)
		if err != nil {
			panic(err)
		}
		groupByValue := s.groupByValue(curr)
		s.currCounts[groupByValue] += 1
		s.currTotals[groupByValue] += fieldValue
	}
	s.returned = true

	result := map[string]string{}

	for k, c := range s.currCounts {
		t := s.currTotals[k]
		avg := t / float64(c)
		result[k] = fmt.Sprintf("%f", avg)
	}

	return result
}
