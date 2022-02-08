package db

import (
	"testing"

	"github.com/smartystreets/assertions/should"
)

func So(t *testing.T, s string) {
	if s != "" {
		t.Fatal(s)
	}
}

func TestMemTable(t *testing.T) {
	mt := NewMemTable()
	t.Run("Has, Put, and Get work as expected", func(t *testing.T) {
		has, err := mt.Has([]byte("hello"))
		So(t, should.BeNil(err))
		So(t, should.BeFalse(has))

		err = mt.Put([]byte("hello"), []byte("world"))
		So(t, should.BeNil(err))
		err = mt.Put([]byte("goodbye"), []byte("sky"))
		So(t, should.BeNil(err))
		err = mt.Put([]byte("apple"), []byte("juice"))
		So(t, should.BeNil(err))

		has, err = mt.Has([]byte("apple"))
		So(t, should.BeNil(err))
		So(t, should.BeTrue(has))

		v, err := mt.Get([]byte("apple"))
		So(t, should.BeNil(err))
		So(t, should.Equal(string(v), "juice"))

		v, err = mt.Get([]byte("notinthere"))
		So(t, should.NotBeNil(err))
		So(t, should.HaveSameTypeAs(err, &NotFoundError{}))
		So(t, should.BeNil(v))
	})

	t.Run("RangeScan returns values in a sorted order", func(t *testing.T) {
		mt.RangeScan([]byte("a"), []byte("z"))
	})

}

//func TestMemTableBenchmark(b *testing.B) {
//
//}
