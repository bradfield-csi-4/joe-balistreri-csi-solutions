package node

func ParseNode(q QueryExpression, nextNode ExecutionNode) ExecutionNode {
	switch q.Name {
	case "SCAN":
		if len(q.Args) != 1 {
			panic("invalid args for scan node")
		}
		return NewScanNode(q.Args[0])
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
