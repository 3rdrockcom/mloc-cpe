package classifier

import (
	"time"

	"github.com/shopspring/decimal"
)

type Credits map[int][]Credit

type Credit struct {
	Date   time.Time
	Amount decimal.Decimal
}
