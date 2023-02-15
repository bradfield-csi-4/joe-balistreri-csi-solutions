package node

import (
	"fmt"
	"strconv"
)

const TOTAL = "__total"

type AggregatorNode struct {
	AggOptions
	underlying  ExecutionNode
	initialized bool
	results     []Row
	resultIndex int
}

type AggOptions struct {
	Aggregators []Aggregator
}

func NewAggregatorNode(underlying ExecutionNode, opts AggOptions) *AggregatorNode {
	res := &AggregatorNode{
		AggOptions: opts,
		underlying: underlying,
	}
	return res
}

func (s *AggregatorNode) Next() Row {
	if !s.initialized {
		for curr := s.underlying.Next(); curr != nil; curr = s.underlying.Next() {
			for _, a := range s.Aggregators {
				a.Add(curr)
			}
		}
		results := map[string]Row{}
		for _, a := range s.Aggregators {
			aResults := a.Results()
			for groupByField, row := range aResults {
				if results[groupByField] == nil {
					results[groupByField] = Row{}
				}
				for k, v := range row {
					results[groupByField][k] = v
				}
			}
		}
		for _, v := range results {
			s.results = append(s.results, v)
		}
		s.initialized = true
	}

	if s.resultIndex >= len(s.results) {
		return nil
	}
	s.resultIndex++
	return s.results[s.resultIndex-1]
}

// AGGREGATOR INTERFACE

type Aggregator interface {
	Add(Row)
	Results() map[string]Row
}

// COUNT AGGREGATOR

type CountAggregator struct {
	currCounts map[string]int
	groupBy    *string
}

func (c *CountAggregator) Add(row Row) {
	groupByValue := groupByValue(row, c.groupBy)
	c.currCounts[groupByValue]++
}

func (c *CountAggregator) Results() map[string]Row {
	results := map[string]Row{}
	valueName := "count"
	groupName := groupByField(c.groupBy)
	for k, c := range c.currCounts {
		row := Row{valueName: strconv.Itoa(c)}
		if k != TOTAL {
			row[groupName] = k
		}
		results[k] = row
	}
	return results
}

func NewCountAggregator(groupBy *string) *CountAggregator {
	return &CountAggregator{currCounts: map[string]int{}, groupBy: groupBy}
}

// AVG AGGREGATOR

type AvgAggregator struct {
	currCounts map[string]int
	currTotals map[string]float64
	groupBy    *string
	field      string
	useSum     bool
}

func (c *AvgAggregator) Add(row Row) {
	groupByValue := groupByValue(row, c.groupBy)
	c.currCounts[groupByValue]++

	fieldValue, err := strconv.ParseFloat(row[c.field], 64)
	if err != nil {
		panic(err)
	}
	c.currTotals[groupByValue] += fieldValue
}

func (a *AvgAggregator) Results() map[string]Row {
	results := map[string]Row{}
	valueName := fmt.Sprintf("avg(%s)", a.field)
	if a.useSum {
		valueName = fmt.Sprintf("sum(%s)", a.field)
	}
	groupName := groupByField(a.groupBy)

	for k, c := range a.currCounts {
		t := a.currTotals[k]
		avg := t / float64(c)

		// handle sum vs avg
		value := avg
		if a.useSum {
			value = t
		}

		row := Row{valueName: strconv.FormatFloat(value, 'f', -1, 64)}
		if k != TOTAL {
			row[groupName] = k
		}
		results[k] = row
	}
	return results
}

func NewAvgAggregator(field string, groupBy *string, useSum bool) *AvgAggregator {
	return &AvgAggregator{
		currCounts: map[string]int{},
		currTotals: map[string]float64{},
		groupBy:    groupBy,
		field:      field,
		useSum:     useSum,
	}
}

// HELPERS

func groupByValue(curr Row, groupBy *string) string {
	if groupBy == nil {
		return TOTAL
	}
	groupByValue, ok := curr[*groupBy]
	if !ok {
		panic("invalid groupBy in count node")
	}
	return groupByValue
}

func groupByField(groupBy *string) string {
	if groupBy == nil {
		return TOTAL
	}
	return *groupBy
}
