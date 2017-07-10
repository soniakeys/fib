package fib

import "fmt"

func (h Heap) dump() {
	if h.Node == nil {
		fmt.Println("empty heap")
		return
	}
	fmt.Printf("min value (top of heap) %v %x\n", h.value, h)
	dump := func(x *Node) {
		fmt.Printf("%v %x\n", x.value, x)
		fmt.Printf("  parent, child, prev, next: %x %x %x %x\n",
			x.parent, x.child, x.prev, x.next)
		fmt.Printf("  rank, mark: %d %t\n", x.rank, x.mark)
	}
	fmt.Println("roots:")
	q := []*Node{h.Node}
	dump(h.Node)
	for r := h.next; r != h.Node; r = r.next {
		q = append(q, r)
		dump(r)
	}
	// then breadth first traversal:
	for level := 1; len(q) > 0; level++ {
		fmt.Printf("level %d:\n", level)
		var q2 []*Node
		for _, h := range q {
			if h.child == nil {
				continue
			}
			q2 = append(q2, h.child)
			dump(h.child)
			for c := h.child.next; c != h.child; c = c.next {
				q2 = append(q2, c)
				dump(c)
			}
		}
		q = q2
	}
}
