// Public domain

package fib_test

import (
	"fmt"

	"github.com/soniakeys/fib"
)

type val string

func (s val) LT(s2 fib.Value) bool {
	return s < s2.(val)
}

func ExampleNode_Value() {
	h := &fib.Heap{}
	x := h.Insert(val("rat"))
	fmt.Println(x.Value())
	// Output:
	// rat
}
