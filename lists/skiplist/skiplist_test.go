package skiplist

import (
	"math/rand"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewSkiplist(t *testing.T) {
	Convey("When NewSkiplist is called", t, func() {
		sl := NewSkiplist(3)
		So(sl.r, ShouldEqual, 3)
		So(len(sl.root.next), ShouldEqual, 3)
	})
}

func TestInsertion(t *testing.T) {
	Convey("When Insert is called", t, func() {
		Convey("When insert is called on an empty list", func() {
			sl := NewSkiplist(3)
			err := sl.Insert(123)
			So(err, ShouldBeNil)
			So(sl.root.next[0].value, ShouldEqual, 123)
		})

		Convey("When a duplicate is inserted", func() {
			sl := NewSkiplist(3)
			err := sl.Insert(123)
			So(err, ShouldBeNil)
			err = sl.Insert(123)
			So(err, ShouldBeError, ErrDuplicateValue)
		})

		Convey("When Insert is called repeatedly", func() {
			sl := NewSkiplist(8)
			for i := 0; i < 100; i++ {
				err := sl.Insert(rand.Int())
				So(err, ShouldBeNil)
			}
		})
	})
}

func TestGet(t *testing.T) {
	Convey("When Get is called", t, func() {
		Convey("When Get is called on an empty list", func() {
			sl := NewSkiplist(8)
			_, err := sl.Get(123)
			So(err, ShouldBeError, ErrValueNotFound)
		})

		Convey("When Get is called on a singleton list", func() {
			sl := NewSkiplist(8)
			err := sl.Insert(123)
			So(err, ShouldBeNil)
			n, err := sl.Get(123)
			So(err, ShouldBeNil)
			So(n, ShouldEqual, 123)
		})
	})
}

func TestDeletion(t *testing.T) {
	Convey("When Delete is called", t, func() {
		Convey("When Delete is called on an empty list", func() {
			sl := NewSkiplist(4)
			err := sl.Delete(123)
			So(err, ShouldBeError, ErrValueNotFound)
		})

		Convey("When Delete is called for an item that does not exist", func() {
			sl := NewSkiplist(8)
			err := sl.Insert(123)
			So(err, ShouldBeNil)
			err = sl.Delete(456)
			So(err, ShouldBeError, ErrValueNotFound)
		})

		Convey("When Delete drains a list", func() {
			sl := NewSkiplist(4)
			vals := []int{1, 2, 3}
			for _, val := range vals {
				err := sl.Insert(val)
				So(err, ShouldBeNil)

			}

			for _, val := range vals {
				err := sl.Delete(val)
				So(err, ShouldBeNil)
			}

			Convey("Re-adding the same items to the now empty list succeeds", func() {
				for _, val := range vals {
					err := sl.Insert(val)
					So(err, ShouldBeNil)
				}

				for _, val := range vals {
					err := sl.Delete(val)
					So(err, ShouldBeNil)
				}

				err := sl.Delete(vals[0])
				So(err, ShouldBeError, ErrValueNotFound)
			})
		})
	})
}
