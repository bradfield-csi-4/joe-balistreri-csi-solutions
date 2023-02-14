package node

import "fmt"

type SelectionNode struct {
	field      string
	operation  string
	value      string
	underlying ExecutionNode
}

var valid_fields = map[string]bool{
	"genres":  true,
	"movieId": true,
	"title":   true,
}

func NewSelectionNode(field, operation, value string, underlying ExecutionNode) *SelectionNode {
	if operation != "EQUALS" {
		panic("unsupported operation")
	}
	if !valid_fields[field] {
		panic("invalid field name")
	}
	return &SelectionNode{
		field:      field,
		operation:  operation,
		value:      value,
		underlying: underlying,
	}
}

func (s *SelectionNode) Next() Row {
	curr := s.underlying.Next()
	for len(curr) != 0 {
		v, ok := curr[s.field]
		if !ok {
			panic(fmt.Sprintf("invalid row, missing field %s", s.field))
		}
		if v == s.value {
			return curr
		}
		curr = s.underlying.Next()
	}
	return nil
}
