package weekday

import (
	"github.com/epointpayment/mloc-cpe/app/models"

	"github.com/shopspring/decimal"
)

// DefaultPrecision is the default numerical precision for results
const DefaultPrecision int32 = 3

// Weekday contains information required to calculate probabilities
type Weekday struct {
	Transactions models.Transactions
	Precision    int32
}

// New creates an instance of the weekday probability service
func New(transactions models.Transactions) *Weekday {
	return &Weekday{
		Transactions: transactions,
		Precision:    DefaultPrecision,
	}
}

// Run executes the calculation
func (w *Weekday) Run() Results {
	return w.calc()
}

// calc calculates the probability for a particular weekday in a week
func (w *Weekday) calc() Results {
	list := make(Results, 7)

	// Aggregate totals for each day in a week
	count := 0
	total := decimal.Zero
	for i := range w.Transactions {
		weekday := int(w.Transactions[i].DateTime.Weekday())

		list[weekday].Weekday = w.Transactions[i].DateTime.Weekday()
		list[weekday].Total = list[weekday].Total.Add(w.Transactions[i].Credit)
		list[weekday].Count++

		total = total.Add(w.Transactions[i].Credit)
		count++
	}

	// Calculate probabilities
	for i := range list {
		list[i].Probability = list[i].Total.Div(total).Round(w.Precision)
	}

	return list
}
