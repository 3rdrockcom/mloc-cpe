package profiler

import (
	"sort"

	"github.com/epointpayment/mloc-cpe/app/services/profiler/probability/day"
	"github.com/epointpayment/mloc-cpe/app/services/profiler/probability/weekday"
	"github.com/juju/errors"

	"github.com/epointpayment/mloc-cpe/app/models"
	"github.com/epointpayment/mloc-cpe/app/services/profiler/classifier"
	"github.com/epointpayment/mloc-cpe/app/services/profiler/probability"

	"github.com/shopspring/decimal"
)

// DefaultPartitions is the number of partitions to split transactions
const DefaultPartitions = 2

func init() {
	classifier.Precision = 5
	probability.Precision = 5
}

// Profiler is a service that profiles a customer's transactions to find an optimal interval
type Profiler struct {
	Transactions models.Transactions
	Partitions   int
	separator    float64
}

// New creates an instance of the profiler service
func New(t models.Transactions, pa int) *Profiler {
	p := &Profiler{
		Transactions: t,
		Partitions:   pa,
		separator:    1 / float64(pa),
	}

	return p
}

// Run processes results to obtain a list of ranked classifications
func (p *Profiler) Run() (res Results, err error) {
	// Split transactions into partitions
	tSplit := p.Transactions.Separator(p.separator)

	// Iterate through partitions
	for a := 0; a < len(tSplit); a++ {
		// Get list of all credit transactions
		transactions := models.Transactions{}
		for _, transaction := range tSplit[a] {
			if transaction.Credit.GreaterThan(decimal.Zero) {
				transactions = append(transactions, transaction)
			}
		}

		// No transactions to process
		if len(transactions) == 0 {
			break
		}

		classifications := make([]Result, 0)

		// Classify account
		cl, err := classifier.New(transactions)
		if err != nil {
			err = errors.Trace(err)
			return res, err
		}
		clr, err := cl.Process()
		if err != nil {
			err = errors.Trace(err)
			return res, err
		}

		for i := range clr {
			// Initialize service
			ps := probability.New(transactions)

			// Prepare classification
			classification := Result{
				Classification: clr[i].Name,
				Match:          clr.GetProbability(i).Mul(decimal.New(100, 0)),
				Credits: Credits{
					AveragePerInterval: clr.GetAveragePerInterval(i),
					Average:            clr.GetAverage(i),
				},
			}

			// Calculate day probabilities for a month
			classification.Statistics.Day = formatDay(ps.RunDay())

			// Calculate weekday probabilities for a week
			if clr[i].Name == "weekly" {
				classification.Statistics.Weekday = formatWeekday(ps.RunWeekday())
			}

			classifications = append(classifications, classification)
		}

		res = append(res, classifications)
	}

	return
}

// formatDay formats the probability results
func formatDay(r day.Results) (p []Day) {
	// Sort by descending rank (high to low)
	sort.Sort(r)

	// Get list of probabilities
	for i := range r {
		if r[i].Probability.Equal(decimal.Zero) {
			break
		}
		p = append(p, Day{
			Day:         r[i].Day,
			Probability: r[i].Probability.Mul(decimal.New(100, 0)),
		})
	}

	return
}

// formatWeekday formats the probability results
func formatWeekday(r weekday.Results) (p []Weekday) {
	// Sort by descending rank (high to low)
	sort.Sort(r)

	// Get list of probabilities
	for i := range r {
		if r[i].Probability.Equal(decimal.Zero) {
			break
		}
		p = append(p, Weekday{
			Weekday:     r[i].Weekday.String(),
			Probability: r[i].Probability.Mul(decimal.New(100, 0)),
		})
	}

	return
}
