package fib_test

import (
	"fmt"

	"github.com/soniakeys/fib"
)

type str string

func (s str) LT(s2 fib.Value) bool {
	return s < s2.(str)
}

func ExampleHeap_DecreaseKey() {
	h := &fib.Heap{}
	x := h.Insert(str("rat"))
	h.Insert(str("cat"))
	fmt.Println(h.Min())

	h.DecreaseKey(x, str("gnat"))
	fmt.Println(h.Min())

	h.DecreaseKey(x, str("bat"))
	fmt.Println(h.Min())
	// Output:
	// cat true
	// cat true
	// bat true
}
