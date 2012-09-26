package logic

import (
	"fmt"
	"strings"
)

func Parse(s string) (Node, error) {
	panic("not implemented")
	var n Node
	s = strings.TrimSpace(s)
	for opstring := range opFuncMap {
		if strings.HasPrefix(s, opstring) {
			n = &Operation{
				Operator: opstring,
			}
			s = s[len(opstring):]
			break
		}
	}
	if n == nil {
		return nil, fmt.Errorf("Unknown operand at start of %s...", s[0:10])
	}
	if strings.HasPrefix(s, "(") {
		return nil, fmt.Errorf("")
	}
	return nil, nil
}
