package logic

import (
	"testing"
	"testing/quick"
)

func TestNot(t *testing.T) {
	l := NewNot(NewLeaf("a"))
	f := func(a bool) bool {
			m := map[string]bool {
				"a": a,
			}
			return l.Eval(m) == !a
		}
	if e := quick.Check(f, nil); e != nil {
		t.Fatalf("%s", e)
	}
}

func TestAnd(t *testing.T) {
	l := NewAnd(NewLeaf("a"), NewLeaf("b"), NewLeaf("c"), NewLeaf("d"))
	f := func(a, b, c, d bool) bool {
			m := map[string]bool {
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
	l := NewOr(NewLeaf("a"), NewLeaf("b"), NewLeaf("c"), NewLeaf("d"))
	f := func(a, b, c, d bool) bool {
			m := map[string]bool {
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
	l := NewIff(NewLeaf("a"), NewLeaf("b"), NewLeaf("c"), NewLeaf("d"))
	f := func(a, b, c, d bool) bool {
			m := map[string]bool {
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
