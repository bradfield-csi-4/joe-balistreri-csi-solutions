package node

type ColumnNameParser struct {
	columns        Row
	columnsToIndex map[string]int
	initialized    bool
}

func (c *ColumnNameParser) AddColumns(row Row) {
	c.columns = row
	c.columnsToIndex = map[string]int{}
	for i, v := range row {
		c.columnsToIndex[v] = i
	}
	c.initialized = true
}
