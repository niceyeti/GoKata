package treap

import (
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

var priority_generator func() int = rand.Int

func init() {
	rand.Seed(time.Now().UnixNano())
}

type TraversalOrder int

const (
	PreOrder TraversalOrder = iota + 1
	PostOrder
	InOrder
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
	default:
		return "", ErrNoSuchTraversalOrder
	}

	return sb.String(), nil
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
			return t.insert(val, node.left, &node.left)
		}

		node.left = &treapNode{
			val:      val,
			priority: priority_generator(),
		}
		*parentLink = t.rotateLeftChild(node)
	} else {
		// Case: val > node.val, so traverse right
		if node.right != nil {
			return t.insert(val, node.right, &node.right)
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
