package db

import (
	"testing"

	"github.com/smartystreets/assertions/should"
)

func TestSkipList(t *testing.T) {
	t.Run("", func(t *testing.T) {
		// Simple Put
		sl := NewSkipList()
		key := []byte("hello")
		value := []byte("goodbye")
		err := sl.Put(key, value)
		So(t, should.BeNil(err))

		// Simple Get
		res, err := sl.Get(key)
		So(t, should.BeNil(err))
		So(t, should.Resemble(res, value))

		// Get not found
		res, err = sl.Get([]byte("not here"))
		So(t, should.NotBeNil(err))
		So(t, should.BeNil(res))

		// Put at end of range
		sl.Put([]byte("mango"), []byte("smoothie"))
		res, _ = sl.Get([]byte("mango"))
		So(t, should.Resemble(res, []byte("smoothie")))

		// Put at start of range
		sl.Put([]byte("abc"), []byte("def"))
		res, _ = sl.Get([]byte("abc"))
		So(t, should.Resemble(res, []byte("def")))

		// Put in middle of range
		sl.Put([]byte("helicopter"), []byte("airplane"))
		res, _ = sl.Get([]byte("helicopter"))
		So(t, should.Resemble(res, []byte("airplane")))

		// Overwrite
		sl.Put([]byte("hello"), []byte("again"))
		res, err = sl.Get([]byte("hello"))
		So(t, should.BeNil(err))
		So(t, should.Resemble(res, []byte("again")))
		res, _ = sl.Get([]byte("abc"))
		So(t, should.Resemble(res, []byte("def")))

		// Has
		has, err := sl.Has([]byte("mango"))
		So(t, should.BeTrue(has))
		So(t, should.BeNil(err))
		has, err = sl.Has([]byte("coconut"))
		So(t, should.BeFalse(has))
		So(t, should.BeNil(err))

		// Delete something that's not there
		sl.Delete([]byte("apple"))
		res, _ = sl.Get([]byte("apple"))
		So(t, should.BeNil(res))

		// Delete something that is there
		sl.Delete([]byte("mango"))
		res, _ = sl.Get([]byte("mango"))
		So(t, should.BeNil(res))
	})

}

//func TestMemTableBenchmark(b *testing.B) {
//
//}
