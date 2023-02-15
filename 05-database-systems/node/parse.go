package node

func ParseNode(q QueryExpression, nextNode ExecutionNode) ExecutionNode {
	switch q.Name {
	case "SCAN":
		if len(q.Args) != 1 {
			panic("invalid args for scan node")
		}
		return NewScanNode(q.Args[0])
	case "LIMIT":
		if len(q.Args) != 1 {
			panic("invalid args for limit node")
		}
		return NewLimitNode(q.Args[0], nextNode)
	case "SORT":
		if len(q.Args) == 0 {
			panic("invalid args for limit node")
		}
		desc := false
		if len(q.Args) == 2 && q.Args[1] == "DESC" {
			desc = true
		}
		return NewSortNode(q.Args[0], desc, nextNode)
	case "DISTINCT":
		if len(q.Args) != 1 {
			panic("invalid args for distinct node")
		}
		return NewDistinctNode(q.Args[0], nextNode)
	case "AGG":
		if len(q.Args) == 0 {
			panic("invalid args for agg")
		}
		var groupBy *string
		if len(q.Args) >= 3 {
			if q.Args[len(q.Args)-2] == "GROUP BY" {
				groupBy = &q.Args[len(q.Args)-1]
			}
		}
		var aggs []Aggregator
	Loop:
		for i := 0; i < len(q.Args); {
			switch q.Args[i] {
			case "COUNT":
				aggs = append(aggs, NewCountAggregator(groupBy))
				i++
			case "AVG":
				aggs = append(aggs, NewAvgAggregator(q.Args[i+1], groupBy, false))
				i += 2
			case "SUM":
				aggs = append(aggs, NewAvgAggregator(q.Args[i+1], groupBy, true))
				i += 2
			default:
				break Loop
			}
		}

		return NewAggregatorNode(nextNode, AggOptions{Aggregators: aggs})
	case "SELECTION":
		if len(q.Args) != 3 {
			panic("invalid args for selection node")
		}
		return NewSelectionNode(q.Args[0], q.Args[1], q.Args[2], nextNode)
	case "PROJECTION":
		return NewProjectionNode(q.Args, nextNode)
	}
	panic("unknown node type")
}
