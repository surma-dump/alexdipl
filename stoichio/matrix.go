package stoichio

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Cell int8
type Matrix [][]Cell

func cleanString(ms string) string {
	ms = regexp.MustCompile("\\];?$").ReplaceAllString(ms, "")
	ms = regexp.MustCompile("^\\[").ReplaceAllString(ms, "")
	return ms
}

func ParseMatrix(ms string) (Matrix, error) {
	ms = cleanString(ms)
	rowstrings := strings.Split(ms, ";")
	r := make([][]Cell, len(rowstrings))
	for i, rowstring := range rowstrings {
		cellstrings := strings.Fields(rowstring)
		r[i] = make([]Cell, len(cellstrings))
		for j, cellstring := range cellstrings {
			v, e := strconv.ParseInt(cellstring, 0, 64)
			if e != nil {
				return nil, e
			}
			r[i][j] = Cell(v)
		}
	}
	return Matrix(r), nil
}

func (m Matrix) String() string {
	s := "["
	rowsep := ""
	for _, row := range m {
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

func ReadFile(file string) (Matrix, []bool, error) {
	f, e := os.Open(file)
	if e != nil {
		return nil, nil, e
	}
	defer f.Close()
	r := bufio.NewReader(f)
	matrixstring, prefix, e := r.ReadLine()
	if e != nil || prefix {
		return nil, nil, fmt.Errorf("Could not read matrix: %s", e)
	}
	irreversiblestring, prefix, e := r.ReadLine()
	if e != nil || prefix {
		return nil, nil, fmt.Errorf("Could not read reactions: %s", e)
	}
	matrix, err := ParseMatrix(string(matrixstring))
	if err != nil {
		return nil, nil, err
	}
	return matrix, ParseIrreversible(string(irreversiblestring)), nil
}
