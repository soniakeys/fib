// Public domain

package fib_test

import (
	"fmt"

	"github.com/soniakeys/fib"
)

type at string

func (s at) LT(s2 fib.Value) bool {
	return s < s2.(at)
}

func ExampleHeap_Meld() {
	h := &fib.Heap{}
	h.Insert(str("rat"))
	h.Insert(str("cat"))
	fmt.Println(h.Min())

	h2 := &fib.Heap{}
	h2.Insert(str("bat"))
	h2.Insert(str("gnat"))
	fmt.Println(h2.Min())

	h.Meld(h2)
	fmt.Println(h.Min())
	fmt.Println(h2.Min())
	// Output:
	// cat true
	// bat true
	// bat true
	// <nil> false
}
