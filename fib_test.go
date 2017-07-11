package fib_test

import (
	"fmt"

	"github.com/soniakeys/fib"
)

func ExampleHeap_Min_empty() {
	h := &fib.Heap{}
	fmt.Println(h.Min())
	// Output:
	// <nil> false
}

type str string

func (s str) LT(s2 fib.Value) bool {
	return s < s2.(str)
}

func ExampleHeap_Insert() {
	h := &fib.Heap{}
	h.Insert(str("bamp"))
	fmt.Println(h.Min())
	// Output:
	// bamp true
}

func ExampleHeap_Meld() {
	// merge two empty heaps
	h1 := &fib.Heap{}
	h2 := &fib.Heap{}
	h1.Meld(h2)
	fmt.Println(h1)
	// merge non-empty into empty
	h2.Insert(str("bamp"))
	h1.Meld(h2)
	fmt.Println(h1.Min())
	// merge empty into non-empty
	h2 = &fib.Heap{}
	h1.Meld(h2)
	fmt.Println(h1.Min())
	// merge non-empty into non-empty
	h2.Insert(str("ash"))
	h1.Meld(h2)
	fmt.Println(h1.Min())

	// Output:
	// &{<nil>}
	// bamp true
	// bamp true
	// ash true
}
