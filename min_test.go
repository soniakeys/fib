// Public domain

package fib_test

import (
	"fmt"

	"github.com/soniakeys/fib"
)

type min string

func (s min) LT(s2 fib.Value) bool {
	return s < s2.(min)
}

func ExampleHeap_Min() {
	h := &fib.Heap{}
	fmt.Println(h.Min())

	h.Insert(min("cat"))
	fmt.Println(h.Min())

	h.Insert(min("rat"))
	fmt.Println(h.Min())
	// Output:
	// <nil> false
	// cat true
	// cat true
}
