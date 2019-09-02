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

func TestFindWith(t *testing.T) {
	var r Tree
	for i := 0; i < 20; i++ {
		r.Insert(myInt(i))
	}
	for i := 0; i < 20; i++ {
		n := r.FindWith(func(n *Node) int {
			v := int(n.Item.(myInt))
			return i - v
		})
		if n == nil || n.Item.(myInt) != myInt(i) {
			t.Fatalf("did not find %v", i)
		}
	}
	if r.FindWith(func(*Node) int { return -1 }) != nil {
		t.Error("found item when always left")
	}
	if r.FindWith(func(*Node) int { return 1 }) != nil {
		t.Error("found item when always right")
	}
}

func TestFindOrInsert(t *testing.T) {
	var r Tree
	for i := 0; i < 20; i++ {
		for j := 0; j < 10; j++ {
			node := r.FindOrInsert(myInt(i))
			if got := int(node.Item.(myInt)); got != i {
				t.Errorf("got insert %d != exp %d", got, i)
			}
		}
	}
	if got := r.Len(); got != 20 {
		t.Errorf("got len %d != exp %d", got, 20)
	}
	i := 0
	for it := IterAt(r.Min()); it.Ok(); it.Right() {
		if got := it.Item().(myInt); got != myInt(i) {
			t.Errorf("got %d != exp %d", got, i)
		}
		i++
	}
}

func TestFindWithOrInsertWith(t *testing.T) {
	var r Tree
	for i := 0; i < 20; i++ {
		for j := 0; j < 10; j++ {
			node := r.FindWithOrInsertWith(
				func(n *Node) int { return i - int(n.Item.(myInt)) },
				func() Item { return myInt(i) },
			)
			if got := int(node.Item.(myInt)); got != i {
				t.Errorf("got insert %d, exp %d", got, i)
			}
		}
	}
	if got := r.Len(); got != 20 {
		t.Errorf("got len %d != exp %d", got, 20)
	}
	i := 0
	for it := IterAt(r.Min()); it.Ok(); it.Right() {
		if got := it.Item().(myInt); got != myInt(i) {
			t.Errorf("got %d != exp %d", got, i)
		}
		i++
	}
}

type intPtr int

func (l *intPtr) Less(r Item) bool {
	return *l < *r.(*intPtr)
}

func newIntPtr(v int) *intPtr {
	i := intPtr(v)
	return &i
}

func TestFix(t *testing.T) {
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

	iter := IterAt(Into(r.Min()))
	for exp := 0; exp < end; exp++ {
		if got := iter.Right().Item.(myInt); got != myInt(exp) {
			t.Errorf("got %d != exp %d", got, exp)
		}
	}

	iter.Reset(Into(r.Max()))
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

	iter.Reset(Into(r.Max()))
	peek, left := iter.PeekLeft(), iter.Left()
	if peek.Item != left.Item {
		t.Error("destructive peek left")
	}
	if peek.Item != myInt(end-1) {
		t.Errorf("got bad peek left from max %d", peek.Item)
	}

	iter.Reset(Into(r.Min()))
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

func BenchmarkFindWithOrInsertWithExisting(b *testing.B) {
	var r Tree
	r.Insert(myInt(1))
	for i := 0; i < b.N; i++ {
		r.FindWithOrInsertWith(
			func(n *Node) int { return 1 - int(n.Item.(myInt)) },
			func() Item { return myInt(1) },
		)
	}
}

func findNum(num int) func(*Node) int {
	return func(n *Node) int { return num - int(n.Item.(myInt)) }
}

func newNum(num int) func() Item {
	return func() Item { return myInt(num) }
}

func BenchmarkFindWithOrInsertWithExistingClosure(b *testing.B) {
	var r Tree
	r.Insert(myInt(1))
	for i := 0; i < b.N; i++ {
		r.FindWithOrInsertWith(
			findNum(1),
			newNum(1),
		)
	}
}

func BenchmarkFindWithOrInsertWithNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r Tree
		r.FindWithOrInsertWith(
			func(n *Node) int { return 1 - int(n.Item.(myInt)) },
			func() Item { return myInt(1) },
		)
	}
}
