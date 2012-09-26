package stoichio

type Vector []Cell

func (v Vector) Supp() []int {
	r := make([]int, 0, len(v))
	for i, val := range v {
		if val != 0 {
			r = append(r, i)
		}
	}
	return r
}
