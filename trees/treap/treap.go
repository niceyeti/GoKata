package treap

import (
	"container/list"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"
)

var (
	ErrNoSuchTraversalOrder error = errors.New("no such traversal order")
	ErrValueNotFound        error = errors.New("value not found")
	ErrDuplicateValue       error = errors.New("duplicate value")
)

// A Treap is a BST whose nodes left/right relationships preserve gt/lt
// bst order relationships, but whose vertical relationships preserve
// min-heap order via tree-rotations on insertion. The result is a simple
// bst randomization property ensuring that a tree's height is lg(n) on
// average, avoiding the degenerate O(n) cases for non-random BSTs.
// Note: this implementation is purely for practice; it does not support
// concurrency and should be abstracted to support arbitrary data types.
type Treap struct {
	root *treapNode
}

type treapNode struct {
	val         int
	priority    int
	left, right *treapNode
}

var priority_generator func() int = func() int {
	return rand.Int()
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

type TraversalOrder int

const (
	PreOrder TraversalOrder = iota + 1
	PostOrder
	InOrder
	BFSOrder
)

// Format returns the prefix, postfix, or inorder representation of the treap.
// BFS is also supported, which is a completely custom-spaced tree representation
// for manual testing/displaying.
func (t *Treap) Format(order TraversalOrder) (string, error) {
	var sb strings.Builder
	visitor := func(node *treapNode) {
		sb.WriteString(fmt.Sprintf("(%d,%d) ", node.val, node.priority))
	}

	switch order {
	case PreOrder:
		t.visitPreOrder(t.root, visitor)
	case InOrder:
		t.visitInOrder(t.root, visitor)
	case PostOrder:
		t.visitPostOrder(t.root, visitor)
	case BFSOrder:
		s := t.formatBFS()
		sb.WriteString(s)
	default:
		return "", ErrNoSuchTraversalOrder
	}

	return sb.String(), nil
}

// The leading bit index of an int is useful because it also describes the
// level of a node in a binary tree, when nodes are numbered in level-order
// starting from 1. It is also floor(log2(n)).
func leadingBitIndex(n uint) (i uint) {
	for n != 0 {
		n = n >> 1
		i++
	}
	return
}

// formatBFS prints the tree vertically using BFS, using a simple procedural
// spacing algorithm to equally distribute the nodes at a given level. This isn't
// the tightest format to visualize parent-child relationships, but is useful
// for manual testing.
func (t *Treap) formatBFS() string {
	// Node width is derived from this format: 5e+00,5e+00 which is from "%1.0e,%1.0e"
	nw := 11
	// Minimum width around nodes, i.e. at the deepest (most crowded) level of the tree.
	mw := 2
	// Tree depth, to determine the line width required to print the widest section of the tree.
	d := t.depth(t.root)
	// Maximum line width, the total space required to evenly space nodes at the deepest level.
	lw := int(math.Exp2(float64(d-1)))*nw + (int(math.Exp2(float64(d-1)))+1)*mw
	// Space around nodes for a given level; this changes for each level.
	sl := 0
	// Max number of nodes on a particular line.
	nl := 0

	var sb, line strings.Builder
	var curLevel uint

	visitor := func(node *treapNode, nodeNumber uint) {
		// Stateful values: the formatting state is fully defined by the height/level in the tree.
		// When a new level is encounted, all the spacing parameters are updated.
		level := leadingBitIndex(nodeNumber)
		if level != curLevel {
			// We reached a new level, so time to recalculate spacing vals
			// Note that a level is defined by the leading-most bit in @level.
			curLevel = level
			// Max number of nodes for this level.
			nl = int(math.Exp2(float64(level - 1)))
			// Padding space surrounding nodes at this level.
			sl = (lw - nl*nw) / (nl + 1)

			// Write out the previous line
			for line.Len() < lw {
				line.WriteString(" ")
			}
			sb.WriteString(line.String() + "\n")
			line.Reset()
		}

		// Line-predecessors is the (maximum) number of preceding nodes on a line
		lp := int(nodeNumber) - nl
		// Absolute starting position is defined by the number of preceding nodes and padding space.
		as := lp*nw + (lp+1)*sl
		for line.Len() < (as - 1) {
			line.WriteString(" ")
		}
		ns := fmt.Sprintf("%1.0e,%1.0e", float64(node.val), float64(node.priority))
		line.WriteString(ns)
	}
	t.visitBFS(visitor)

	// Write any remaining content from the last line.
	if line.Len() > 0 {
		sb.WriteString(line.String() + "\n")
	}

	return sb.String()
}

func (t *Treap) visitBFS(fn func(*treapNode, uint)) {
	if t.root == nil {
		return
	}

	type bfsData struct {
		// Number is assigned to nodes in level-order, starting from 1.
		// This number is equivalent to heap-array indices, such that spatial
		// relations can be known, since a node's left child is 2*number and right child
		// is 2*number+1, its height is floor(lg(number)), etc.
		number uint
		node   *treapNode
	}

	q := list.New()
	q.PushBack(bfsData{
		number: 1,
		node:   t.root})

	for q.Len() > 0 {
		f := q.Front()
		item := f.Value.(bfsData)
		q.Remove(f)

		fn(item.node, item.number)

		if item.node.left != nil {
			q.PushBack(bfsData{
				number: item.number * 2,
				node:   item.node.left,
			})
		}

		if item.node.right != nil {
			q.PushBack(
				bfsData{
					number: item.number*2 + 1,
					node:   item.node.right,
				})
		}
	}
}

func (t *Treap) depth(node *treapNode) int {
	if node == nil {
		return 0
	}

	return max(
		t.depth(node.left)+1,
		t.depth(node.right)+1)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (t *Treap) visitPreOrder(node *treapNode, fn func(*treapNode)) {
	if node == nil {
		return
	}

	fn(node)
	t.visitPreOrder(node.left, fn)
	t.visitPreOrder(node.right, fn)
}

func (t *Treap) visitPostOrder(node *treapNode, fn func(*treapNode)) {
	if node == nil {
		return
	}

	t.visitPostOrder(node.left, fn)
	t.visitPostOrder(node.right, fn)
	fn(node)
}

func (t *Treap) visitInOrder(node *treapNode, fn func(*treapNode)) {
	if node == nil {
		return
	}

	t.visitInOrder(node.left, fn)
	fn(node)
	t.visitInOrder(node.right, fn)
}

func (t *Treap) Insert(val int) error {
	if t.root == nil {
		t.root = &treapNode{
			val:      val,
			priority: 0,
			right:    nil,
			left:     nil,
		}
		return nil
	}

	return t.insert(val, &t.root)
}

// TODO: if this alg works, simplify by passing only parentLink, since it also contains @node as its value.
// TODO: what if priorities are not unique?
func (t *Treap) insert(val int, parentLink **treapNode) error {
	node := *parentLink
	if node.val == val {
		return ErrDuplicateValue
	}

	// TODO: rotation when priorities are equal

	// val < node.val, so traverse left
	if val < node.val {
		if node.left == nil {
			node.left = &treapNode{
				val:      val,
				priority: priority_generator(),
			}
		} else if err := t.insert(val, &node.left); err != nil {
			return err
		}
		*parentLink = t.rotateLeftChild(node)
	} else {
		// Case: val > node.val, so traverse right
		if node.right == nil {
			node.right = &treapNode{
				val:      val,
				priority: priority_generator(),
			}
		} else if err := t.insert(val, &node.right); err != nil {
			return err
		}
		*parentLink = t.rotateRightChild(node)
	}

	return nil
}

func (t *Treap) rotateLeftChild(node *treapNode) *treapNode {
	if node.priority < node.left.priority {
		// priorities already obey heap-order, so just return
		return node
	}

	leftChild := node.left
	node.left = leftChild.right
	leftChild.right = node

	return leftChild
}

func (t *Treap) rotateRightChild(node *treapNode) *treapNode {
	if node.priority < node.right.priority {
		// priorities already obey heap-order, so just return
		return node
	}

	rightChild := node.right
	node.right = rightChild.left
	rightChild.left = node

	return rightChild
}

// Get retrieves an item in the tree if it exists, else returns -math.MaxInt.
// Passing the value to retrieve and returning it when found is redundant,
// this is just for a demo. A properly abstracted treap would search by id.
// TODO: abstract the treap to support arbitrary data types 1) using an
// Equals() or Id() interface, or 2) using templating.
func (t *Treap) Get(val int) int {
	if node := t.get(val, t.root); node != nil {
		return node.val
	}

	return -math.MaxInt
}

func (t *Treap) get(val int, node *treapNode) *treapNode {
	if node == nil {
		return nil
	}

	if node.val == val {
		return node
	}

	if val < node.val {
		return t.get(val, node.left)
	}

	return t.get(val, node.right)
}
