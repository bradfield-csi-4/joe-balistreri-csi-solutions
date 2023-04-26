package encoding

import (
	"testing"

	"github.com/smartystreets/assertions/should"
)

func So(t *testing.T, s string) {
	t.Helper()
	if s != "" {
		t.Fatal(s)
	}
}

func TestToAndFromBinary(t *testing.T) {
	t.Run("should work as a bijection", func(t *testing.T) {
		cmds := []*Command{
			{Op: GET, Key: []byte("abcdef")},
			{Op: SET, Key: []byte("abcdef"), Value: []byte("xyz")},
		}

		for _, cmd := range cmds {
			b := cmd.ToBinaryV1()

			cmd2, err := CommandFromBinary(b)
			So(t, should.BeNil(err))
			So(t, should.Resemble(cmd, cmd2))
		}
	})
}
