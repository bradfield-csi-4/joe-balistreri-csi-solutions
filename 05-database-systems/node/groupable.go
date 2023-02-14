package node

type Groupable struct {
	groupBy *string
}

func (s *Groupable) groupByValue(curr Row) string {
	if s.groupBy == nil {
		return "total"
	}
	groupByValue, ok := curr[*s.groupBy]
	if !ok {
		panic("invalid groupBy in count node")
	}
	return groupByValue
}

func (s *Groupable) groupByField() string {
	if s.groupBy == nil {
		return "total"
	}
	return *s.groupBy
}
