package node

type QueryExpression struct {
	Name string
	Args []string
}

type ExecutionNode interface {
	Next() map[string]string // TODO: how to represent a row??; is a string pointer a bad idea?
	// Init()
	// Close()
	// Inputs(_, _)?
}
