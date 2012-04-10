package logic

const (
	NOT = iota
	AND
	OR
	IFF
)

type Node interface {
	Eval(Configuration) bool
	String() string
}

type Operation struct {
	Operator OperatorFunc
	OpSymbol string
	Operands []Node
}

func DefaultMap(n Node) map[string]bool {
	r := make(map[string]bool)
	queue := make([]Node, 1, 20)
	queue[0] = n
	for len(queue) != 0 {
		cur := queue[len(queue)-1]
		queue = queue[0:len(queue)-1]

		switch x := cur.(type) {
		case Leaf:
			r[string(x)] = false
		case *Operation:
			queue = append(queue, x.Operands...)
		}
	}
	return r
}

func NewOperation(opidx int) *Operation {
	op := &Operation {
		Operands: make([]Node, 0, 10),
	}
	switch opidx {
	case NOT:
		op.Operator = not
		op.OpSymbol = "!"
	case AND:
		op.Operator = and
		op.OpSymbol = "^"
	case OR:
		op.Operator = or
		op.OpSymbol = "v"
	case IFF:
		op.Operator = not
		op.OpSymbol = "<=>"
	}
	return op
}

func (o *Operation) Eval(config Configuration) bool {
	return o.Operator(o.Operands, config)
}

func (o *Operation) PushOperand(n Node) {
	o.Operands = append(o.Operands, n)
}

func (o *Operation) String() string {
	r := "("
	ops := ""
	for _, operand := range o.Operands {
		r += ops+operand.String()
		ops = " " + o.OpSymbol + " "
	}
	return r+")"
}

type OperatorFunc func([]Node, Configuration) bool

type Leaf string

func (l Leaf) Eval(config Configuration) bool {
	return config[string(l)]
}

func (l Leaf) String() string {
	return string(l)
}

func NewLeaf(name string) Node {
	return Leaf(name)
}

type Configuration map[string]bool

func not(operands []Node, config Configuration) bool {
	if len(operands) != 1 {
		panic("`not` only takes 1 argument")
	}
	return !operands[0].Eval(config)
}

func NewNot(op Node) Node {
	return &Operation {
		Operator: not,
		OpSymbol: "!",
		Operands: []Node {op},
	}
}

func or(operands []Node, config Configuration) bool {
	if len(operands) == 0 {
		panic("Zero arguments for `or`")
	}
	for _, operand := range operands {
		if operand.Eval(config) {
			return true
		}
	}
	return false
}

func NewOr(op... Node) Node {
	return &Operation {
		Operator: or,
		OpSymbol: "v",
		Operands: op,
	}
}

func and(operands []Node, config Configuration) bool {
	if len(operands) == 0 {
		panic("Zero arguments for `and`")
	}
	for _, operand := range operands {
		if !operand.Eval(config) {
			return false
		}
	}
	return true
}

func NewAnd(op... Node) Node {
	return &Operation {
		Operator: and,
		OpSymbol: "^",
		Operands: op,
	}
}

func iff(operands []Node, config Configuration) bool {
	if len(operands) == 0 {
		panic("Zero arguments for `iff`")
	}
	r := true;
	for _, operand := range operands {
		op := operand.Eval(config)
		r = r == op
	}
	return r
}

func NewIff(op... Node) Node {
	return &Operation {
		Operator: iff,
		OpSymbol: "<=>",
		Operands: op,
	}
}

