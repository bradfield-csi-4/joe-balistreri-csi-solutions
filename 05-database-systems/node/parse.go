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
	case "COUNT":
		var field *string
		if len(q.Args) == 2 {
			if q.Args[0] != "GROUP BY" {
				panic("invalid args for count")
			}
			field = &q.Args[1]
		}
		return NewAggregatorNode(nextNode, AggOptions{Aggregators: []Aggregator{NewCountAggregator(field)}})
	case "AVG":
		if len(q.Args) == 0 {
			panic("invalid args for avg")
		}
		var groupBy *string
		if len(q.Args) == 3 {
			if q.Args[1] != "GROUP BY" {
				panic("invalid args for avg")
			}
			groupBy = &q.Args[2]
		}
		return NewAggregatorNode(nextNode, AggOptions{Aggregators: []Aggregator{NewAvgAggregator(q.Args[0], groupBy, false)}})
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
