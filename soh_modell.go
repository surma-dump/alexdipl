package main

import (
	"./logic"
	"./stoichio"
	"fmt"
	"github.com/surma/goptions"
	"log"
)

const METABOL = "m_%d_t=%d"
const REACTION = "r_%d_t=%d"

const VERSION = "0.1"

func main() {
	options := struct {
		InputFile     string `goptions:"-i, --input, description='File to read', obligatory"`
		TimeLimit int `goptions:"-t, --time, description='Maximum number of timesteps (default: 10)'"`
		goptions.Help `goptions:"-h, --help, description='Show this help'"`
	} {
		TimeLimit: 10,
	}

	err := goptions.Parse(&options)
	if err != nil {
		log.Printf("Error: %s", err)
		goptions.PrintHelp()
		return
	}

	matrix, irreversible, err := stoichio.ReadFile(options.InputFile)
	if err != nil {
		log.Fatalf("Could not read file: %s", err)
	}

	t_in, t_out := make([]int, 0), make([]int, 0)
	sourceset := make([]int, 0)
	for i := 0; i < matrix.NumCols(); i++ {
		col := matrix.Col(i)
		if len(col.Supp()) != 1 {
			continue
		}
		supp := col.Supp()
		if col[supp[0]] > 0 || (col[supp[0]] < 0 && !irreversible[i]){
			t_in = append(t_in, i)
			sourceset = append(sourceset, supp[0])
		}
		if col[supp[0]] < 0 || (col[supp[0]] > 0 && !irreversible[i]){
			t_out = append(t_out, i)
		}
	}
	log.Printf("t_in: %#v", t_in)
	log.Printf("t_out: %#v", t_out)
	log.Printf("sorceset: %#v", sourceset)

	a1 := logic.NewOperation(logic.AND)
	for i := 0; i < matrix.NumRows(); i++ {
		if contains(sourceset, i) {
			a1.PushOperands(logic.NewLeaf(fmt.Sprintf(METABOL, i, 0)))
		} else {
			a1.PushOperands(logic.NewOperation(logic.NOT, logic.NewLeaf(fmt.Sprintf(METABOL, i, 0))))
		}
	}

	a2 := logic.NewOperation(logic.AND)
	a3 := logic.NewOperation(logic.AND)
	a4 := logic.NewOperation(logic.AND)
	a5 := logic.NewOperation(logic.AND)
	a6 := logic.NewOperation(logic.AND)
	a4.PushOperands(generateA4(0, matrix))
	a5.PushOperands(generateA5(0, matrix))
	for t := 1; t <= options.TimeLimit; t++ {
		a2.PushOperands(generateA2(t, matrix, irreversible))
		a3.PushOperands(generateA3(t, matrix, irreversible))
		a4.PushOperands(generateA4(t, matrix))
		a5.PushOperands(generateA5(t, matrix))
	}
	a6.PushOperands(generateA6(options.TimeLimit, matrix)...)
	log.Printf("A1: %s", a1)
	log.Printf("A2: %s", a2)
	log.Printf("A3: %s", a3)
	log.Printf("A4: %s", a4)
	log.Printf("A5: %s", a5)
	log.Printf("A6: %s", a6)
	a := logic.NewOperation(logic.AND, a1, a2, a3, a4, a5, a6)
	log.Printf("A: %s", a)
	log.Printf("CNF(A): %s", logic.CNF(a))
	log.Printf("SAT(A): %s", logic.FormatSAT(a))


}

func generateA2(t int, matrix stoichio.Matrix, irreversible []bool) logic.Node {
	m := logic.NewOperation(logic.AND)
	for j := 0; j < matrix.NumCols(); j++ {
		alpha := logic.NewOperation(logic.AND)
		for i := 0; i < matrix.NumRows(); i++ {
			if matrix[i][j] < 0 || (matrix[i][j] > 0 && !irreversible[j]) {
				alpha.PushOperands(logic.NewLeaf(fmt.Sprintf(METABOL, i, t-1)))
			}
		}
		if len(alpha.Operands) > 0 {
			m.PushOperands(logic.NewOperation(logic.IF,
				logic.NewLeaf(fmt.Sprintf(REACTION, j, t)),
				alpha))
		}
	}
	return m
}

func generateA3(t int, matrix stoichio.Matrix, irreversible []bool) logic.Node {
	m := logic.NewOperation(logic.AND)
	for j := 0; j < matrix.NumCols(); j++ {
		beta := logic.NewOperation(logic.AND)
		for i := 0; i < matrix.NumRows(); i++ {
			if matrix[i][j] > 0 || (matrix[i][j] < 0 && !irreversible[j]) {
				beta.PushOperands(logic.NewLeaf(fmt.Sprintf(METABOL, i, t)))
			}
		}
		if len(beta.Operands) > 0 {
			m.PushOperands(logic.NewOperation(logic.IF,
				logic.NewLeaf(fmt.Sprintf(REACTION, j, t)),
				beta))
		}
	}
	return m
}

func generateA4(t int, matrix stoichio.Matrix) logic.Node {
	m := logic.NewOperation(logic.AND)
	for i := 0; i < matrix.NumRows(); i++ {
		m.PushOperands(logic.NewOperation(logic.IF,
			logic.NewLeaf(fmt.Sprintf(METABOL, i, t)),
			logic.NewLeaf(fmt.Sprintf(METABOL, i, t+1))))
	}
	return m
}

func generateA5(t int, matrix stoichio.Matrix) logic.Node {
	m := logic.NewOperation(logic.AND)
	for j := 0; j < matrix.NumCols(); j++ {
		m.PushOperands(logic.NewOperation(logic.IF,
			logic.NewLeaf(fmt.Sprintf(REACTION, j, t)),
			logic.NewLeaf(fmt.Sprintf(REACTION, j, t+1))))
	}
	return m
}

func generateA6(t int, matrix stoichio.Matrix) []logic.Node {
	m := make([]logic.Node, 0)
	for i := 0; i < matrix.NumRows(); i++ {
		m = append(m, logic.NewLeaf(fmt.Sprintf(METABOL, i, t)))
	}
	return m
}

func contains(a []int, i int) bool {
	for _, v := range a {
		if v == i {
			return true
		}
	}
	return false
}

