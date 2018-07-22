package day

import (
	"github.com/epointpayment/mloc-cpe/app/models"

	"github.com/shopspring/decimal"
)

// DefaultPrecision is the default numerical precision for results
const DefaultPrecision int32 = 3

// Day contains information required to calculate probabilities
type Day struct {
	Transactions models.Transactions
	Precision    int32
}

// New creates an instance of the day probability service
func New(transactions models.Transactions) *Day {
	return &Day{
		Transactions: transactions,
		Precision:    DefaultPrecision,
	}
}

// Run executes the calculation
func (d *Day) Run() Results {
	return d.calc()
}

// calc calculates the probability for a particular day in a month
func (d *Day) calc() Results {
	list := make(Results, 31)

	// Aggregate totals for each day in a month
	count := 0
	total := decimal.Zero
	for i := range d.Transactions {
		day := d.Transactions[i].DateTime.Day() - 1

		list[day].Day = d.Transactions[i].DateTime.Day()
		list[day].Total = list[day].Total.Add(d.Transactions[i].Credit)
		list[day].Count++

		total = total.Add(d.Transactions[i].Credit)
		count++
	}

	// Calculate probabilities
	for i := range list {
		list[i].Probability = list[i].Total.Div(total).Round(d.Precision)
	}

	return list
}
