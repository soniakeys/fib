package fib_test

import (
	"fmt"

	"github.com/soniakeys/fib"
)

type anim string

func (s anim) LT(s2 fib.Value) bool {
	return s < s2.(anim)
}

func ExampleHeap_Insert() {
	h := &fib.Heap{}
	h.Insert(anim("cat"))      // return value not used
	a := h.Insert(anim("rat")) // return value kept
	fmt.Println(h.Min())

	h.DecreaseKey(a, anim("bat")) // Node "a" used
	fmt.Println(h.Min())

	h.Delete(a) // Node "a" used
	fmt.Println(h.Min())
	// Output:
	// cat true
	// bat true
	// cat true
}
