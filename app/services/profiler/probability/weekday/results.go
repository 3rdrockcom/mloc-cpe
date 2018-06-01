package weekday

import (
	"time"

	"github.com/shopspring/decimal"
)

type Results []Result

func (r Results) Len() int           { return len(r) }
func (r Results) Less(i, j int) bool { return r[i].Probability.GreaterThan(r[j].Probability) }
func (r Results) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }

type Result struct {
	Weekday     time.Weekday    `json:"weekday"`
	Count       int             `json:"count"`
	Total       decimal.Decimal `json:"total"`
	Probability decimal.Decimal `json:"probability"`
}
