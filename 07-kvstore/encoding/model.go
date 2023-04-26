package encoding

type OP int

const (
	GET OP = iota + 1
	SET
	MSG
)

type Command struct {
	Op    OP
	Key   []byte
	Value []byte
}

func (c *Command) Args() [][]byte {
	switch c.Op {
	case GET:
		return [][]byte{c.Key}
	case SET:
		return [][]byte{c.Key, c.Value}
	case MSG:
		return [][]byte{c.Value}
	}
	panic("invalid op")
}

func (c *Command) AddArgs(args [][]byte) *Command {
	switch c.Op {
	case GET:
		c.Key = args[0]
		return c
	case SET:
		c.Key = args[0]
		c.Value = args[1]
		return c
	case MSG:
		c.Value = args[0]
		return c
	}
	panic("invalid op")
}
