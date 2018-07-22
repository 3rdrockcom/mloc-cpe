package day

import "github.com/shopspring/decimal"

// Results contains a list of result probabilities
type Results []Result

// Sort results by descending probability (high to low)
func (r Results) Len() int           { return len(r) }
func (r Results) Less(i, j int) bool { return r[i].Probability.GreaterThan(r[j].Probability) }
func (r Results) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }

// Result contains information about the probability for a particular day
type Result struct {
	Day         int
	Count       int
	Total       decimal.Decimal
	Probability decimal.Decimal
}
