package encoding

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var r = regexp.MustCompile("^(get|GET|set|SET) ([a-zA-Z0-9]+)(?:=(.+))?$")

func CommandFromString(s string) (*Command, error) {
	parts := r.FindStringSubmatch(s)
	if len(parts) != 4 {
		return nil, errors.New("invalid input")
	}

	cmd := Command{
		Key: []byte(parts[2]),
	}

	op := strings.ToLower(parts[1])
	switch op {
	case "get":
		cmd.Op = GET
		if parts[3] != "" {
			return nil, errors.New("cannot set a value with GET")
		}
	case "set":
		cmd.Op = SET
		cmd.Value = []byte(parts[3])
	default:
		return nil, fmt.Errorf("invalid command %s", op)
	}

	return &cmd, nil
}

const VERSION = 1

func (c *Command) ToBinaryV1() []byte {
	buf := bytes.Buffer{}
	varIntBuf := make([]byte, binary.MaxVarintLen64)

	// encode version
	n := binary.PutUvarint(varIntBuf, VERSION)
	buf.Write(varIntBuf[:n])

	// encode operation
	n = binary.PutUvarint(varIntBuf, uint64(c.Op))
	buf.Write(varIntBuf[:n])

	// encode the number of args
	args := c.Args()
	n = binary.PutUvarint(varIntBuf, uint64(len(args)))
	buf.Write(varIntBuf[:n])

	// encode arguments
	for _, arg := range c.Args() {
		// encode length
		n = binary.PutUvarint(varIntBuf, uint64(len(arg)))
		buf.Write(varIntBuf[:n])
		// encode bytes
		buf.Write(arg)
	}
	return buf.Bytes()
}

func CommandFromBinary(b []byte) (*Command, error) {
	r := bytes.NewReader(b)

	// check version
	version, err := binary.ReadUvarint(r)
	if err != nil {
		return nil, err
	}
	if version != 1 {
		panic("unexpected version")
	}

	// get op
	cmd := Command{}
	op, err := binary.ReadUvarint(r)
	if err != nil {
		panic(err)
	}
	cmd.Op = OP(op)

	nArgs, err := binary.ReadUvarint(r)
	if err != nil {
		panic(err)
	}

	var args [][]byte

	for i := 0; i < int(nArgs); i++ {
		argLen, err := binary.ReadUvarint(r)
		if err != nil {
			panic(err)
		}
		if argLen == 0 {
			args = append(args, []byte(""))
			continue
		}
		next := make([]byte, argLen)
		n, err := r.Read(next)
		if err != nil {
			panic(err)
		}
		if n != len(next) {
			panic("mismatch!")
		}
		args = append(args, next)
	}

	return cmd.AddArgs(args), nil
}
