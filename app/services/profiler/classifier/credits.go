package classifier

import (
	"time"

	"github.com/shopspring/decimal"
)

// Credits contains a list of credits by id
type Credits map[int][]Credit

// Credit contains information about a credit transaction
type Credit struct {
	Date   time.Time
	Amount decimal.Decimal
}
