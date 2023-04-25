package node

type DistinctNode struct {
	ColumnNameParser
	field      string
	seen       map[string]bool
	underlying ExecutionNode
}

func NewDistinctNode(field string, underlying ExecutionNode) *DistinctNode {
	return &DistinctNode{
		underlying: underlying,
		field:      field,
		seen:       map[string]bool{},
	}
}

func (s *DistinctNode) Next() Row {
	if !s.initialized {
		columns := s.underlying.Next()
		s.AddColumns(columns)
		return columns
	}

	curr := s.underlying.Next()
	if curr == nil {
		return nil
	}

	v := curr[s.columnsToIndex[s.field]]
	for s.seen[v] {
		curr := s.underlying.Next()
		if curr == nil {
			return nil
		}
		v = curr[s.columnsToIndex[s.field]]
	}

	s.seen[v] = true
	return curr
}
