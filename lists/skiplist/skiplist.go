package skiplist

import (
	"errors"
	"math"
	"math/rand"
)

// Skiplist is an ordered list for which nodes possess up to r forward
// pointers, each 'skipping' to the next node with r' >= r, where r is some
// small number. All nodes at least point to their immediate sibling.
// For example, if r == 4 for some node, then the node possesses four
// forward pointers: a traditional 'next' sibling pointer, one skipping
// to the next node for which r' >= 2, one skipping to the next r' >= 3,
// and one pointing to the next node for which r' >= 4. Only a diagram makes
// any sense, but note the pointers behave like a light ray, attaching to the
// next header block that 'blocks' them.
//
// This strategy ensures that a list has the same average search and
// insertion complexity as a randomized bst, such as a treap. Searching is
// done rapidly using quasi binary search.
//
// Usage and applications: https://en.wikipedia.org/wiki/Skip_list
//
// Diagram for a skiplist with r = 4:
//
// -        | 3 | ------------> | 3 | --------------------------------> | 3 | ------------> nil
// - Node   | 2 | ------------> | 2 | ------------> | 2 | ------------> | 2 | --> | 2 | --> nil
// - Ptrs   | 1 | ------------> | 1 | --> | 1 | --> | 1 | ------------> | 1 | --> | 1 | --> nil
// -        | 0 | --> | 0 | --> | 0 | --> | 0 | --> | 0 | --> | 0 | --> | 0 | --> | 0 | --> nil
// - Vals:    2         5         9         19        45        47        54        62
//
//	Note the invariant that the first node contains R pointers, where R is the maximum
//	number of pointers. This isn't strictly necessary, but simplifies code. This can
//	also be implemented via a permanent empty dummy-start node, e.g. containing only
//	a sentinel value such as -int.MaxInt.
//
// A personal misconception of my own is that the rank, r, of a pointer relates to the
// a pow(2,r) number of nodes to skip. There is no such constraint. As shown, the
// pointers merely point to the next node of equal or greater r. However, the probability
// distribution does ensure that pow(2,r) is the average distribution of link hops,
// which is why the data structure has O(lg(n)). Likewise, binary search is implemented
// without a parititon point, but merely by looking ahead from the highest rank pointer
// down to the lower ranked pointer.
type Skiplist struct {
	root *skipNode
	r    int
}

type skipNode struct {
	next  []*skipNode
	value int
}

var (
	ErrDuplicateValue error      = errors.New("duplicate value")
	ErrValueNotFound  error      = errors.New("value not found")
	rand_generator    func() int = rand.Int
)

// A new SkipList is created with a fixed r, much like a hash-table's size
// is determined up front. A dynamically-sized r-Skiplist could be implemented,
// much the same as rehashing is performed on a hashtable when it reaches a
// certain load factor; for Skiplists this would be about r=lg(N).
func NewSkiplist(r int) *Skiplist {
	return &Skiplist{
		r: r,
		root: &skipNode{
			next: make([]*skipNode, r),
			// Note: this is a palliative; max-int is merely in valid-but-unlikely values range.
			// It is best to code defensively such that the sentinal/root node's value is never
			// evaluated at all, since the node itself is permanent and can be handled distinctly.
			value: -math.MaxInt,
		},
	}
}

// TODO: this is just a demo, since it is obviously redundant to search
// by value... for the same value, haha. The Skiplist should be abstracted
// to match list items based on an Id() interface or other comparable mechanism.
func (sl *Skiplist) Get(n int) (int, error) {
	ptrs := sl.search(n)
	if ptrs[0] == nil || ptrs[0].value != n {
		return 0, ErrValueNotFound
	}

	return ptrs[0].value, nil
}

/*
Search populates and returns a pointer slice of size r, for which each
entry is the first node of that rank prior to n in the list ordering.

For straightforward search, the 0th value in the slice contains the node
less than or equal to the value. Hence the 0th index always contains:
- the value searched for, if it exists
- the first node prior that node's ordered location, if it does not.
*/
func (sl *Skiplist) search(n int) []*skipNode {
	node := sl.root.next[0]
	if node == nil {
		// List is empty, since sentinel/root points to nil.
		return sl.root.next
	}

	ptrs := make([]*skipNode, sl.r)
	for rank := sl.r - 1; rank >= 0; rank-- {
		// Search for the last node at this level prior to the passed value, or nil
		for node.next[rank] != nil && node.next[rank].value < n {
			node = node.next[rank]
		}
		ptrs[rank] = node
	}

	return ptrs
}

// Insert threads in a new node, whose header size is randomly generated in (0,r].
// Per skiplist structure, the new node's header entries are required to point
// to each next node for that entry's skip value.
func (sl *Skiplist) Insert(n int) error {
	ptrs := sl.search(n)
	// TODO: when would ptrs[0] be nil?
	if ptrs[0] != nil && ptrs[0].value == n {
		return ErrDuplicateValue
	}

	hdrSize := rand_generator() % sl.r
	newNode := &skipNode{
		next: make([]*skipNode, hdrSize),
	}

	// Thread the new node into the previous node's headers,
	// only up to hdrSize in the ptr array.
	for i := 0; i < hdrSize; i++ {
		newNode.next[i] = ptrs[i].next[i]
		ptrs[i].next[i] = newNode
	}

	return nil
}

// Delete removes a node from the skiplist.
// Deletion is merely the inverse of insertion: point
// all parent pointers to one's children, even if they are nil.
func (sl *Skiplist) Delete(n int) error {
	ptrs := sl.search(n)
	// When is ptrs[0] nil?
	if ptrs[0] == nil || ptrs[0].value != n {
		return ErrValueNotFound
	}

	// Thread all parent pointers to the node's successors
	target := ptrs[0].next[0]
	for i := 0; i < len(target.next); i++ {
		ptrs[i].next[i] = target.next[i]
		// Nillify all ptrs to prevent mem leaks and release memory
		target.next[i] = nil
	}

	return nil
}