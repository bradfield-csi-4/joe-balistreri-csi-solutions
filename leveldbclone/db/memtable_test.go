package db

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

func TestMemTable(t *testing.T) {
	t.Run("Has, Put, and Get work as expected", func(t *testing.T) {
		mt := NewMemTable()
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
		mt := NewMemTable()
		mt.Put([]byte("hello"), []byte("world"))
		mt.Put([]byte("goodbye"), []byte("sky"))
		mt.Put([]byte("apple"), []byte("juice"))

		i, err := mt.RangeScan([]byte("a"), []byte("z"))
		So(t, should.BeNil(err))

		So(t, should.Equal(string(i.Key()), "apple"))
		So(t, should.Equal(string(i.Value()), "juice"))
		So(t, should.BeTrue(i.Next()))

		So(t, should.Equal(string(i.Key()), "goodbye"))
		So(t, should.Equal(string(i.Value()), "sky"))
		So(t, should.BeTrue(i.Next()))

		So(t, should.Equal(string(i.Key()), "hello"))
		So(t, should.Equal(string(i.Value()), "world"))
		So(t, should.BeFalse(i.Next()))

		So(t, should.BeNil(i.Error()))
	})

	t.Run("RangeScan slices a subset of data", func(t *testing.T) {
		mt := NewMemTable()
		mt.Put([]byte("hello"), []byte("world"))
		mt.Put([]byte("goodbye"), []byte("sky"))
		mt.Put([]byte("apple"), []byte("juice"))

		i, err := mt.RangeScan([]byte("d"), []byte("h"))
		So(t, should.BeNil(err))

		So(t, should.Equal(string(i.Key()), "goodbye"))
		So(t, should.Equal(string(i.Value()), "sky"))
		So(t, should.BeFalse(i.Next()))

		So(t, should.BeNil(i.Error()))
	})
}

//func TestMemTableBenchmark(b *testing.B) {
//
//}
