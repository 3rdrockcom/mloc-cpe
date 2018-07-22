package profiler

import (
	"github.com/shopspring/decimal"
)

// Results containts a list of results
type Results [][]Result

// Result contains information about a classification
type Result struct {
	Classification string          `json:"classification"`
	Match          decimal.Decimal `json:"match"`
	Credits        Credits         `json:"credits"`
	Statistics     Statistics      `json:"statistics"`
}

// Credits contains information about transactions
type Credits struct {
	AveragePerInterval decimal.Decimal `json:"average_per_interval"`
	Average            decimal.Decimal `json:"average"`
}

// Statistics contains statistics for a month and week (for weekly classifications only)
type Statistics struct {
	Day     []Day     `json:"day"`
	Weekday []Weekday `json:"weekday,omitempty"`
}

// Day contains information about day's probability for a month
type Day struct {
	Day         int             `json:"day"`
	Probability decimal.Decimal `json:"probability"`
}

// Weekday contains information about weekday's probability for a week
type Weekday struct {
	Weekday     string          `json:"weekday"`
	Probability decimal.Decimal `json:"probability"`
}
