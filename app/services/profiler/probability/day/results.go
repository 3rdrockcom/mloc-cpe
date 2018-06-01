package day

import "github.com/shopspring/decimal"

type Results []Result

func (r Results) Len() int           { return len(r) }
func (r Results) Less(i, j int) bool { return r[i].Probability.GreaterThan(r[j].Probability) }
func (r Results) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }

type Result struct {
	Day         int             `json:"day"`
	Count       int             `json:"count"`
	Total       decimal.Decimal `json:"total"`
	Probability decimal.Decimal `json:"probability"`
}
