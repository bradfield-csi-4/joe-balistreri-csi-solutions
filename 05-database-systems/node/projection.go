package node

type ProjectionNode struct {
	ColumnNameParser
	fields     []string
	underlying ExecutionNode
}

func NewProjectionNode(fields []string, underlying ExecutionNode) *ProjectionNode {
	for _, field := range fields {
		if !valid_fields[field] {
			panic("invalid field name")
		}
	}
	return &ProjectionNode{
		fields:     fields,
		underlying: underlying,
	}
}

func (s *ProjectionNode) Next() Row {
	if !s.initialized {
		columns := s.underlying.Next()
		s.AddColumns(columns)
		return s.selection(columns)
	}

	curr := s.underlying.Next()
	if curr == nil {
		return nil
	}
	return s.selection(curr)
}

func (s *ProjectionNode) selection(row Row) Row {
	result := Row{}
	for _, v := range s.fields {
		result = append(result, row[s.columnsToIndex[v]])
	}
	return result
}
