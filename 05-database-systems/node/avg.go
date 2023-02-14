package node

import (
	"fmt"
	"strconv"
)

type AvgNode struct {
	Groupable
	field       string
	underlying  ExecutionNode
	initialized bool
	results     []Row
	resultIndex int
}

func NewAvgNode(underlying ExecutionNode, field string, groupBy *string) *AvgNode {
	res := &AvgNode{
		field:      field,
		Groupable:  Groupable{groupBy: groupBy},
		underlying: underlying,
	}
	return res
}

func (s *AvgNode) Next() Row {
	if !s.initialized {
		currCounts := map[string]int{}
		currTotals := map[string]float64{}

		for curr := s.underlying.Next(); curr != nil; curr = s.underlying.Next() {
			fieldValue, err := strconv.ParseFloat(curr[s.field], 64)
			if err != nil {
				panic(err)
			}
			groupByValue := s.groupByValue(curr)
			currCounts[groupByValue]++
			currTotals[groupByValue] += fieldValue
		}
		valueName := fmt.Sprintf("avg(%s)", s.field)
		groupName := s.groupByField()
		for k, c := range currCounts {
			t := currTotals[k]
			avg := t / float64(c)
			s.results = append(s.results, Row{
				groupName: k,
				valueName: fmt.Sprintf("%f", avg),
			})
		}
	}

	if s.resultIndex >= len(s.results) {
		return nil
	}
	s.resultIndex++
	return s.results[s.resultIndex-1]
}
