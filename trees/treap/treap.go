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

type Treap struct {
	root *treapNode
}

type treapNode struct {
	val         int
	priority    int
	left, right *treapNode
}

// var priority_generator func() int = rand.Int
var priority_generator func() int = func() int {
	return rand.Int() % 1000
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

// Pretty prints a treap using 150 char wide terminal, which allows only modestly sized tree depths.
/*
	root
		c1
			A
				D
				Nil
			B
		c2
			E
			F


					root
				 1.1E9,4.2E5

		  0.1E9,6.2E5   2.1E9,7.2E5

    BFS visitor pattern:
	- writer maintains spacing info
	-

	Alg:
	- the number of a node tells its center position on a line among n peers, where n is a power of 2
	- some n positions will be missing; these should consume constant width
	- the spacing between nodes on a line is determined by the number of nodes to pack in
	   * very few nodes means lots of spacing
	   * very many nodes means tight spacing
	   * equal spacing between nodes on a line
	   * Result appearance:
	           				*
	        		*               *
                *       *       *       *
	          *   *   *   *   *   *   *   *
		* Maths:
		    Defs:
			    nc: the number of nodes on a line (always some power of 2, minus 1)
		   	    nw: the printed width of a node (including space printed for missing/null nodes)
			    lw: the max/total line width (e.g. 180, or infinite based on 2**d where d is tree depth)
			    sw: width between/around nodes
				mw: minimum width between/around nodes
				as: absolute start position for printing node
				nn: node number
		    Note that 'as' is the only parameter a printing routine needs.
		    Then:
				mw = 2
			    lw = 2**d * mw + (2**d+1) * mw
			    sw = (lw -  nc * nw) / (nc + 1)
				as = (nn - 2**(d-1)) * (sw + nw)
*/

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

func (t *Treap) formatBFS() string {
	// Node width is derived from this format: 5e+00,5e+00 which is from "%1.0e,%1.0e"
	nw := 11
	// Minimum width around nodes, i.e. at the deepest (most crowded) level of the tree.
	mw := 2
	// Tree depth
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
		lp := (int(nodeNumber) - nl)
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
		// relations can be known, since a nodes left child is 2*number and right child
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

	return t.insert(val, t.root, &t.root)
}

// TODO: if this alg works, simplify by passing only parentLink, since it also contains @node as its value.
// TODO: what if priorities are not unique?
func (t *Treap) insert(val int, node *treapNode, parentLink **treapNode) error {
	// TODO: remove me or simplify nil cases
	if node == nil {
		panic("node cannot be nil in insert()")
	}

	if node.val == val {
		return ErrDuplicateValue
	}

	if val < node.val {
		if node.left != nil {
			if err := t.insert(val, node.left, &node.left); err != nil {
				return err
			}
			*parentLink = t.rotateLeftChild(node)
			return nil
		}

		node.left = &treapNode{
			val:      val,
			priority: priority_generator(),
		}
		*parentLink = t.rotateLeftChild(node)
	} else {
		// Case: val > node.val, so traverse right
		if node.right != nil {
			if err := t.insert(val, node.right, &node.right); err != nil {
				return err
			}
			*parentLink = t.rotateRightChild(node)
			return nil
		}

		node.right = &treapNode{
			val:      val,
			priority: priority_generator(),
		}
		*parentLink = t.rotateRightChild(node)
	}

	return nil
}

// TODO: simplify error cases, e.g. node.left must not be nil in this func
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
