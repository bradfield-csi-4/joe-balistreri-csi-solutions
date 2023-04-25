package node

type SelectionNode struct {
	ColumnNameParser
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
	if !s.initialized {
		columns := s.underlying.Next()
		s.AddColumns(columns)
		return columns
	}

	curr := s.underlying.Next()
	fieldI := s.columnsToIndex[s.field]
	for len(curr) != 0 {
		if curr[fieldI] == s.value {
			return curr
		}
		curr = s.underlying.Next()
	}
	return nil
}
