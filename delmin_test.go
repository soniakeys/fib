// Public domain

package fib_test

import (
	"fmt"

	"github.com/soniakeys/fib"
)

type sval string

func (s sval) LT(s2 fib.Value) bool {
	return s < s2.(sval)
}

func ExampleHeap_DeleteMin() {
	h := &fib.Heap{}
	h.Insert(sval("rat"))
	h.Insert(sval("cat"))
	fmt.Println(h.Min())

	h.DeleteMin()
	fmt.Println(h.Min())

	h.DeleteMin()
	fmt.Println(h.Min())
	// Output:
	// cat true
	// rat true
	// <nil> false
}
