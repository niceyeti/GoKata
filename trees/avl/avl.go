package avl

import (
	"container/list"
	"errors"
	"fmt"
	"math"
	"strings"
)

type Node struct {
	left, right *Node
	data        int
	// Height is defined as the longest path from this node to a leaf (thus zero, if it is a leaf).
	height int
}

type AvlTree struct {
	root      *Node
	nodeCount int
}

var (
	ErrDuplicateItem error = errors.New("duplicate item")
)

const allowedImbalance = 1

func NewTree() *AvlTree {
	return &AvlTree{}
}

func (t *AvlTree) Insert(n int) {
	if t.root == nil {
		t.root = &Node{
			data:   n,
			height: 1,
		}
		return
	}

	t.insert(&t.root, n)
}

func (t *AvlTree) insert(node **Node, n int) error {
	// base case
	if *node == nil {
		*node = &Node{
			data:   n,
			height: 0,
		}
		return nil
	}

	if n == (*node).data {
		return ErrDuplicateItem
	}

	if n < (*node).data {
		t.insert(&(*node).left, n)
	} else {
		t.insert(&(*node).right, n)
	}

	t.balance(node)

	return nil
}

func (t *AvlTree) balance(node **Node) {
	if (*node).left == nil || (*node).right == nil {
		return
	}

	// TODO: still some lurking nils here, need to simplify
	if leftImbalance(*node) {
		if outerLeftDeeper(*node) {
			// run left outer single rotation
			rotateWithLeftChild(node)
		} else {
			// run left inner double rotation
			doubleRotateWithLeftChild(node)
		}
	} else if rightImbalance(*node) {
		if outerRightDeeper(*node) {
			// run left outer single rotation
			rotateWithRightChild(node)
		} else {
			// run left inner double rotation
			doubleRotateWithRightChild(node)
		}
	}
}

func leftImbalance(node *Node) bool {
	return node.left.height-node.right.height > allowedImbalance
}

func rightImbalance(node *Node) bool {
	return node.right.height-node.left.height > allowedImbalance
}

func outerLeftDeeper(node *Node) bool {
	return node.left.left != nil &&
		node.left.right != nil &&
		node.left.left.height > node.left.right.height
}

func outerRightDeeper(node *Node) bool {
	return node.right.left != nil &&
		node.right.right != nil &&
		node.right.right.height > node.right.left.height
}

// The rotation funcs are best understood via diagram only.
func rotateWithLeftChild(root **Node) {
	k2 := *root
	k1 := k2.left
	k2.left = k1.right
	k1.right = k2
	*root = k1

	setHeight(k1)
	setHeight(k2)
}

// The rotation funcs are best understood via diagram only.
func rotateWithRightChild(root **Node) {
	k2 := *root
	k1 := k2.right
	k2.right = k1.left
	k1.left = k2
	*root = k1

	setHeight(k1)
	setHeight(k2)
}

func setHeight(node *Node) {
	node.height = 1 + max(height(node.left), height(node.right))
}

// The double rotation operations can be performed via two single
// rotations, though a pencil example is necessary to demonstrate.
func doubleRotateWithLeftChild(node **Node) {
	rotateWithRightChild(&(*node).left)
	rotateWithLeftChild(node)
}

// The double rotation operations can be performed via two single
// rotations, though a pencil example is necessary to demonstrate.
func doubleRotateWithRightChild(node **Node) {
	rotateWithLeftChild(&(*node).right)
	rotateWithRightChild(node)
}

func height(node *Node) int {
	if node == nil {
		return -1
	}
	return node.height
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func (t *AvlTree) Delete(n int) {
	t.delete(&t.root, n)
}

func (t *AvlTree) delete(node **Node, n int) {
	// TODO: handle deletion of root

	if *node == nil {
		// item not found
		return
	}

	if n < (*node).data {
		t.delete(&(*node).left, n)
	} else if n > (*node).data {
		t.delete(&(*node).right, n)
	} else if (*node).left != nil && (*node).right != nil {
		// Simply
		(*node).data = findMin((*node).right).data
		t.delete(&(*node).right, (*node).data)
	} else {

	}

	t.balance(node)
}

// formatBFS prints the tree vertically using BFS, using a simple procedural
// spacing algorithm to equally distribute the nodes at a given level. This isn't
// the tightest format to visualize parent-child relationships, but is useful
// for manual testing.
func (t *AvlTree) FormatBFS() string {
	if t.root == nil {
		return "<empty>"
	}

	// Node width is derived from this format: 5e+00,5e+00 which is from "%1.0e,%1.0e"
	nw := 11
	// Minimum width around nodes, i.e. at the deepest (most crowded) level of the tree.
	mw := 2
	// Tree depth, to determine the line width required to print the widest section of the tree.
	d := t.root.height
	// Maximum line width, the total space required to evenly space nodes at the deepest level.
	lw := int(math.Exp2(float64(d-1)))*nw + (int(math.Exp2(float64(d-1)))+1)*mw
	// Space around nodes for a given level; this changes for each level.
	sl := 0
	// Max number of nodes on a particular line.
	nl := 0

	var sb, line strings.Builder
	var curLevel uint

	visitor := func(node *Node, nodeNumber uint) {
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
		ns := fmt.Sprintf("%1.0e", float64(node.data))
		line.WriteString(ns)
	}
	t.visitBFS(visitor)

	// Write any remaining content from the last line.
	if line.Len() > 0 {
		sb.WriteString(line.String() + "\n")
	}

	return sb.String()
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

func (t *AvlTree) visitBFS(fn func(*Node, uint)) {
	if t.root == nil {
		return
	}

	type bfsData struct {
		// Number is assigned to nodes in level-order, starting from 1.
		// This number is equivalent to heap-array indices, such that spatial
		// relations can be known, since a node's left child is 2*number and right child
		// is 2*number+1, its height is floor(lg(number)), etc.
		number uint
		node   *Node
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

func findMin(node *Node) *Node {
	if node.left == nil {
		return node
	}
	return findMin(node.left)
}

// Find returns a node given its value; obviously this is
// redundant, it is purely for demonstration.
func (t *AvlTree) Find(n int) *Node {
	return t.find(t.root, n)
}

func (t *AvlTree) find(node *Node, n int) *Node {
	if node == nil || node.data == n {
		return node
	}
	if n < node.data {
		return t.find(node.left, n)
	}
	return t.find(node.right, n)
}
