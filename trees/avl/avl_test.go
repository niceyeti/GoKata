package avl

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFind(t *testing.T) {
	Convey("Find tests", t, func() {
		Convey("When Find is called for an item that exists", func() {
			t := NewTree()
			err := t.Insert(7)
			So(err, ShouldBeNil)
			node := t.Find(7)
			So(node, ShouldNotBeNil)
			So(node.data, ShouldEqual, 7)
		})

		Convey("When Find is called for an item that does not exist", func() {
			t := NewTree()
			err := t.Insert(7)
			So(err, ShouldBeNil)
			node := t.Find(5)
			So(node, ShouldBeNil)
		})

		Convey("When Find is called on an empty tree", func() {
			t := NewTree()
			node := t.Find(5)
			So(node, ShouldBeNil)
		})

		Convey("When Find is called on a complex tree", func() {
			t := NewTree()
			vals := []int{4, 2, 6, 1, 3, 5, 7}
			for _, v := range vals {
				err := t.Insert(v)
				So(err, ShouldBeNil)
			}

			for _, v := range vals {
				node := t.Find(v)
				So(node, ShouldNotBeNil)
				So(node.data, ShouldEqual, v)
			}

			// For good measure, search for something that doesn't exist
			// and would be deep in the tree.
			node := t.Find(367)
			So(node, ShouldBeNil)
			node = t.Find(-1)
			So(node, ShouldBeNil)
		})
	})
}

func TestInsert(t *testing.T) {
	Convey("Insertion tests", t, func() {
		Convey("When calling insertion on an empty tree", func() {
			t := NewTree()
			err := t.Insert(7)
			So(err, ShouldBeNil)
			So(t.root.height, ShouldEqual, 0)
		})

		Convey("When inserting duplicate items", func() {
			t := NewTree()
			err := t.Insert(6)
			So(err, ShouldBeNil)
			err = t.Insert(6)
			So(err, ShouldBeError, ErrDuplicateItem)
		})

		Convey("When inserting items and left double rotation required", func() {
			t := NewTree()

			err := t.Insert(4)
			/*
				After Insert(4):

					4 (root)
			*/
			So(err, ShouldBeNil)
			So(t.root.data, ShouldEqual, 4)
			So(t.root.height, ShouldEqual, 0)
			So(t.root.left, ShouldBeNil)
			So(t.root.right, ShouldBeNil)

			err = t.Insert(1)
			/*
				After Insert(1):
						4 (root)
					   /
					  1
			*/
			So(err, ShouldBeNil)
			So(t.root.left, ShouldNotBeNil)
			So(t.root.left.data, ShouldEqual, 1)
			So(t.root.left.height, ShouldEqual, 0)
			So(t.root.height, ShouldEqual, 1)

			err = t.Insert(2)
			/*
				After Insert(2) and double rotation:
						2 (new root)
					   / \
					  1   4
			*/
			So(err, ShouldBeNil)
			nodeTwo := t.Find(2)
			So(nodeTwo, ShouldNotBeNil)
			So(nodeTwo.height, ShouldEqual, 1)
			So(nodeTwo, ShouldEqual, t.root)
			So(t.root.right, ShouldNotBeNil)
			So(t.root.right.data, ShouldEqual, 4)
			So(t.root.right.height, ShouldEqual, 0)
			So(t.root.left.height, ShouldEqual, 0)
			So(t.root.left.data, ShouldEqual, 1)

			treeStr := t.FormatDFS(InOrder)
			So(treeStr, ShouldEqual, "1 2 4 ")
		})

		Convey("When inserting items and right double rotation required", func() {
			t := NewTree()
			err := t.Insert(4)
			/*
				After Insert(4):

					4 (root)
			*/
			So(err, ShouldBeNil)
			So(t.root.data, ShouldEqual, 4)
			So(t.root.height, ShouldEqual, 0)
			So(t.root.left, ShouldBeNil)
			So(t.root.right, ShouldBeNil)

			err = t.Insert(7)
			/*
				After Insert(7):

						4 (root)
					     \
					      7
			*/
			So(err, ShouldBeNil)
			So(t.root.right, ShouldNotBeNil)
			So(t.root.right.data, ShouldEqual, 7)
			So(t.root.right.height, ShouldEqual, 0)
			So(t.root.height, ShouldEqual, 1)

			err = t.Insert(5)
			/*
				After Insert(5) and double rotation:
						5 (new root)
					   / \
					  4   7
			*/
			So(err, ShouldBeNil)
			nodeFive := t.Find(5)
			So(nodeFive, ShouldNotBeNil)
			So(nodeFive.height, ShouldEqual, 1)
			So(nodeFive, ShouldEqual, t.root)
			So(t.root.right, ShouldNotBeNil)
			So(t.root.right.data, ShouldEqual, 7)
			So(t.root.right.height, ShouldEqual, 0)
			So(t.root.left.data, ShouldEqual, 4)
			So(t.root.left.height, ShouldEqual, 0)

			treeStr := t.FormatDFS(InOrder)
			So(treeStr, ShouldEqual, "4 5 7 ")
		})

		Convey("When inserting items and left single rotation required", func() {
			t := NewTree()

			err := t.Insert(4)
			/*
				After Insert(4):

					4 (root)
			*/
			So(err, ShouldBeNil)
			So(t.root.data, ShouldEqual, 4)
			So(t.root.height, ShouldEqual, 0)
			So(t.root.left, ShouldBeNil)
			So(t.root.right, ShouldBeNil)

			err = t.Insert(2)
			/*
				After Insert(2):
						4 (root)
					   /
					  2
			*/
			So(err, ShouldBeNil)
			So(t.root.left, ShouldNotBeNil)
			So(t.root.left.data, ShouldEqual, 2)
			So(t.root.left.height, ShouldEqual, 0)
			So(t.root.height, ShouldEqual, 1)

			err = t.Insert(1)
			/*
				After Insert(2) and single rotation:
						2 (new root)
					   / \
					  1   4
			*/
			So(err, ShouldBeNil)
			nodeTwo := t.Find(2)
			So(nodeTwo, ShouldNotBeNil)
			So(nodeTwo.height, ShouldEqual, 1)
			So(nodeTwo, ShouldEqual, t.root)
			So(t.root.right, ShouldNotBeNil)
			So(t.root.right.data, ShouldEqual, 4)
			So(t.root.right.height, ShouldEqual, 0)
			So(t.root.left.height, ShouldEqual, 0)
			So(t.root.left.data, ShouldEqual, 1)

			treeStr := t.FormatDFS(InOrder)
			So(treeStr, ShouldEqual, "1 2 4 ")
		})

		Convey("When inserting items and right single rotation required", func() {
			t := NewTree()

			err := t.Insert(1)
			/*
				After Insert(1):

					1 (root)
			*/
			So(err, ShouldBeNil)
			So(t.root.data, ShouldEqual, 1)
			So(t.root.height, ShouldEqual, 0)
			So(t.root.left, ShouldBeNil)
			So(t.root.right, ShouldBeNil)

			err = t.Insert(2)
			/*
				After Insert(1):
						1 (root)
					     \
					      2
			*/
			So(err, ShouldBeNil)
			So(t.root.right, ShouldNotBeNil)
			So(t.root.right.data, ShouldEqual, 2)
			So(t.root.right.height, ShouldEqual, 0)
			So(t.root.height, ShouldEqual, 1)

			err = t.Insert(4)
			/*
				After Insert(4) and single rotation:
						2 (new root)
					   / \
					  1   4
			*/
			So(err, ShouldBeNil)
			nodeTwo := t.Find(2)
			So(nodeTwo, ShouldNotBeNil)
			So(nodeTwo.height, ShouldEqual, 1)
			So(nodeTwo, ShouldEqual, t.root)
			So(t.root.right, ShouldNotBeNil)
			So(t.root.right.data, ShouldEqual, 4)
			So(t.root.right.height, ShouldEqual, 0)
			So(t.root.left.height, ShouldEqual, 0)
			So(t.root.left.data, ShouldEqual, 1)

			treeStr := t.FormatDFS(InOrder)
			So(treeStr, ShouldEqual, "1 2 4 ")
		})
	})
}
