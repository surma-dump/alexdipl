package logic

const (
	NOT = "!"
	AND = "^"
	OR  = "v"
	IFF = "<=>"
)

var (
	opFuncMap = map[string]opFunc{
		NOT: not,
		AND: and,
		OR:  or,
		IFF: iff,
	}
)

type Node interface {
	Eval(Configuration) bool
	String() string
}

type Operation struct {
	Operator string
	Operands []Node
}

func DefaultMap(n Node) map[string]bool {
	r := make(map[string]bool)
	queue := make([]Node, 1, 20)
	queue[0] = n
	for len(queue) != 0 {
		cur := queue[len(queue)-1]
		queue = queue[0 : len(queue)-1]

		switch x := cur.(type) {
		case Leaf:
			r[string(x)] = false
		case *Operation:
			queue = append(queue, x.Operands...)
		}
	}
	return r
}

func NewOperation(opidx string, ops ...Node) *Operation {
	op := &Operation{
		Operands: ops,
		Operator: opidx,
	}
	return op
}

func (o *Operation) Eval(config Configuration) bool {
	return opFuncMap[o.Operator](o.Operands, config)
}

func (o *Operation) PushOperands(n ...Node) {
	o.Operands = append(o.Operands, n...)
}

func (o *Operation) String() string {
	r := o.Operator + "("
	ops := ""
	for _, operand := range o.Operands {
		r += ops + operand.String()
		ops = ", "
	}
	return r + ")"
}

type opFunc func([]Node, Configuration) bool

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

func iff(operands []Node, config Configuration) bool {
	if len(operands) == 0 {
		panic("Zero arguments for `iff`")
	}
	r := true
	for _, operand := range operands {
		op := operand.Eval(config)
		r = r == op
	}
	return r
}
