package logic

import (
	"testing"
	"testing/quick"
)

func TestNot(t *testing.T) {
	l := NewOperation(NOT, NewLeaf("a"))
	f := func(a bool) bool {
		m := map[string]bool{
			"a": a,
		}
		return l.Eval(m) == !a
	}
	if e := quick.Check(f, nil); e != nil {
		t.Fatalf("%s", e)
	}
}

func TestAnd(t *testing.T) {
	l := NewOperation(AND, NewLeaf("a"), NewLeaf("b"), NewLeaf("c"), NewLeaf("d"))
	f := func(a, b, c, d bool) bool {
		m := map[string]bool{
			"a": a,
			"b": b,
			"c": c,
			"d": d,
		}
		return l.Eval(m) == (a && b && c && d)
	}
	if e := quick.Check(f, nil); e != nil {
		t.Fatalf("%s", e)
	}
}

func TestOr(t *testing.T) {
	l := NewOperation(OR, NewLeaf("a"), NewLeaf("b"), NewLeaf("c"), NewLeaf("d"))
	f := func(a, b, c, d bool) bool {
		m := map[string]bool{
			"a": a,
			"b": b,
			"c": c,
			"d": d,
		}
		return l.Eval(m) == (a || b || c || d)
	}
	if e := quick.Check(f, nil); e != nil {
		t.Fatalf("%s", e)
	}
}

func TestIff(t *testing.T) {
	l := NewOperation(IFF, NewLeaf("a"), NewLeaf("b"), NewLeaf("c"), NewLeaf("d"))
	f := func(a, b, c, d bool) bool {
		m := map[string]bool{
			"a": a,
			"b": b,
			"c": c,
			"d": d,
		}
		return l.Eval(m) == (a == b == c == d)
	}
	if e := quick.Check(f, nil); e != nil {
		t.Fatalf("%s", e)
	}
}

// This is stupid!
func TestSimplify(t *testing.T) {
	l := NewOperation(NOT, NewLeaf("a"))
	r := Simplify(l)
	if r.String() != "!(a)" {
		t.Fatalf("Simplify(%s) returned %s", l, r)
	}

	l = NewOperation(IFF, NewOperation(NOT, NewLeaf("a")), NewOperation(OR, NewLeaf("a"), NewLeaf("b")))
	r = Simplify(l)
	if r.String() != "^(v(!(!(a)), v(a, b)), v(!(v(a, b)), !(a)))" {
		t.Fatalf("Simplify(%s) returned %s", l, r)
	}

	l = NewOperation(NOT, NewOperation(AND, NewOperation(OR, NewLeaf("a"), NewLeaf("b")), NewOperation(OR, NewLeaf("c"), NewLeaf("d"))))
	r = Simplify(l)
	if r.String() != "!(^(v(a, b), v(c, d)))" {
		t.Fatalf("Simplify(%s) returned %s", l, r)
	}

}

// This is stupid, as well. God, who am I?!
func TestDeMorgan(t *testing.T) {
	var l, r Node

	l = NewOperation(AND, NewLeaf("a"), NewLeaf("b"))
	r = DeMorgan(l)
	if r.String() != "^(a, b)" {
		t.Fatalf("DeMorgan(%s) returned %s", l, r)
	}

	l = NewOperation(NOT, NewLeaf("a"))
	r = DeMorgan(l)
	if r.String() != "!(a)" {
		t.Fatalf("DeMorgan(%s) returned %s", l, r)
	}

	l = NewOperation(NOT, NewOperation(NOT, NewLeaf("a")))
	r = DeMorgan(l)
	if r.String() != "a" {
		t.Fatalf("DeMorgan(%s) returned %s", l, r)
	}

	l = NewOperation(NOT, NewOperation(AND, NewLeaf("a"), NewLeaf("b")))
	r = DeMorgan(l)
	if r.String() != "v(!(a), !(b))" {
		t.Fatalf("DeMorgan(%s) returned %s", l, r)
	}

	l = NewOperation(NOT, NewOperation(AND, NewOperation(OR, NewLeaf("a"), NewLeaf("b")), NewOperation(OR, NewLeaf("c"), NewLeaf("d"))))
	r = DeMorgan(l)
	if r.String() != "v(^(!(a), !(b)), ^(!(c), !(d)))" {
		t.Fatalf("Simplify(%s) returned %s", l, r)
	}

}

func TestCNF(t *testing.T) {
	var l, r Node

	l = NewLeaf("a")
	r = CNF(l)
	if r.String() != "^(v(a))" {
		t.Fatalf("CNF(%s) returned %s", l, r)
	}

	l = NewOperation(NOT, NewLeaf("a"))
	r = CNF(l)
	if r.String() != "^(v(!(a)))" {
		t.Fatalf("CNF(%s) returned %s", l, r)
	}

	l = NewOperation(AND, NewLeaf("a"), NewLeaf("b"))
	r = CNF(l)
	if r.String() != "^(v(a), v(b))" {
		t.Fatalf("CNF(%s) returned %s", l, r)
	}

	l = NewOperation(OR, NewLeaf("a"), NewLeaf("b"))
	r = CNF(l)
	if r.String() != "^(v(a, b))" {
		t.Fatalf("CNF(%s) returned %s", l, r)
	}

	l = NewOperation(OR, NewOperation(AND, NewLeaf("a"), NewLeaf("b")), NewOperation(AND, NewLeaf("c"), NewLeaf("d")))
	r = CNF(l)
	if r.String() != "^(v(a, c), v(b, c), v(a, d), v(b, d))" {
		t.Fatalf("CNF(%s) returned %s", l, r)
	}

	l = NewOperation(AND, NewOperation(NOT, NewOperation(OR, NewLeaf("a"), NewLeaf("b"))), NewOperation(OR, NewLeaf("c"), NewLeaf("d")))
	r = CNF(l)
	if r.String() != "^(v(!(a)), v(!(b)), v(c, d))" {
		t.Fatalf("CNF(%s) returned %s", l, r)
	}

	l = NewOperation(NOT, NewOperation(AND, NewOperation(OR, NewLeaf("a"), NewLeaf("b")), NewOperation(OR, NewLeaf("c"), NewLeaf("d"))))
	r = CNF(l)
	if r.String() != "^(v(!(a), !(c)), v(!(b), !(c)), v(!(a), !(d)), v(!(b), !(d)))" {
		t.Fatalf("CNF(%s) returned %s", l, r)
	}

}
