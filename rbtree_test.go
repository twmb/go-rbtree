package rbtree

import (
	"math/rand"
	"testing"
)

type myInt int

func (l myInt) Less(r Item) bool {
	return l < r.(myInt)
}

func TestRandom(t *testing.T) {
	const bound = 50
	rng := rand.New(rand.NewSource(0))
	var r Tree
	for i := 0; i < 10000; i++ {
		r.Insert(myInt(rng.Intn(bound)))
	}
	for i := 0; i < 500; i++ {
		if rng.Intn(3) == 0 {
			if found := r.Find(myInt(rng.Intn(bound))); found != nil {
				r.Delete(found)
			}
		} else {
			r.Insert(myInt(rng.Intn(bound)))
		}
	}
	for r.Len() > 0 {
		if rng.Intn(2) == 0 {
			r.Delete(r.Min())
		} else {
			r.Delete(r.Max())
		}
	}

	if r.Min() != nil {
		t.Error("expected len 0 min to be nil")
	}
	if r.Max() != nil {
		t.Error("expected len 0 max to be nil")
	}
}

func TestFind(t *testing.T) {
	var r Tree
	for i := 0; i < 20; i++ {
		r.Insert(myInt(i))
	}
	for i := 20 - 1; i >= 0; i-- {
		max := r.Max()
		if max == nil || max.Item.(myInt) != myInt(i) {
			t.Fatal("could not get max")
		}
		found := r.Find(max.Item.(myInt))
		if found == nil {
			t.Fatal("could not find max")
		}
		r.Delete(found)
	}

	if r.Find(myInt(21)) != nil {
		t.Error("found unexpected 21")
	}
}

func TestFindFn(t *testing.T) {
	var r Tree
	for i := 0; i < 20; i++ {
		r.Insert(myInt(i))
	}
	for i := 0; i < 20; i++ {
		n := r.FindFn(func(n *Node) int {
			v := int(n.Item.(myInt))
			return i - v
		})
		if n == nil || n.Item.(myInt) != myInt(i) {
			t.Fatalf("did not find %v", i)
		}
	}
	if r.FindFn(func(*Node) int { return -1 }) != nil {
		t.Error("found item when always left")
	}
	if r.FindFn(func(*Node) int { return 1 }) != nil {
		t.Error("found item when always right")
	}
}

type intPtr int

func (l *intPtr) Less(r Item) bool {
	return *l < *r.(*intPtr)
}

func TestFix(t *testing.T) {
	newIntPtr := func(v int) *intPtr {
		i := intPtr(v)
		return &i
	}

	var r Tree
	r.Insert(newIntPtr(1))
	r.Insert(newIntPtr(2))
	r.Insert(newIntPtr(3))
	r.Insert(newIntPtr(4))
	r.Insert(newIntPtr(9))
	r.Insert(newIntPtr(8))
	r.Insert(newIntPtr(7))
	r.Insert(newIntPtr(6))
	r.Insert(newIntPtr(5))

	max := r.Max()
	*max.Item.(*intPtr) = 0
	r.Fix(max)

	var exp intPtr
	for iter := IterAt(r.Min()); iter.Ok(); iter.Left() {
		if got := *iter.Item().(*intPtr); got != exp {
			t.Errorf("got %d != exp %d", got, exp)
		}
		exp++
	}
}

func TestIter(t *testing.T) {
	var r Tree
	r.Insert(myInt(0))
	r.Insert(myInt(1))
	r.Insert(myInt(2))
	r.Insert(myInt(3))
	r.Insert(myInt(4))
	r.Insert(myInt(5))
	r.Insert(myInt(6))
	r.Insert(myInt(7))
	r.Insert(myInt(8))
	r.Insert(myInt(9))
	const end = 10

	iter := IterAt(r.Before(r.Min()))
	for exp := 0; exp < end; exp++ {
		if got := iter.Right().Item.(myInt); got != myInt(exp) {
			t.Errorf("got %d != exp %d", got, exp)
		}
	}

	iter.Reset(r.After(r.Max()))
	for exp := end - 1; exp >= 0; exp-- {
		if got := iter.Left().Item.(myInt); got != myInt(exp) {
			t.Errorf("got %d != exp %d", got, exp)
		}
	}

	iter.Reset(r.Min())
	for exp := 0; exp < end; exp++ {
		if got := iter.Item().(myInt); got != myInt(exp) {
			t.Errorf("got %d != exp %d", got, exp)
		}
		iter.Right()
	}

	iter.Reset(r.Max())
	for exp := end - 1; exp >= 0; exp-- {
		if got := iter.Item().(myInt); got != myInt(exp) {
			t.Errorf("got %d != exp %d", got, exp)
		}
		iter.Left()
	}

	var exp int
	for iter := IterAt(r.Min()); iter.Ok(); iter.Right() {
		if got := iter.Item().(myInt); got != myInt(exp) {
			t.Errorf("got %d != exp %d", got, exp)
		}
		exp++
	}

	iter.Reset(r.After(r.Max()))
	peek, left := iter.PeekLeft(), iter.Left()
	if peek.Item != left.Item {
		t.Error("destructive peek left")
	}
	if peek.Item != myInt(end-1) {
		t.Errorf("got bad peek left from max %d", peek.Item)
	}

	iter.Reset(r.Before(r.Min()))
	peek, right := iter.PeekRight(), iter.Right()
	if peek.Item != right.Item {
		t.Error("destructive peek right")
	}
	if peek.Item != myInt(0) {
		t.Errorf("got bad peek right from min %d", peek.Item)
	}
}

func BenchmarkInsertCase4(b *testing.B) {
	b.Run("baseline", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var r Tree
			r.Insert(myInt(1))
			r.Insert(myInt(2))
		}
	})
	b.Run("thru_case4", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var r Tree
			r.Insert(myInt(1))
			r.Insert(myInt(2))
			r.Insert(myInt(3))
		}
	})
}
