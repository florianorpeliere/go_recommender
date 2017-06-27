package neighborhood_model

import (
	"sort"
)

type Pair struct {
	Key int
	Value float64
}

type PairList []Pair
func (p PairList) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p PairList) Len() int { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }

func sortMapByValue(m map[int]float64) PairList {
	p := make(PairList, len(m))
	i := 0
	for k, v := range m {
		p[i] = Pair{Key: k, Value: v}
		i++
	}
	sort.Sort(sort.Reverse(p))
	return p
}