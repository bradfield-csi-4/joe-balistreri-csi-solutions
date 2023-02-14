package node

type ProjectionNode struct {
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
	curr := s.underlying.Next()
	if curr == nil {
		return nil
	}
	result := Row{}
	for _, f := range s.fields {
		result[f] = curr[f]
	}
	return result
}
