package main

import (
	"./logic"
	"./stoichio"
	"fmt"
	"github.com/voxelbrain/goptions"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
)

const METABOL = "m_%d_t=%d"
const REACTION = "r_%d_t=%d"

const VERSION = "0.1"

func main() {
	options := struct {
		InputFile     string `goptions:"-i, --input, description='File to read', obligatory"`
		TimeLimit     int    `goptions:"-t, --time, description='Maximum number of timesteps (default: 10)'"`
		Targetset     string `goptions:"-z, --targetset, description='Comma-separated list of metabolite indices'"`
		Verbosity     []bool `goptions:"-v, --verbose, description='Increase verbosity'"`
		SAT           bool   `goptions:"-s, --output-sat, description='Output in SAT format instead of human-readable CNF'"`
		goptions.Help `goptions:"-h, --help, description='Show this help'"`
	}{
		TimeLimit: 10,
	}

	err := goptions.Parse(&options)
	if err != nil {
		log.Printf("Error: %s", err)
		goptions.PrintHelp()
		return
	}

	z := []string{}
	if len(options.Targetset) > 0 {
		z = strings.Split(options.Targetset, ",")
		for i, elem := range z {
			z[i] = strings.TrimSpace(elem)
		}
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
		if col[supp[0]] > 0 || (col[supp[0]] < 0 && !irreversible[i]) {
			t_in = append(t_in, i)
			sourceset = append(sourceset, supp[0])
		}
		if col[supp[0]] < 0 || (col[supp[0]] > 0 && !irreversible[i]) {
			t_out = append(t_out, i)
		}
	}
	if len(options.Verbosity) >= 2 {
		log.Printf("t_in: %#v", t_in)
		log.Printf("t_out: %#v", t_out)
		log.Printf("sourceset: %#v", sourceset)
		log.Printf("z: %#v", z)
	}

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
	a6.PushOperands(generateA6(options.TimeLimit, z)...)
	a := logic.NewOperation(logic.AND, a1, a2, a3, a4, a5)
	if len(z) > 0 {
		a.PushOperands(a6)
	}

	if options.SAT {
		sat, table := logic.FormatSAT(a)
		list := sortTable(table)
		fmt.Println("Table:")
		tw := tabwriter.NewWriter(os.Stdout, 4, 4, 1, ' ', tabwriter.AlignRight)
		for _, v := range list {
			fmt.Fprintf(tw, "\t%d\t=>\t%s\t\n", v.Id, v.Name)
		}
		tw.Flush()
		fmt.Println("")
		fmt.Println("SAT:")
		fmt.Println(sat)
	} else {
		if len(options.Verbosity) >= 1 {
			log.Printf("A1:\n%s", a1)
			log.Printf("A2:\n%s", a2)
			log.Printf("A3:\n%s", a3)
			log.Printf("A4:\n%s", a4)
			log.Printf("A5:\n%s", a5)
			log.Printf("A6:\n%s", a6)
			log.Printf("A:\n%s", a)
		}
		log.Printf("CNF(A):\n%s", logic.CNF(a))
	}

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

func generateA6(t int, targetset []string) []logic.Node {
	m := make([]logic.Node, 0)
	for _, idx := range targetset {
		i, e := strconv.ParseInt(idx, 10, 64)
		if e != nil {
			log.Fatalf("Invalid integer in target set: %s", idx)
		}
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

type ListEntry struct {
	Name string
	Id   int
}

type List []ListEntry

func (t List) Len() int {
	return len(t)
}

func (t List) Less(i, j int) bool {
	return t[i].Id < t[j].Id
}

func (t List) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func sortTable(hashtable map[string]int) List {
	list := make(List, 0, len(hashtable))
	for k, v := range hashtable {
		list = append(list, ListEntry{
			Name: k,
			Id:   v,
		})
	}
	sort.Sort(list)
	return list
}
