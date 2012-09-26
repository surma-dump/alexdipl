package logic

func operandMap(parent *Operation, f func(Node) Node) *Operation {
	r := &Operation{}
	*r = *parent
	r.Operands = make([]Node, len(parent.Operands))
	for i := range parent.Operands {
		r.Operands[i] = f(parent.Operands[i])
	}
	return r
}

// Reduces expressions to using only AND, OR and NOT
func Simplify(n Node) Node {
	if x, ok := n.(*Operation); ok {
		if len(x.Operands) == 0 {
			return nil
		}
		if x.Operator != NOT && len(x.Operands) == 1 {
			return Simplify(x.Operands[0])
		}

		if x.Operator == IF {
			op1 := x.Operands[0]
			op2 := x.Operands[1]
			if len(x.Operands) > 2 {
				op2 = Simplify(op2)
			}
			// a -> b <=> (!a v b)
			return NewOperation(OR, NewOperation(NOT, op1), op2)
		}

		if x.Operator == IFF {
			op1 := x.Operands[0]
			op2 := x.Operands[1]
			if len(x.Operands) > 2 {
				op2 = Simplify(op2)
			}
			// a <-> b <=> (a -> b) ^ (b -> a)
			return Simplify(NewOperation(AND, NewOperation(IF, op1, op2), NewOperation(IF, op2, op1)))
		}
		return operandMap(x, Simplify)
	}
	return n
}

// Uses DeMorgan until NOTs are only applied to literals (i.e. leafs)
// Implies simplify, assumes simplifiability
func DeMorgan(n Node) Node {
	n = Simplify(n)
	if x, ok := n.(*Operation); ok {
		if len(x.Operands) == 0 {
			return nil
		}
		if x.Operator != NOT {
			return operandMap(x, DeMorgan)
		}
		if _, ok := x.Operands[0].(Leaf); ok {
			return x
		}
		x = x.Operands[0].(*Operation)
		switch x.Operator {
		case OR:
			r := NewOperation(AND, x.Operands...)
			return operandMap(r, func(n Node) Node {
				return DeMorgan(NewOperation(NOT, n))
			})
		case AND:
			r := NewOperation(OR, x.Operands...)
			return operandMap(r, func(n Node) Node {
				return DeMorgan(NewOperation(NOT, n))
			})
		case NOT:
			return DeMorgan(x.Operands[0])
		default:
			panic("Unexpected Operator type while DeMorganing")
		}
	}
	return n
}

// Converts expression to CNF
// Implies DeMorgan()
//
// Also: This function is ugly as shit. Kill it with fire.
func CNF(n Node) Node {
	n = DeMorgan(n)
	if x, ok := n.(*Operation); ok {
		// Child is definitely a leaf is a leaf, we're done
		if x.Operator == NOT {
			return NewOperation(AND, NewOperation(OR, x))
		}
		// CNFify child nodes so assumptions below can be made
		x = operandMap(x, CNF)
		switch x.Operator {
		case OR:
			and := NewOperation(AND)
			count := make([]int, len(x.Operands))
			maxcount := make([]int, len(x.Operands))
			minterms := x.Operands
			for i := range minterms {
				maxcount[i] = len(minterms[i].(*Operation).Operands)
			}
			for count[len(count)-1] < maxcount[len(count)-1] {
				or := NewOperation(OR)
				for i, idx := range count {
					maxterms := minterms[i].(*Operation).Operands[idx].(*Operation).Operands
					for i := range maxterms {
						or.PushOperands(maxterms[i])
					}
				}
				and.PushOperands(or)

				carry := 1
				for i := range count {
					count[i] += carry
					carry = 0
					if count[i] >= maxcount[i] && i != len(count)-1 {
						carry = 1
						count[i] = 0
					}
				}
			}
			return and
		case AND:
			and := NewOperation(AND)
			for i := range x.Operands {
				and.PushOperands(x.Operands[i].(*Operation).Operands...)
			}
			return and
		default:
			panic("Unexpected Operator type while CNFing")
		}
	}
	return NewOperation(AND, NewOperation(OR, n))
}
