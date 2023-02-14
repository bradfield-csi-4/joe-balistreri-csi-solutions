package node

type DistinctNode struct {
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

func (s *DistinctNode) Next() map[string]string {
	curr := s.underlying.Next()
	if curr == nil {
		return nil
	}

	v, ok := curr[s.field]
	for !ok || s.seen[v] {
		curr := s.underlying.Next()
		if curr == nil {
			return nil
		}
		v, ok = curr[s.field]
	}

	s.seen[v] = true
	return curr
}
