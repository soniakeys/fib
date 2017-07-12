// Public domain

package fib_test

import (
	"fmt"

	"github.com/soniakeys/fib"
)

func ExampleHeap() {
	h := &fib.Heap{}
	fmt.Println(h.Node == nil)
	// Output:
	// true
}
