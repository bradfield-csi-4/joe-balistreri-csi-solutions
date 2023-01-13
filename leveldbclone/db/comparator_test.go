package db

import (
	"fmt"
	"testing"

	"github.com/smartystreets/assertions/should"
)

func TestComparator(t *testing.T) {
	for _, params := range []struct {
		a, b     []byte
		expected int
	}{
		{[]byte{}, []byte{}, 0},
		{[]byte{}, []byte{1}, -1},
		{[]byte{1}, []byte{}, 1},
		{[]byte{1}, []byte{1}, 0},
		{[]byte{1}, []byte{2}, -1},
		{[]byte{2}, []byte{1}, 1},
		{[]byte("hello"), []byte("goodbye"), 1},
		{[]byte("Hello"), []byte("goodbye"), -1},
	} {
		t.Run(fmt.Sprintf("%s - %s", params.a, params.b), func(t *testing.T) {
			So(t, should.Equal(compareBytes(params.a, params.b), params.expected))
		})
	}

}
