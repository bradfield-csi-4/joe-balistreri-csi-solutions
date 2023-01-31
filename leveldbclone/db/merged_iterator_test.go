package db

import (
	"testing"

	"github.com/smartystreets/assertions/should"
)

func TestMergedIterator(t *testing.T) {
	t.Run("Properly handles being empty", func(t *testing.T) {
		m := &MergedIterator{}
		So(t, should.BeFalse(m.Next()))
		So(t, should.BeNil(m.Value()))
		So(t, should.BeNil(m.Key()))
	})
	t.Run("Reads responses in the correct order", func(t *testing.T) {
		l1 := NewLinkedList()
		l1.Put([]byte("a"), []byte("a"))
		l1.Put([]byte("c"), []byte("c"))
		l1.Put([]byte("d"), []byte("d"))
		l1.Put([]byte("e"), nil)
		it1, _ := l1.RangeScan(nil, nil)
		So(t, should.BeTrue(it1.Next())) // initialize

		l2 := NewLinkedList()
		l2.Put([]byte("a"), []byte("A"))
		l2.Put([]byte("b"), []byte("B"))
		l2.Put([]byte("c"), []byte("C"))
		l2.Put([]byte("e"), []byte("E"))
		it2, _ := l2.RangeScan(nil, nil)
		So(t, should.BeTrue(it2.Next())) // initialize

		l3 := NewLinkedList()
		l3.Put([]byte("a"), []byte("1"))
		l3.Put([]byte("c"), []byte("2"))
		l3.Put([]byte("e"), []byte("3"))
		l3.Put([]byte("f"), []byte("4"))
		it3, _ := l3.RangeScan(nil, nil)
		So(t, should.BeTrue(it3.Next())) // initialize

		m := &MergedIterator{its: []Iterator{it1, it2, it3}}

		So(t, should.BeTrue(m.Next()))
		So(t, should.Equal(string(m.Key()), "a"))
		So(t, should.Equal(string(m.Value()), "a"))

		So(t, should.BeTrue(m.Next()))
		So(t, should.Equal(string(m.Key()), "b"))
		So(t, should.Equal(string(m.Value()), "B"))

		So(t, should.BeTrue(m.Next()))
		So(t, should.Equal(string(m.Key()), "c"))
		So(t, should.Equal(string(m.Value()), "c"))

		So(t, should.BeTrue(m.Next()))
		So(t, should.Equal(string(m.Key()), "d"))
		So(t, should.Equal(string(m.Value()), "d"))

		// So(t, should.BeTrue(m.Next()))
		// So(t, should.Equal(string(m.Key()), "e"))
		// So(t, should.Equal(string(m.Value()), "E"))

		So(t, should.BeTrue(m.Next()))
		So(t, should.Equal(string(m.Key()), "f"))
		So(t, should.Equal(string(m.Value()), "4"))

		So(t, should.BeFalse(m.Next()))
	})
}
