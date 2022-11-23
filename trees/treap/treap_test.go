package treap

import (
	"math/rand"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// buildSimpleTreap returns a specific three-node treap for testing:
//
//					(4,0)   <-root  (value,priority)
//		   	       /     \
//	           (2,2)     (6,4)
//
// The values were chosen to allow testing violations of vals and priorities.
func buildSimpleTreap() *Treap {
	t := &Treap{}
	i := 0
	priority_generator = func() int {
		i++
		return t.root.priority + i*2
	}
	defer func() {
		priority_generator = rand.Int
	}()

	_ = t.Insert(4)
	_ = t.Insert(2)
	_ = t.Insert(6)

	return t
}

// TODO: Get is only partially built. These will need completion when Treap api is abstracted completely.
func TestGet(te *testing.T) {
	Convey("Get tests", te, func() {
		t := buildSimpleTreap()

		Convey("When Get() called for existing values", func() {
			vals := []int{2, 4, 6}
			for _, v := range vals {
				result := t.Get(v)
				So(result, ShouldEqual, v)
			}
		})

		Convey("When no such left child", func() {
			v := t.Get(3)
			So(v, ShouldNotEqual, 3)
		})

		Convey("When no such right child", func() {
			v := t.Get(5)
			So(v, ShouldNotEqual, 5)
		})
	})
}

func TestInsertion(te *testing.T) {
	// Insertion tests are a bit complicated since we must validate the structure of
	// the tree to ensure its invariants; but doing so requires overriding internal
	// assumptions, such as priority generation.
	Convey("Insertion tests", te, func() {
		Convey("When treap is empty", func() {
			t := Treap{}
			err := t.Insert(3)
			So(err, ShouldBeNil)
			So(t.root.val, ShouldEqual, 3)
			So(t.root.priority, ShouldEqual, 0)
		})

		Convey("When a duplicate value is added", func() {
			t := Treap{}
			err := t.Insert(3)
			So(err, ShouldBeNil)
			err = t.Insert(3)
			So(err, ShouldBeError, ErrDuplicateValue)
		})

		Convey("When a simple tree is built", func() {
			priorities := [2]int{1, 2}
			i := 0
			priority_generator = func() int {
				next := priorities[i]
				i++
				return next
			}
			defer func() {
				priority_generator = rand.Int
			}()

			t := Treap{}

			err := t.Insert(4)
			So(err, ShouldBeNil)
			So(t.root.val, ShouldEqual, 4)
			So(t.root.priority, ShouldEqual, 0)

			// Insert 2, overriding its priority as 1
			err = t.Insert(2)
			So(err, ShouldBeNil)
			So(t.root.left.val, ShouldEqual, 2)
			So(t.root.left.priority, ShouldEqual, 1)

			// Insert 5, overridiing its priority as 2
			err = t.Insert(5)
			So(err, ShouldBeNil)
			So(t.root.right.val, ShouldEqual, 5)
			So(t.root.right.priority, ShouldEqual, 2)
		})

		Convey("When a left-rotation on left child is required", func() {
			t := buildSimpleTreap()

			priority_generator = func() int {
				return 1
			}
			defer func() {
				priority_generator = rand.Int
			}()
			// Inserting 1 with priority 1 places a node as root's left-child's
			// left child, violating heap order and forcing a left rotation.
			err := t.Insert(1)
			So(err, ShouldBeNil)
			So(t.root.left.val, ShouldEqual, 1)
			So(t.root.left.priority, ShouldEqual, 1)
			So(t.root.left.left, ShouldBeNil)
			So(t.root.left.right, ShouldNotBeNil)
			So(t.root.left.right.val, ShouldEqual, 2)
			So(t.root.left.right.priority, ShouldEqual, 2)
			So(t.root.left.right.left, ShouldBeNil)
			So(t.root.left.right.right, ShouldBeNil)
		})

		Convey("When a right-rotation on left child is required", func() {
			t := buildSimpleTreap()

			priority_generator = func() int {
				return 1
			}
			defer func() {
				priority_generator = rand.Int
			}()
			// Inserting 3 with priority 1 places a node as root's left-child's
			// right child, violating heap order and forcing a right rotation.
			err := t.Insert(3)
			So(err, ShouldBeNil)
			So(t.root.left.val, ShouldEqual, 3)
			So(t.root.left.priority, ShouldEqual, 1)
			So(t.root.left.right, ShouldBeNil)
			So(t.root.left.left, ShouldNotBeNil)
			So(t.root.left.left.val, ShouldEqual, 2)
			So(t.root.left.left.priority, ShouldEqual, 2)
			So(t.root.left.left.left, ShouldBeNil)
			So(t.root.left.left.right, ShouldBeNil)
		})

		Convey("When a left-rotation on right child is required", func() {
			t := buildSimpleTreap()

			priority_generator = func() int {
				return 1
			}
			defer func() {
				priority_generator = rand.Int
			}()

			err := t.Insert(5)
			So(err, ShouldBeNil)
			So(t.root.right.val, ShouldEqual, 5)
			So(t.root.right.priority, ShouldEqual, 1)
			So(t.root.right.left, ShouldBeNil)
			So(t.root.right.right, ShouldNotBeNil)
			So(t.root.right.right.val, ShouldEqual, 6)
			So(t.root.right.right.priority, ShouldEqual, 4)
			So(t.root.right.right.left, ShouldBeNil)
			So(t.root.right.right.right, ShouldBeNil)
		})

		Convey("When a right-rotation on right child is required", func() {
			t := buildSimpleTreap()

			priority_generator = func() int {
				return 1
			}
			defer func() {
				priority_generator = rand.Int
			}()

			err := t.Insert(7)
			So(err, ShouldBeNil)
			So(t.root.right.val, ShouldEqual, 7)
			So(t.root.right.priority, ShouldEqual, 1)
			So(t.root.right.right, ShouldBeNil)
			So(t.root.right.left, ShouldNotBeNil)
			So(t.root.right.left.val, ShouldEqual, 6)
			So(t.root.right.left.priority, ShouldEqual, 4)
			So(t.root.right.left.left, ShouldBeNil)
			So(t.root.right.left.right, ShouldBeNil)
		})

		Convey("When random trees are generated, all trees are both bst-ordered and heap-ordered", func() {
			for n := 0; n < 4; n++ {
				t := Treap{}
				for i := 0; i < 100; i++ {
					_ = t.Insert(rand.Int() % 10000)
					//So(err, ShouldBeNil)
				}

				So(isHeap(t.root, te), ShouldBeTrue)
				So(isBST(t.root), ShouldBeTrue)

				// Uncomment to view small trees (n < 10) for manual verification.
				//s, _ := t.Format(BFSOrder)
				//ioutil.WriteFile("test.txt", []byte(s), 0777)
			}
		})
	})
}

// Verifies that all nodes are in min-heap order, such that every
// node's priority is less than its children.
func isHeap(node *treapNode, t *testing.T) bool {
	if node == nil {
		return true
	}

	if node.left != nil && node.priority > node.left.priority {
		t.Logf("Violation at (%d,%d) with left (%d,%d)\n", node.val, node.priority, node.left.val, node.left.priority)
		return false
	}
	if node.right != nil && node.priority > node.right.priority {
		return false
	}

	return isHeap(node.left, t) && isHeap(node.right, t)
}

// Verifies that all nodes are in bst-order, such that all of a node's subtree
// have values less than the node, and vice versa for the left subtree.
func isBST(node *treapNode) bool {
	if node == nil {
		return true
	}

	if node.left != nil && node.left.val > node.val {
		return false
	}
	if node.right != nil && node.val > node.right.val {
		return false
	}

	return isBST(node.left) && isBST(node.right)
}

func toString(order TraversalOrder) string {
	switch order {
	case PreOrder:
		return "PreOrder"
	case InOrder:
		return "InOrder"
	case PostOrder:
		return "PostOrder"
	default:
		return "UNKOWN ORDER"
	}
}

func TestFormat(te *testing.T) {
	Convey("When various ordered formats are requested", te, func() {
		Convey("When treap is empty", func() {
			t := Treap{}
			for _, order := range []TraversalOrder{PreOrder, InOrder, PostOrder, BFSOrder} {
				result, err := t.Format(order)
				So(err, ShouldBeNil)
				So(result, ShouldEqual, "")
			}
		})

		Convey("When an invalid traversal order is passed", func() {
			t := Treap{}
			_, err := t.Format(TraversalOrder(-1))
			So(err, ShouldBeError, ErrNoSuchTraversalOrder)
		})

		Convey("When a simple treap is built", func() {
			tests := []struct {
				order    TraversalOrder
				expected string
			}{
				{
					order:    PreOrder,
					expected: "(4,0) (2,2) (6,4) ",
				},
				{
					order:    InOrder,
					expected: "(2,2) (4,0) (6,4) ",
				},
				{
					order:    PostOrder,
					expected: "(2,2) (6,4) (4,0) ",
				},
				{
					order:    BFSOrder,
					expected: "                            \n       4e+00,0e+00          \n 2e+00,2e+00  6e+00,4e+00\n",
				},
			}

			t := buildSimpleTreap()
			for _, test := range tests {
				Convey("When calling Format("+toString(test.order)+")", func() {
					result, err := t.Format(test.order)
					So(err, ShouldBeNil)
					So(result, ShouldEqual, test.expected)
				})
			}
		})
	})
}
