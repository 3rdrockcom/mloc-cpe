package profiler

import (
	"github.com/shopspring/decimal"
)

type Results [][]Result

type Result struct {
	Classification string          `json:"classification"`
	Match          decimal.Decimal `json:"match"`
	Credits        Credits         `json:"credits"`
	Statistics     Statistics      `json:"statistics"`
}

type Credits struct {
	AveragePerInterval decimal.Decimal `json:"average_per_interval"`
	Average            decimal.Decimal `json:"average"`
}

type Statistics struct {
	Day     []Day     `json:"day"`
	Weekday []Weekday `json:"weekday,omitempty"`
}

type Day struct {
	Day         int             `json:"day"`
	Probability decimal.Decimal `json:"probability"`
}

type Weekday struct {
	Weekday     string          `json:"weekday"`
	Probability decimal.Decimal `json:"probability"`
}
