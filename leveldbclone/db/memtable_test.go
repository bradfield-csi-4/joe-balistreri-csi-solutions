package db

import (
	"fmt"
	"math/rand"
	"os/exec"
	"testing"

	"github.com/smartystreets/assertions/should"
)

func So(t *testing.T, s string) {
	t.Helper()
	if s != "" {
		t.Fatal(s)
	}
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

var NewMemTable = func() (DB, func()) {
	return NewKVStore(randSeq(6) + "test")
}

// var NewMemTable = NewSkipList

func TestMain(m *testing.M) {
	m.Run()
	fmt.Println(exec.Command("bash", "-c", "rm *test.wal").Run())
	fmt.Println(exec.Command("bash", "-c", "rm *test.sst").Run())
}

var memtableTypes = []func() (DB, func()){
	func() (DB, func()) {
		return NewKVStore(randSeq(6) + "test")
	},
	func() (DB, func()) {
		return NewSkipList(), func() {}
	},
	func() (DB, func()) {
		return NewLinkedList(), func() {}
	},
}

func TestMemTable(t *testing.T) {
	for _, newMemTable := range memtableTypes {
		t.Run("Has, Put, and Get work as expected", func(t *testing.T) {
			mt, done := newMemTable()
			defer done()
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
			So(t, should.Equal(err, ErrNotFound))
			So(t, should.BeNil(v))
		})

		t.Run("RangeScan handles an empty db", func(t *testing.T) {
			mt, done := newMemTable()
			defer done()
			i, _ := mt.RangeScan([]byte{}, []byte{})
			So(t, should.BeNil(i.Key()))
			So(t, should.BeNil(i.Value()))
			So(t, should.BeFalse(i.Next()))
		})

		t.Run("RangeScan returns values in a sorted order", func(t *testing.T) {
			mt, done := newMemTable()
			defer done()
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
			mt, done := newMemTable()
			defer done()
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

		t.Run("Delete works properly", func(t *testing.T) {
			mt, done := newMemTable()
			defer done()
			mt.Put([]byte("hello"), []byte("world"))
			mt.Put([]byte("goodbye"), []byte("sky"))
			mt.Put([]byte("apple"), []byte("juice"))

			has, _ := mt.Has([]byte("goodbye"))
			So(t, should.BeTrue(has))

			err := mt.Delete([]byte("goodbye"))
			So(t, should.BeNil(err))

			// can delete twice
			err = mt.Delete([]byte("goodbye"))
			So(t, should.BeNil(err))

			has, _ = mt.Has([]byte("goodbye"))
			So(t, should.BeFalse(has))

			v, err := mt.Get([]byte("goodbye"))
			So(t, should.BeNil(v))
			So(t, should.NotBeNil(err))

			i, err := mt.RangeScan([]byte("a"), []byte("z"))
			So(t, should.BeNil(err))

			So(t, should.Equal(string(i.Key()), "apple"))
			So(t, should.Equal(string(i.Value()), "juice"))
			So(t, should.BeTrue(i.Next()))

			So(t, should.Equal(string(i.Key()), "hello"))
			So(t, should.Equal(string(i.Value()), "world"))
			So(t, should.BeFalse(i.Next()))

			So(t, should.BeNil(i.Error()))
		})
	}
}

//func TestMemTableBenchmark(b *testing.B) {
//
//}
