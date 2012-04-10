package main

import (
	"stoichio/logic"
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	input = flag.String("i", "", "Input to read from.")
)

func main() {
	flag.Parse()

	matrixstring, irreversiblestring := ReadInputStrings()
	stoichio := ParseMatrix(matrixstring)
	irreversible := ParseIrreversible(irreversiblestring)
	checkSanity(stoichio, irreversible)
	l := generateLogic(stoichio, irreversible)
	cnf := logic.CNF(l)
	s := formatSAT(cnf)
	fmt.Printf("Logic: %s\nSAT:\n%s\n", l, s)
}

func checkSanity(stoichio StoichioMatrix, irreversible []bool) {
	nr := len(stoichio)
	if nr <= 0 {
		panic("Matrix with 0 rows")
	}
	nc := len(stoichio[0])
	if nc <= 0 {
		panic("Matrix with 0 columns")
	}
	for i := range stoichio {
		if len(stoichio[i]) != nc {
			panic("Not all rows have same length")
		}
	}
	if nc != len(irreversible) {
		panic("Size of irreversible-reactions-list doesn't match total number of reactions")
	}
}

func ReadInputStrings() (string, string) {
	var r *bufio.Reader

	if *input == "" || *input == "-" {
		r = bufio.NewReader(os.Stdin)
	} else {
		f, e := os.Open(*input)
		if e != nil {
			panic("Could not open input file")
		}
		defer f.Close()
		r = bufio.NewReader(f)
	}
	matrixstring, prefix, e := r.ReadLine()
	if e != nil || prefix {
		panic("Could not read matrix line")
	}
	irreversiblestring, prefix, e := r.ReadLine()
	if e != nil || prefix {
		panic("Could not read reaction line")
	}
	return string(matrixstring), string(irreversiblestring)
}

type Cell int8
type StoichioMatrix [][]Cell

func ParseMatrix(ms string) StoichioMatrix {
	ms = cleanString(ms)
	rowstrings := strings.Split(ms, ";")
	r := make([][]Cell, len(rowstrings))
	for i, rowstring := range rowstrings {
		cellstrings := strings.Fields(rowstring)
		r[i] = make([]Cell, len(cellstrings))
		for j, cellstring := range cellstrings {
			v, e := strconv.ParseInt(cellstring, 0, 64)
			if e != nil {
				panic("Invalid cell value: " + cellstring)
			}
			r[i][j] = Cell(v)
		}
	}
	return StoichioMatrix(r)
}

func cleanString(ms string) string {
	ms = regexp.MustCompile("\\];?$").ReplaceAllString(ms, "")
	ms = regexp.MustCompile("^\\[").ReplaceAllString(ms, "")
	return ms
}

func (sm StoichioMatrix) String() string {
	s := "["
	rowsep := ""
	for _, row := range sm {
		s += rowsep
		cellsep := ""
		for _, cell := range row {
			s += fmt.Sprintf("%s%d", cellsep, cell)
			cellsep = " "
		}
		rowsep = "; "
	}
	s += "]"
	return s
}

func ParseIrreversible(is string) []bool {
	is = cleanString(is)
	fields := strings.Fields(is)
	r := make([]bool, len(fields))
	for i, field := range fields {
		if field == "1" {
			r[i] = true
		}
	}
	return r
}

func generateLogic(stoichio StoichioMatrix, irreversible []bool) logic.Node {
	reversiblemap := map[string]string{}
	root := logic.NewOperation(logic.AND)
	for _, metabol := range stoichio {
		metaboliteins, metaboliteouts := logic.NewOperation(logic.OR), logic.NewOperation(logic.OR)
		// Traverse metabolites
		for reactionidx, reaction := range metabol {
			reactionname := strconv.Itoa(reactionidx+1)
			if !irreversible[reactionidx] {
				if _, ok := reversiblemap[reactionname+"f"]; !ok {
					reversiblemap[reactionname+"f"] = strconv.Itoa(len(irreversible) + 2*len(reversiblemap))
					reversiblemap[reactionname+"b"] = strconv.Itoa(len(irreversible) + 2*len(reversiblemap)+1)
				}
				if reaction > 0 {
					metaboliteins.PushOperands(logic.NewLeaf(reversiblemap[reactionname+"f"]))
					metaboliteouts.PushOperands(logic.NewLeaf(reversiblemap[reactionname+"b"]))
				} else if reaction < 0 {
					metaboliteins.PushOperands(logic.NewLeaf(reversiblemap[reactionname+"b"]))
					metaboliteouts.PushOperands(logic.NewLeaf(reversiblemap[reactionname+"f"]))
				}
			} else {
				if reaction > 0 {
					metaboliteins.PushOperands(logic.NewLeaf(reactionname))
				} else if reaction < 0 {
					metaboliteouts.PushOperands(logic.NewLeaf(reactionname))
				}
			}
		}
		root.PushOperands(logic.NewOperation(logic.IFF, metaboliteins, metaboliteouts))
	}

	// Traverse irreversible reactions
	for reactionidx, isIrreversible := range irreversible {
		if isIrreversible {
			continue
		}
		varname := strconv.Itoa(reactionidx+1)
		if _, ok := reversiblemap[varname+"f"]; !ok {
			reversiblemap[varname+"f"] = strconv.Itoa(len(irreversible) + 2*len(reversiblemap))
			reversiblemap[varname+"b"] = strconv.Itoa(len(irreversible) + 2*len(reversiblemap)+1)
		}
		in := logic.NewLeaf(reversiblemap[varname+"f"])
		out := logic.NewLeaf(reversiblemap[varname+"b"])
		reaction := logic.NewLeaf(varname)
		exclusion := logic.NewOperation(logic.NOT, logic.NewOperation(logic.AND, in, out))
		implication := logic.NewOperation(logic.IFF, reaction, logic.NewOperation(logic.OR, in, out))
		root.PushOperands(exclusion)
		root.PushOperands(implication)
	}
	return root
}


func formatSAT(n logic.Node) string {
	s := ""
	if _, ok := n.(logic.Leaf); ok {
		return n.String()
	}
	x := n.(*logic.Operation)
	switch x.Operator {
	case logic.AND:
		for _, op := range x.Operands {
			s += formatSAT(op)+" 0 \n"
		}
	case logic.OR:
		sep := ""
		for _, op := range x.Operands {
			s += sep+formatSAT(op)
			sep = " "
		}
	case logic.NOT:
		s = "-" + x.Operands[0].String()
	default:
		panic("This is not possibe o_O")
	}
	return s
}
