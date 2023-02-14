package node

type Groupable struct {
	groupBy *string
}

func (s *Groupable) groupByValue(curr map[string]string) string {
	if s.groupBy == nil {
		return "total"
	}
	groupByValue, ok := curr[*s.groupBy]
	if !ok {
		panic("invalid groupBy in count node")
	}
	return groupByValue
}
