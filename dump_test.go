package fib

import (
	"bytes"
	"fmt"
	"testing"
)

func (h Heap) validate(t *testing.T) {
	n := h.Node
	if n == nil {
		return
	}
	if n.parent != nil {
		t.Fatalf("root %v parent non-nil", n.value)
	}
	n.validateSibs(t) // recursive
	min := n
	for n = n.next; n != h.Node; n = n.next {
		if n.parent != nil {
			t.Fatalf("root %v parent non-nil", n.value)
		}
		if n.value.LT(min.value) {
			min = n
		}
	}
	if min != h.Node {
		t.Fatalf("heap min at %v but min sibling is %v", h.value, min.value)
	}
}

func (n *Node) validateSibs(t *testing.T) (numSibs int) {
	for x := n; ; {
		numSibs++
		if x.next.prev != x {
			t.Fatalf("node %v not sibling linked", x.value)
		}
		nch := 0
		if c := x.child; c != nil {
			if c.parent != x {
				t.Fatalf("node %v not parent linked", c.value)
			}
			if c.value.LT(x.value) {
				t.Fatalf("node %v LT parent %v", c.value, x.value)
			}
			nch = x.child.validateSibs(t) // recurse
		}
		if nch != x.rank {
			t.Fatalf("node %v stores rank=%d, but there are %d children",
				x.value, x.rank, nch)
		}
		x = x.next
		if x == n {
			break
		}
	}
	return
}

func (h Heap) dump() {
	if h.Node == nil {
		fmt.Println("empty heap")
		return
	}
	fmt.Printf("min value (top of heap) %v %p\n", h.value, h.Node)
	dump := func(x *Node) {
		fmt.Printf("%v %p\n", x.value, x)
		fmt.Printf("  parent, child, prev, next: %p %p %p %p\n",
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

func (n *Node) str() string {
	if n == nil {
		return "<nil>"
	}
	return fmt.Sprint(n.value)
}

func (h Heap) str() string {
	if h.Node == nil {
		return "empty heap"
	}
	b := &bytes.Buffer{}
	fmt.Fprintf(b, "min value (top of heap) %v\n", h.value)
	dump := func(x *Node) {
		fmt.Fprintf(b, "%s:\n", x.str())
		fmt.Fprintf(b, "  parent, child, prev, next: %s %s %s %s\n",
			x.parent.str(), x.child.str(), x.prev.str(), x.next.str())
		fmt.Fprintf(b, "  rank, mark: %d %t\n", x.rank, x.mark)
	}
	fmt.Fprintln(b, "roots:")
	q := []*Node{h.Node}
	dump(h.Node)
	for r := h.next; r != h.Node; r = r.next {
		q = append(q, r)
		dump(r)
	}
	// then breadth first traversal:
	for level := 1; len(q) > 0; level++ {
		fmt.Fprintf(b, "level %d:\n", level)
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
	return b.String()
}

// link c as a child of p
func link(p, c *Node) {
	c.parent = p
	if p.child == nil {
		c.next = c
		c.prev = c
		p.child = c
	} else {
		meld1(p.child, c)
	}
	p.rank++
}

type Int int

func (a Int) LT(b Value) bool { return a < b.(Int) }

func TestFig3(t *testing.T) {
	h := &Heap{}
	t.Run("empty", h.validate)

	p := h.Insert(Int(3))
	c := &Node{value: Int(4)}
	link(p, c)
	link(p, &Node{value: Int(5)})
	link(c, &Node{value: Int(14)})

	p = h.Insert(Int(6))
	c = &Node{value: Int(7)}
	link(p, c)
	link(p, &Node{value: Int(18)})
	link(c, &Node{value: Int(11)})

	p = h.Insert(Int(8))
	link(p, &Node{value: Int(10)})

	h.Insert(Int(12))
	t.Run("populated", h.validate)

	got := fmt.Sprint(h.DeleteMin())
	want := "3 true"
	if got != want {
		t.Fatal("got: ", got, ", want: ", want)
	}
	t.Run("result", h.validate)

	got = h.str()
	want = `min value (top of heap) 4
roots:
4:
  parent, child, prev, next: <nil> 14 5 5
  rank, mark: 3 false
5:
  parent, child, prev, next: <nil> 12 4 4
  rank, mark: 1 false
level 1:
14:
  parent, child, prev, next: 4 <nil> 6 8
  rank, mark: 0 false
8:
  parent, child, prev, next: 4 10 14 6
  rank, mark: 1 false
6:
  parent, child, prev, next: 4 7 8 14
  rank, mark: 2 false
12:
  parent, child, prev, next: 5 <nil> 12 12
  rank, mark: 0 false
level 2:
10:
  parent, child, prev, next: 8 <nil> 10 10
  rank, mark: 0 false
7:
  parent, child, prev, next: 6 11 18 18
  rank, mark: 1 false
18:
  parent, child, prev, next: 6 <nil> 7 7
  rank, mark: 0 false
level 3:
11:
  parent, child, prev, next: 7 <nil> 11 11
  rank, mark: 0 false
level 4:
`
	if got != want {
		t.Fatal("got: ", got, ", want: ", want)
	}
}

func TestFig5(t *testing.T) {
	h := &Heap{}
	p := h.Insert(Int(4))
	n7 := &Node{value: Int(7)}
	link(p, n7)
	link(p, &Node{value: Int(8)})
	p = n7
	link(p, &Node{value: Int(12)})
	c := &Node{value: Int(10)}
	link(p, c)
	link(p, &Node{value: Int(9)})
	link(c, &Node{value: Int(15)})
	p = h.Insert(Int(3))
	link(p, &Node{value: Int(14)})
	link(p, &Node{value: Int(5)})
	t.Run("(a)", h.validate)

	if err := h.DecreaseKey(c, Int(6)); err != nil {
		t.Fatal(err)
	}
	t.Run("(b)", h.validate)
	h.Delete(n7)
	t.Run("(c)", h.validate)
	got := h.str()
	want := `min value (top of heap) 3
roots:
3:
  parent, child, prev, next: <nil> 14 12 4
  rank, mark: 2 false
4:
  parent, child, prev, next: <nil> 8 3 6
  rank, mark: 1 false
6:
  parent, child, prev, next: <nil> 15 4 9
  rank, mark: 1 false
9:
  parent, child, prev, next: <nil> <nil> 6 12
  rank, mark: 0 false
12:
  parent, child, prev, next: <nil> <nil> 9 3
  rank, mark: 0 false
level 1:
14:
  parent, child, prev, next: 3 <nil> 5 5
  rank, mark: 0 false
5:
  parent, child, prev, next: 3 <nil> 14 14
  rank, mark: 0 false
8:
  parent, child, prev, next: 4 <nil> 8 8
  rank, mark: 0 false
15:
  parent, child, prev, next: 6 <nil> 15 15
  rank, mark: 0 false
level 2:
`
	if got != want {
		t.Fatal("got: ", got, ", want: ", want)
	}
}

func TestFig6(t *testing.T) {
	h := &Heap{}
	p := h.Insert(Int(2))
	c := &Node{value: Int(4), mark: true}
	link(p, c)
	link(p, &Node{value: Int(20)})
	p = c // 4
	c = &Node{value: Int(5), mark: true}
	link(p, c)
	link(p, &Node{value: Int(8)})
	link(p, &Node{value: Int(11)})
	p = c // 5
	c = &Node{value: Int(9), mark: true}
	link(p, c)
	link(p, &Node{value: Int(6)})
	link(p, &Node{value: Int(14)})
	p = c // 9
	c = &Node{value: Int(10)}
	link(p, c)
	link(p, &Node{value: Int(16)})
	link(c, &Node{value: Int(12)})
	link(c, &Node{value: Int(15)})
	t.Run("(a)", h.validate)
	if err := h.DecreaseKey(c, Int(7)); err != nil {
		t.Fatal(err)
	}
	t.Run("(b)", h.validate)
	got := h.str()
	want := `min value (top of heap) 2
roots:
2:
  parent, child, prev, next: <nil> 20 7 4
  rank, mark: 1 false
4:
  parent, child, prev, next: <nil> 8 2 5
  rank, mark: 2 true
5:
  parent, child, prev, next: <nil> 6 4 9
  rank, mark: 2 true
9:
  parent, child, prev, next: <nil> 16 5 7
  rank, mark: 1 true
7:
  parent, child, prev, next: <nil> 12 9 2
  rank, mark: 2 false
level 1:
20:
  parent, child, prev, next: 2 <nil> 20 20
  rank, mark: 0 false
8:
  parent, child, prev, next: 4 <nil> 11 11
  rank, mark: 0 false
11:
  parent, child, prev, next: 4 <nil> 8 8
  rank, mark: 0 false
6:
  parent, child, prev, next: 5 <nil> 14 14
  rank, mark: 0 false
14:
  parent, child, prev, next: 5 <nil> 6 6
  rank, mark: 0 false
16:
  parent, child, prev, next: 9 <nil> 16 16
  rank, mark: 0 false
12:
  parent, child, prev, next: 7 <nil> 15 15
  rank, mark: 0 false
15:
  parent, child, prev, next: 7 <nil> 12 12
  rank, mark: 0 false
level 2:
`
	if got != want {
		t.Fatal("got: ", got, ", want: ", want)
	}
}
