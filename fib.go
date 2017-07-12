// Public domain

// Fib implements a Fibonacci heap.
//
// Implementation follows Fredman and Tarjan's "Fibonacci Heaps and Their Uses
// in Improved Network Optimization Algorithms", JACM 34:3, 1987.
package fib

import "errors"

// Value is an interface for a value stored in the heap.  Fredman and Tarjan
// call a value an "item" with a "real-valued" key.  While a floating point
// key could be prescribed, it is sufficient that the values are orderable.
//
// Stored values must simply implement this less-than method, LT.
//
// An implementation of LT on a type will typically type assert the argument
// to the receiver type and proceed with the comparison.
type Value interface {
	LT(Value) bool
}

// A node in a Fibonacci heap, holding a single value.
//
// Fredman and Tarjan have a parameter i for functions insert, decrease key,
// and delete.  For decrease key and delete, they add "This operation assumes
// that the position of i in h is known."  This clarifies that the argument
// of insert is a value, or an "item with a key" as F&T say, but that the
// argument of decrease key and delete must allow a reference to a value
// within the heap.  A *Node represents such a reference to a value.
// A *Node is thus returned by Insert and passed to DecreaseKey and Delete.
//
// The Node passed to DecreaseKey or Delete must be a Node created in and
// still present in the receiver heap.  Otherwise the Heap will likely be
// corrupted.  For tracking whether a Node is still in the Heap, remember
// that DeleteMin also removes Nodes.
type Node struct {
	value      Value
	parent     *Node // CLRS, Fredman and Tarjan use simply "p"
	child      *Node
	prev, next *Node // CLRS, Fredman and Tarjan use "left", "right"
	rank       int   // CLRS, Wikipedia use "degree"
	mark       bool
}

// Value is an accessor, or getter, for the Value stored in a Node.
func (n Node) Value() Value { return n.value }

// Heap represents a Fibonacci heap.
//
// The zero value of Heap is a valid empty Fibonacci heap.
// There is no constructor provided.  Use {} or new.
// To test if Heap h is empty, test h.Node == nil.
//
// (Note that while a Heap consisting of a nil *Node is valid, a nil *Heap
// is not a valid Fibonacci heap and will panic most Heap methods.)
type Heap struct{ *Node }

// Insert creates a new Node for Value v, adds it to receiver Heap h, and
// returns the newly created Node.
//
// Keep the return value if you might need to pass it to DecreaseKey or
// Delete.
func (h *Heap) Insert(v Value) *Node {
	x := &Node{value: v}
	if h.Node == nil {
		x.next = x
		x.prev = x
		h.Node = x
	} else {
		meld1(h.Node, x)
		if x.value.LT(h.value) {
			h.Node = x
		}
	}
	return x
}

// add a single node to a non-empty list.
// the added node does not need self-linked next and prev pointers
func meld1(list, single *Node) {
	list.prev.next = single
	single.prev = list.prev
	single.next = list
	list.prev = single
}

// Meld merges two Heaps.
//
// Meld merges all nodes of h2 into h.  Heap h2 is left empty.
//
// The two heaps must be different heaps.  Melding a heap to itself
// will corrupt the heap.
func (h *Heap) Meld(h2 *Heap) {
	switch {
	case h.Node == nil:
		*h = *h2
	case h2.Node != nil:
		meld2(h.Node, h2.Node)
		if h2.value.LT(h.value) {
			*h = *h2
		}
	}
	h2.Node = nil
}

// meld two non-empty node lists
func meld2(a, b *Node) {
	a.prev.next = b
	b.prev.next = a
	a.prev, b.prev = b.prev, a.prev
}

// Min returns the minimum value in a Heap.
//
// It returns the minimum and ok = true as long as the heap is not empty.
// Otherwise it returns a nil Value and ok = false.
func (h Heap) Min() (min Value, ok bool) {
	if h.Node == nil {
		return
	}
	return h.value, true
}

// DeleteMin deletes the miminum value from a heap.
//
// It returns the deleted minimum and ok = true as long as the heap is not
// empty.  Otherwise the heap is left empty and the method returns a nil
// Value interface and ok = false.
//
// For a non-empty Heap, the minimum value is stored in the Node at h.Node
// and it is that exact Node that will be removed from the Heap.
func (h *Heap) DeleteMin() (min Value, ok bool) {
	if h.Node == nil {
		return
	}
	min = h.value // return value

	// "Linking Step" of F&T
	// F&T and CLRS both reference n, a total number of nodes in the heap
	// and suggest a function of log(n) as a bound for an array of root nodes
	// with unique rank.  Code here uses a map instead which still gives
	// amortized O(1) access time but is simpler and safer.
	roots := map[int]*Node{}
	add := func(r *Node) {
		r.prev = r
		r.next = r
		for {
			x, ok := roots[r.rank]
			if !ok {
				break
			}
			delete(roots, r.rank)
			// r, x are single Nodes with same rank.  "link" them.
			if x.value.LT(r.value) {
				r, x = x, r
			}
			// r has minimum Value. meld x with children of r
			x.parent = r
			x.mark = false
			if r.child == nil {
				x.next = x
				x.prev = x
				r.child = x
			} else {
				meld1(r.child, x)
			}
			r.rank++
		}
		roots[r.rank] = r
	}
	// add operations are performed on a virtual list consisting of the
	// root nodes other than the minimum node and the children of the
	// minimum node.  there is no point in actually constructing this list
	// as it is processed sequentially by adding nodes to the "roots" map.
	for r := h.next; r != h.Node; {
		n := r.next
		add(r)
		r = n
	}
	// add any children of minimum
	if c := h.child; c != nil {
		c.parent = nil
		r := c.next
		add(c)
		for r != c {
			n := r.next
			r.parent = nil
			add(r)
			r = n
		}
	}
	if len(roots) == 0 {
		h.Node = nil
		return min, true
	}
	// link roots, finding one with (new) minimum value
	var mv *Node
	var d int
	for d, mv = range roots {
		break
	}
	delete(roots, d)
	mv.next = mv
	mv.prev = mv
	for _, r := range roots {
		r.prev = mv
		r.next = mv.next
		mv.next.prev = r
		mv.next = r
		if r.value.LT(mv.value) {
			mv = r
		}
	}
	h.Node = mv      // set receiver to new min
	return min, true // return old min
}

// DecreaseKey stores a new Value in Node n.
//
// Node n must be a node in Heap h.  The new value v must be less than or
// equal to the existing value.
//
// If the existing value is LT the new value, the method returns an error.
func (h *Heap) DecreaseKey(n *Node, v Value) error {
	if n.value.LT(v) {
		return errors.New("DecreaseKey new value greater than existing value")
	}
	n.value = v      // store it
	if n == h.Node { // if it was min before, it's still min.
		return nil
	}
	p := n.parent
	if p == nil {
		if v.LT(h.value) {
			h.Node = n
		}
		return nil
	}
	h.cutAndMeld(n)
	return nil
}

func (h Heap) cut(x *Node) {
	// cut loc from parent
	p := x.parent
	p.rank--
	if p.rank == 0 {
		p.child = nil
	} else {
		p.child = x.next
		x.prev.next = x.next
		x.next.prev = x.prev
	}
	if p.parent == nil {
		return
	}
	if !p.mark { // parent is losing first child: mark
		p.mark = true
		return
	}
	// parent is losing second child: cascade
	h.cutAndMeld(p)
}

func (h Heap) cutAndMeld(x *Node) {
	h.cut(x)
	x.parent = nil
	meld1(h.Node, x)
}

// Delete removes the specified node n from heap h.
//
// Node n must be a node in Heap h.
func (h *Heap) Delete(n *Node) {
	// Delete here follows F&T rather than CLSR.  CLSR takes the easy route
	// of always going through DeleteMin.  F&T only calls DeleteMin if the
	// node being deleted is indeed the minimum.  They claim that keeps it
	// O(1) unless the minimum is being deleted.  I'm not sure it's quite O(1)
	// becase parent pointers must be set to nill for all the children
	// that become roots, but that's still cheaper than calling DeleteMin.
	// The two approaches remain fundamentally different because F&T avoid
	// any linking/consolidate steps unless it is actually the minimum that
	// is being deleted.
	p := n.parent
	if p == nil {
		if n == h.Node {
			h.DeleteMin()
			return
		}
		n.prev.next = n.next
		n.next.prev = n.prev
	} else {
		h.cut(n) // cut n from parent, but don't add it as a root
	}
	c := n.child
	if c == nil {
		return
	}
	// add children as roots
	for {
		c.parent = nil
		c = c.next
		if c == n.child {
			break
		}
	}
	meld2(h.Node, c)
}
