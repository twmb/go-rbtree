package rbtree_test

import (
	"fmt"

	"github.com/twmb/go-rbtree"
)

type Int int

func (l Int) Less(r rbtree.Item) bool {
	return l < r.(Int)
}

func ExampleIter() {
	var r rbtree.Tree

	for i := 0; i < 10; i++ {
		r.Insert(Int(i)) // Int provides Less on int
	}

	// Declaring iter here to show we can reset it later.
	var iter rbtree.Iter

	// We can start iterating _at_ the max.
	fmt.Println("iterating down...")
	for iter = rbtree.IterAt(r.Max()); iter.Ok(); iter.Left() {
		fmt.Println(iter.Item())
	}

	// Or, we can start iterating just before the min so that
	// our first call to Right will start at the min.
	iter.Reset(rbtree.Before(r.Min()))
	fmt.Println("iterating up...")
	for {
		next := iter.Right()
		if next == nil { // when right returns nil, we are done
			break
		}
		fmt.Println(next.Item)
	}

	// Putting it together...
	fmt.Println("iterating down very inefficiently...")
	for it := rbtree.IterAt(r.Min()); it.Ok(); it.Right() {
		if it.Item() == r.Max().Item {
			fmt.Println(it.Item())
			r.Delete(it.Node())
			it.Reset(rbtree.Before(r.Min()))
		}
	}

	// Output:
	// iterating down...
	// 9
	// 8
	// 7
	// 6
	// 5
	// 4
	// 3
	// 2
	// 1
	// 0
	// iterating up...
	// 0
	// 1
	// 2
	// 3
	// 4
	// 5
	// 6
	// 7
	// 8
	// 9
	// iterating down very inefficiently...
	// 9
	// 8
	// 7
	// 6
	// 5
	// 4
	// 3
	// 2
	// 1
	// 0

}
