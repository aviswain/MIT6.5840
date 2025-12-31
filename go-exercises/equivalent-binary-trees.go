package main

import (
	"golang.org/x/tour/tree"
	"fmt"
)

func walkInOrder(t *tree.Tree, ch chan int) {
	if t == nil {
		return
	}

	walkInOrder(t.Left, ch)
	ch <- t.Value
	walkInOrder(t.Right, ch)
}
// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
	go func() {
		walkInOrder(t, ch)
		close(ch)
	}()
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
	ch1 := make(chan int)
	ch2 := make(chan int)

	Walk(t1, ch1)
	Walk(t2, ch2)

	for {
		val1, ok1 := <-ch1
		val2, ok2 := <-ch2

		// Both channels close at the same time...
		if !ok1 && !ok2 {
			return true
		}

		// Only one channel closed — different sizes
		if !ok1 || !ok2 {
			return false
		}

		if val1 != val2 {
			return false
		}
	}
}

func main() {
	fmt.Println(Same(tree.New(1), tree.New(1)))
	fmt.Println(Same(tree.New(1), tree.New(2)))
}
