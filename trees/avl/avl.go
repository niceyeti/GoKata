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
	// Height is defined as the longest path from this node to a leaf (thus zero if it is a leaf).
	height int
}

// AvlTrees implement a balance property ensuring that sibling subtrees
// do not differ in height by more than one (a modifiable parameter),
// such that operations are O(lg(n)) on average. The balance property
// is implemented using basic rotation operations. Treaps and skiplists
// are competing data structures; to determine which one is better one
// would pencil out their mem/alg complexity. AVL trees are nice because
// they are deterministic and do not require any external dependencies.
// Treaps and skiplists both require external rand sys dependencies.
// NOTE: this is an exercise, this tree has not been fully evaluated for
// correctness, performance, nor concurrent usage.
type AvlTree struct {
	root      *Node
	nodeCount int
}

var (
	ErrDuplicateItem error = errors.New("duplicate item")
	ErrItemNotFound  error = errors.New("item not found")
)

// DFSOrder merely describes a tree traversal order, per classic tree ops.
type DFSOrder int

const (
	PreOrder DFSOrder = iota + 1
	PostOrder
	InOrder
)

// The allowed difference between right/left subtrees.
const allowedImbalance = 1

// NewTree returns an empty AVL tree.
func NewTree() *AvlTree {
	return &AvlTree{}
}

// Insert a new item in the tree.
func (t *AvlTree) Insert(n int) error {
	return t.insert(&t.root, n)
}

func (t *AvlTree) insert(node **Node, n int) (err error) {
	// base case
	if *node == nil {
		*node = &Node{
			data:   n,
			height: 0,
		}
		t.nodeCount++
		return
	}

	if n == (*node).data {
		err = ErrDuplicateItem
		return
	}

	if n < (*node).data {
		err = t.insert(&(*node).left, n)
	} else {
		err = t.insert(&(*node).right, n)
	}

	if err != nil {
		return
	}

	setHeight(*node)

	t.balance(node)

	return nil
}

func (t *AvlTree) balance(node **Node) {
	leftHeight := height((*node).left)
	rightHeight := height((*node).right)

	if leftHeight-rightHeight > allowedImbalance {
		// TODO: still some lurking nils here, need to simplify
		if outerLeftDeeper(*node) {
			// outer single rotation
			rotateWithLeftChild(node)
		} else {
			// inner double rotation
			doubleRotateWithLeftChild(node)
		}
	} else if rightHeight-leftHeight > allowedImbalance {
		if outerRightDeeper(*node) {
			// outer single rotation
			rotateWithRightChild(node)
		} else {
			// inner double rotation
			doubleRotateWithRightChild(node)
		}
	}
}

func outerLeftDeeper(node *Node) bool {
	return height(node.left.left) > height(node.left.right)
}

func outerRightDeeper(node *Node) bool {
	return height(node.right.right) > height(node.right.left)
}

// The rotation funcs are best understood via diagram.
func rotateWithLeftChild(root **Node) {
	k2 := *root
	k1 := k2.left
	k2.left = k1.right
	k1.right = k2
	*root = k1

	// Note: this order of height updates is required.
	setHeight(k2)
	setHeight(k1)
}

// The rotation funcs are best understood via diagram.
func rotateWithRightChild(root **Node) {
	k2 := *root
	k1 := k2.right
	k2.right = k1.left
	k1.left = k2
	*root = k1

	// Note: this order of height updates is required.
	setHeight(k2)
	setHeight(k1)
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

// Delete removes an item from the tree, if it exists.
func (t *AvlTree) Delete(n int) error {
	err := t.delete(&t.root, n)
	if err == nil {
		t.nodeCount--
	}
	return err
}

func (t *AvlTree) delete(node **Node, n int) (err error) {
	defer func() {
		if err == nil && *node != nil {
			t.balance(node)
		}
	}()

	if *node == nil {
		// item not found
		err = ErrItemNotFound
		return
	}

	if n < (*node).data {
		err = t.delete(&(*node).left, n)
		return
	}
	if n > (*node).data {
		err = t.delete(&(*node).right, n)
		return
	}

	// Target found and has both children.
	if (*node).left != nil && (*node).right != nil {
		// Deletion strategy: target's value is replaced by its min-right successor,
		// to preserve BST order, and then that min-right successor node is itself deleted.
		// TODO: this introduces a bias whereby a succession of deletions
		// selects the right-inner child as replacement, thus making the right tree
		// shallower over time. I have not considered the full effects.
		(*node).data = findMin((*node).right).data
		// err intentionally discarded because we know the item exists from the previous line
		_ = t.delete(&(*node).right, (*node).data)
		return
	}

	// Target found and has only a left child.
	if (*node).left != nil {
		// The node is merely in line to its children and removable.
		left := (*node).left
		// Nil out the node pointers to allow its garbage collection
		(*node).left = nil
		(*node).right = nil
		*node = left
		return
	}

	// Target found but only has right child OR no children (a leaf).
	// For these cases, the node is merely in line to its children
	// and can be removed directly.
	right := (*node).right
	// Nil out the node pointers to allow node's garbage collection
	(*node).left = nil
	(*node).right = nil
	*node = right

	return
}

type nodeVisitor func(*Node)

func (t *AvlTree) FormatDFS(order DFSOrder) string {
	sb := strings.Builder{}
	visitor := nodeVisitor(func(node *Node) {
		sb.WriteString(fmt.Sprintf("%d ", node.data))
	})

	switch order {
	case PreOrder:
		preorder(t.root, visitor)
	case InOrder:
		inorder(t.root, visitor)
	case PostOrder:
		postorder(t.root, visitor)
	default:
		panic("DFSOrder not found")
	}

	return sb.String()
}

func preorder(node *Node, visitor nodeVisitor) {
	if node == nil {
		return
	}
	visitor(node)
	preorder(node.left, visitor)
	preorder(node.right, visitor)
}

func inorder(node *Node, visitor nodeVisitor) {
	if node == nil {
		return
	}
	inorder(node.left, visitor)
	visitor(node)
	inorder(node.right, visitor)
}

func postorder(node *Node, visitor nodeVisitor) {
	if node == nil {
		return
	}
	postorder(node.left, visitor)
	postorder(node.right, visitor)
	visitor(node)
}

// formatBFS prints the tree vertically using BFS, using a simple procedural
// spacing algorithm to equally distribute the nodes at a given level. This isn't
// the tightest format to visualize parent-child relationships, but is useful
// for manual testing.
func (t *AvlTree) FormatBFS() string {
	if t.root == nil {
		return "<empty>"
	}

	// Space-char is printed between/around nodes. It is often best to print a non-blank char to
	// prevent editors from chomping leading space or converting spaces to tabs, etc.
	spaceChar := "."
	// Node width is derived from this format: 5e+00,5e+00 which is from "%1.0e,%1.0e", or 3 or "%3d"
	nw := 3
	// Minimum width around nodes, i.e. at the deepest (most crowded) level of the tree.
	mw := 2
	// Tree depth, to determine the line width required to print the widest section of the tree.
	d := t.root.height + 1
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
				line.WriteString(spaceChar)
			}
			sb.WriteString(line.String() + "\n")
			line.Reset()
		}

		// Line-predecessors is the (maximum) number of preceding nodes on a line
		lp := int(nodeNumber) - nl
		// Absolute starting position is defined by the number of preceding nodes and padding space.
		as := lp*nw + (lp+1)*sl
		for line.Len() < (as - 1) {
			line.WriteString(spaceChar)
		}
		//ns := fmt.Sprintf("%1.0e", float64(node.data))
		ns := fmt.Sprintf("%3d", node.data)
		ns = strings.Replace(ns, " ", spaceChar, -1)
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
// Input   Returned value
// 0001         1
// 0100         3
// 1001011      7
// 1000000      7
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
			q.PushBack(bfsData{
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
// Returns nil if not found.
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
