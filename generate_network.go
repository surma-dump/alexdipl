package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/voxelbrain/goptions"
)

var (
	options = struct {
		NumMetabolites int `goptions:"-m, --num-metabolites, obligatory, description='Number of metabolites in the generated network'"`
		NumReactions   int `goptions:"-r, --num-reactions, obligatory, description='Number of reactions in the generated network'"`
		HumanReadable bool `goptions:"-u, --human, description='Output matrix with newlines'"`
		Seed int64 `goptions:"-s, --seed, description='Seed for the random generator'"`
		goptions.Help  `goptions:"-h, --help, description='Show this help'"`
	}{
		Seed: time.Now().UnixNano(),
	}
)

func init() {
	goptions.ParseAndFail(&options)
	rand.Seed(int64(options.Seed))
}

type Matrix struct {
	data          []int
	Width, Height int
}

func main() {
	mx := NewMatrix(options.NumReactions, options.NumMetabolites)
	for m := 0; m < options.NumMetabolites; m++ {
		perm := rand.Perm(options.NumReactions)
		subperm := perm[0:rand.Intn(len(perm)-2)+2]
		cut := rand.Intn(len(subperm)-1)+1
		source, dest := subperm[0:cut], subperm[cut:]
		for _, r := range source {
			mx.Set(r, m, -1)
		}
		for _, r := range dest {
			mx.Set(r, m, 1)
		}
	}

	fmt.Printf("%s\n", mx)
	fmt.Printf("[")
	for r := 0; r < mx.Width; r++ {
		fmt.Printf("1 ")
	}
	fmt.Printf("]\n")
}

func NewMatrix(width, height int) *Matrix {
	return &Matrix{
		data:   make([]int, width*height),
		Width:  width,
		Height: height,
	}
}

func (m *Matrix) Set(x, y int, v int) {
	m.data[y*m.Width+x] = v
}

func (m *Matrix) Get(x, y int) int {
	return m.data[y*m.Width+x]
}

func (mx *Matrix) String() string {
	s := "[ "
	if options.HumanReadable{
		s += "\n"
	}
	sep := ""
	for m := 0; m < mx.Height; m++ {
		s += sep
		for r := 0; r < mx.Width; r++ {
			s += fmt.Sprintf("%+d ", mx.Get(r, m))
		}
		if options.HumanReadable{
			sep = "\n"
		} else {
			sep = "; "
		}
	}
	return s + "]"
}
