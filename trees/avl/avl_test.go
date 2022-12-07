package avl

import (
	"math"
	"math/rand"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

// TODO: memory leak benchmarking, perf benchmarking

func TestFormatting(t *testing.T) {
	Convey("Test recursive formatters", t, func() {
		t := NewTree()
		vals := []int{1, 2, 3, 4, 5, 6, 7, 8}
		for _, v := range vals {
			err := t.Insert(v)
			So(err, ShouldBeNil)
		}
		/*
			The resulting tree:

							4
						  /    \
						2        6
					  /   \    /   \
					 1     3  5     7
					                 \
									  8
		*/

		Convey("Test preorder traversal", func() {
			s := t.FormatDFS(PreOrder)
			So(s, ShouldEqual, "4 2 1 3 6 5 7 8 ")
		})

		Convey("Test inorder traversal", func() {
			s := t.FormatDFS(InOrder)
			So(s, ShouldEqual, "1 2 3 4 5 6 7 8 ")
		})

		Convey("Test postorder traversal", func() {
			s := t.FormatDFS(PostOrder)
			So(s, ShouldEqual, "1 3 2 5 8 7 6 4 ")
		})

		Convey("Unrecognized traversal should panic", func() {
			forcePanic := func() {
				t.FormatDFS(DFSOrder(-1))
			}
			So(forcePanic, ShouldPanic)
		})

		Convey("Test bfs traversal", func() {
			s := t.FormatBFS()
			So(s, ShouldEqual, `..........................................
....................4.....................
.............2..............6.............
.......1........3........5........7.......
......................................8
`)

			emptyTree := NewTree()
			s = emptyTree.FormatBFS()
			So(s, ShouldEqual, "<empty>")

			visits := 0
			emptyTree.visitBFS(func(node *Node, nodeNum uint) {
				visits++
			})
			So(visits, ShouldEqual, 0)
		})
	})
}

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

func TestDelete(t *testing.T) {
	Convey("Deletion tests", t, func() {
		Convey("When Delete is called on a non-existent item", func() {
			Convey("When the the item is highly nested and attempting to delete non-existent items", func() {
				t := NewTree()
				for i := 0; i < 32; i++ {
					err := t.Insert(i)
					So(err, ShouldBeNil)
				}

				err := t.Delete(-1)
				So(err, ShouldBeError, ErrItemNotFound)
				err = t.Delete(32)
				So(err, ShouldBeError, ErrItemNotFound)
			})

			Convey("When Delete is called on an empty tree", func() {
				t := NewTree()
				err := t.Delete(123)
				So(err, ShouldBeError, ErrItemNotFound)
			})
		})

		Convey("When a manually defined tree has items deleted", func() {
			t := NewTree()
			for i := 1; i <= 8; i++ {
				err := t.Insert(i)
				So(err, ShouldBeNil)
			}
			So(t.nodeCount, ShouldEqual, 8)
			/*
				The resulting tree:

								4
							  /    \
							2        6
						  /   \    /   \
						 1     3  5     7
						                 \
										  8
			*/

			// Delete node with a right child only
			err := t.Delete(7)
			So(err, ShouldBeNil)
			So(t.nodeCount, ShouldEqual, 7)
			So(t.FormatDFS(PreOrder), ShouldEqual, "4 2 1 3 6 5 8 ")
			So(t.FormatDFS(PostOrder), ShouldEqual, "1 3 2 5 8 6 4 ")

			// Delete node with no children
			err = t.Delete(8)
			So(err, ShouldBeNil)
			So(t.nodeCount, ShouldEqual, 6)
			So(t.FormatDFS(PreOrder), ShouldEqual, "4 2 1 3 6 5 ")
			So(t.FormatDFS(PostOrder), ShouldEqual, "1 3 2 5 6 4 ")

			// Delete node with left child only
			err = t.Delete(6)
			So(err, ShouldBeNil)
			So(t.nodeCount, ShouldEqual, 5)
			So(t.FormatDFS(PreOrder), ShouldEqual, "4 2 1 3 5 ")
			So(t.FormatDFS(PostOrder), ShouldEqual, "1 3 2 5 4 ")

			// Delete node with two children
			err = t.Delete(2)
			So(err, ShouldBeNil)
			So(t.nodeCount, ShouldEqual, 4)
			So(t.FormatDFS(PreOrder), ShouldEqual, "4 3 1 5 ")
			So(t.FormatDFS(PostOrder), ShouldEqual, "1 3 5 4 ")

			// Delete root node
			err = t.Delete(4)
			So(err, ShouldBeNil)
			So(t.nodeCount, ShouldEqual, 3)
			So(t.FormatDFS(PreOrder), ShouldEqual, "3 1 5 ")
			So(t.FormatDFS(PostOrder), ShouldEqual, "1 5 3 ")
		})

		Convey("When Delete is called on existing items", func() {
			t := NewTree()
			for i := 0; i < 32; i++ {
				err := t.Insert(i)
				So(err, ShouldBeNil)
			}

			err := t.Delete(0)
			So(err, ShouldBeNil)

			err = t.Delete(31)
			So(err, ShouldBeNil)

			err = t.Delete(15)
			So(err, ShouldBeNil)
		})

		Convey("When Delete empties a tree", func() {
			t := NewTree()

			// Builds and deletes tree twice, to test edge cases such as initial
			// state (nodeCount zero, root nil) and valid return to initial state.
			for i := 0; i < 2; i++ {
				for i := 0; i < 32; i++ {
					err := t.Insert(i)
					So(err, ShouldBeNil)
				}
				So(t.nodeCount, ShouldEqual, 32)

				for i := 0; i < 32; i++ {
					err := t.Delete(i)
					So(err, ShouldBeNil)
				}

				So(t.nodeCount, ShouldEqual, 0)
				So(t.root, ShouldBeNil)
			}
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
			err := t.Insert(1)
			So(err, ShouldBeNil)
			err = t.Insert(1)
			So(err, ShouldBeError, ErrDuplicateItem)

			// Add some items for a deeper tree
			for i := 2; i <= 8; i++ {
				err := t.Insert(i)
				So(err, ShouldBeNil)
			}

			err = t.Insert(8)
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

		Convey("When progressively building a large contrived tree", func() {
			tcs := []struct {
				val                                 int
				expectedPreOrder, expectedPostOrder string
			}{
				{
					1,
					"1 ",
					"1 ",
				},
				{
					2,
					"1 2 ",
					"2 1 ",
				},
				{
					3,
					"2 1 3 ",
					"1 3 2 ",
				},
				{
					4,
					"2 1 3 4 ",
					"1 4 3 2 ",
				},
				{
					5,
					"2 1 4 3 5 ",
					"1 3 5 4 2 ",
				},
				{
					6,
					"4 2 1 3 5 6 ",
					"1 3 2 6 5 4 ",
				},
				{
					7,
					"4 2 1 3 6 5 7 ",
					"1 3 2 5 7 6 4 ",
				},
				{
					8,
					"4 2 1 3 6 5 7 8 ",
					"1 3 2 5 8 7 6 4 ",
				},
			}

			t := NewTree()
			for _, tc := range tcs {
				err := t.Insert(tc.val)
				So(err, ShouldBeNil)

				// A tree's structure is fully specified by its pre-order and post-order
				// traversal, thus verifying both of these verifies their structure.
				pre := t.FormatDFS(PreOrder)
				So(pre, ShouldEqual, tc.expectedPreOrder)

				post := t.FormatDFS(PostOrder)
				So(post, ShouldEqual, tc.expectedPostOrder)
			}
		})

		Convey("When building large trees sequentially (stress test)", func() {
			t := NewTree()
			for i := 0; i < 256; i++ {
				err := t.Insert(i)
				So(err, ShouldBeNil)
			}
			So(t.root.height, ShouldEqual, 8)

			t = NewTree()
			for i := 255; i >= 0; i-- {
				err := t.Insert(i)
				So(err, ShouldBeNil)
			}
			So(t.root.height, ShouldEqual, 8)

			rand.Seed(time.Now().UnixNano())
			t = NewTree()
			n := 256
			for i := 0; i < n; i++ {
				err := t.Insert(rand.Int())
				for err != nil {
					err = t.Insert(rand.Int())
				}
			}

			// It can be shown that a tree's height is at least lg(n), and at most
			// 1.44 * lg(N+2) - 1.328.
			maxHeight := math.Ceil(1.44*math.Log2(float64(n)+2.0) - 1.328)
			minHeight := math.Floor(math.Log2(float64(n)))
			So(t.root.height >= int(minHeight) || t.root.height <= int(maxHeight), ShouldBeTrue)
		})
	})
}
