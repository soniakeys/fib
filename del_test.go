package fib_test

import (
	"fmt"

	"github.com/soniakeys/fib"
)

type sv string

func (s sv) LT(s2 fib.Value) bool {
	return s < s2.(sv)
}

func ExampleHeap_Delete() {
	h := &fib.Heap{}
	r := h.Insert(sv("rat"))
	c := h.Insert(sv("cat"))
	fmt.Println(h.Min())

	h.Delete(r)
	fmt.Println(h.Min())

	h.Delete(c)
	fmt.Println(h.Min())
	// Output:
	// cat true
	// cat true
	// <nil> false
}
