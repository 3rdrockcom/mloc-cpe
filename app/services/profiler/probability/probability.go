package probability

import (
	"sort"

	"github.com/epointpayment/mloc-cpe/app/models"
	"github.com/epointpayment/mloc-cpe/app/services/profiler/probability/day"
	"github.com/epointpayment/mloc-cpe/app/services/profiler/probability/weekday"
)

// Precision is the default numerical precision for results
var Precision int32 = 3

// Probability is a service that calculates the probability
type Probability struct {
	Transactions models.Transactions
	Precision    int32
}

// New creates an instance of the probability service
func New(t models.Transactions) *Probability {
	// Sort transactions by date
	sort.Sort(t)

	return &Probability{
		Transactions: t,
		Precision:    Precision,
	}
}

// RunDay calculates the daily probability in a month
func (p *Probability) RunDay() day.Results {
	d := day.New(p.Transactions)
	d.Precision = p.Precision

	return d.Run()
}

// RunWeekday calculates the weekday probability in a week
func (p *Probability) RunWeekday() weekday.Results {
	w := weekday.New(p.Transactions)
	w.Precision = p.Precision

	return w.Run()
}
