# Treap

A treap is a binary search tree (BST) that attempts to ensure O(log(n)) operations by implementing a randomization strategy to distribute nodes up/down such that it has the same structural distribution as a randomly-generated bst.
They are similar to a kd-tree in that nodes are inserted based on multidimensional keys, except that the treap itself generates half of the key.
The name is a combination of tree+heap, all of which was intentionally done to confuse programmers attempting to read hand-waving descriptions of the data structure. :P

Formally, a treap consists of nodes containing a comparable stored value, x, as well as a randomly-generated value, y.
The treap enforces both bst and heap-order structural properties:
1) Nodes obey left/right bst-ordering: node A's left subtree contains x' < A.x and its right subtree contains x' > A.x
2) Nodes obey min-heap order: A.y < parent.y

Nodes are inserted in bst-order based on their x value, and then assigned a randomly-generated y value by which the node is moved up/down to achieve heap order.

And thank God you are asking how this is even possible, since many authors (*cough cough* Weiss) simply handwave, despite the fact
that this insertion routine appears to break the bst-property: nodes cannot be inserted in BST order and then simply
percolated up/down to achieve heap order, since these two requirements are seemingly inconsistent with one another.
However, a node can be moved up/down using the appropriate tree rotations to do so, and as usual with trees,
a few examples suffice to demonstrate.

(Assuming that min-heap property is maintained over y.)
Cases when inserting a new node A with value x and random y, once its bst-position is found:
1) y < parent.Y, RotateUp:
    - make A's right child the left child of its parent
    - make parent the right child of A
    - maintain A's left child
2) y > parent.Y, RotateDown:
    - take the lesser y child of A, C, and make A the right child of C
    - repeat on A, until it is in a consistent location.

TODO: cornercase when y == parent.y or child.y






