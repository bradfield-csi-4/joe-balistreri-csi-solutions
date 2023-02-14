package node

type QueryExpression struct {
	Name string
	Args []string
}

type Row map[string]string

type ExecutionNode interface {
	Next() Row // TODO: how to represent a row??; is a string pointer a bad idea?
	// Init()
	// Close()
	// Inputs(_, _)?
}
