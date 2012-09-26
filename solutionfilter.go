package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	helpFlag = flag.Bool("help", false, "Show this help")
)

func main() {
	flag.Parse()

	if *helpFlag || flag.NArg() != 1 {
		fmt.Printf("Usage: solutionfilter <file>\n")
		flag.PrintDefaults()
		return
	}

	file := flag.Arg(0)
	solutions := parseSolutions(file)
	filter(solutions)
	for _, s := range solutions {
		if s == nil {
			continue
		}
		fmt.Printf("%s\n", s)
	}
}

func filter(sol []Solution) {
	for i := range sol {
		if sol[i].IsEmpty() {
			sol[i] = nil
		}
	}
	for i := range sol {
		if sol[i] == nil {
			continue
		}
		for j := range sol {
			if sol[j] == nil {
				continue
			}
			if i != j && sol[i].Contains(sol[j]) {
				sol[j] = nil
			}
		}
	}
}

type Solution []bool

func parseSolutions(file string) []Solution {
	f, e := os.Open(file)
	if e != nil {
		panic(e)
	}
	defer f.Close()
	r := bufio.NewReaderSize(f, 10000)
	solutions := make([]Solution, 0, 10)
	for line, _, e := r.ReadLine(); e == nil; line, _, e = r.ReadLine() {
		vars := strings.Fields(string(line))
		solution := make([]bool, len(vars)-1)
		for i, v := range vars {
			if v[0] != '-' && v != "0" {
				solution[i] = true
			}
		}
		solutions = append(solutions, solution)
	}
	return solutions
}

func (sol Solution) String() string {
	s := ""
	for i, b := range sol {
		if !b {
			s += "-"
		}
		s += fmt.Sprintf("%d ", i+1)
	}
	return s
}

func (sol Solution) Contains(sol2 Solution) bool {
	if sol2 == nil || sol == nil {
		return false
	}
	for i := range sol {
		if sol[i] && !sol2[i] {
			return false
		}
	}
	return true
}

func (sol Solution) IsEmpty() bool {
	for i := range sol {
		if sol[i] {
			return false
		}
	}
	return true
}
