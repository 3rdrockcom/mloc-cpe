package profiler

import (
	"sort"

	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/models"
	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/services/profiler/classifier"
	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/services/profiler/probability"
)

var Debug bool

type Profiler struct {
	Transactions models.Transactions
	Partitions   int
	separator    float64
}

func New(t models.Transactions, pa int) *Profiler {
	p := &Profiler{
		Transactions: t,
		Partitions:   pa,
		separator:    1 / float64(pa),
	}

	return p
}

func (p *Profiler) Run() Results {
	res := make(Results, 0)

	tSplit := p.Transactions.Separator(p.separator)

	for a := 0; a < len(tSplit); a++ {
		transactions := models.Transactions{}
		for _, transaction := range tSplit[a] {
			if transaction.Credit > 0 {
				transactions = append(transactions, transaction)
			}
		}

		if len(transactions) == 0 {
			break
		}

		entries := make([]Result, 0)

		// Classify account
		cl, err := classifier.NewClassifier(transactions)
		if err != nil {
			panic(err)
		}
		clr := cl.Process()

		for i := range clr {
			entry := Result{}

			entry.Classification = clr[i].Name
			entry.Match = clr.GetProbability(i) * 100

			entry.Credits = Credits{
				AveragePerInterval: clr.GetAveragePerInterval(i),
				Average:            clr.GetAverage(i),
			}

			// Statistics
			entry.Statistics = Statistics{}

			pr := probability.New(transactions)

			probDay := pr.RunDay()
			sort.Sort(probDay)
			pd := make([]Day, 0)
			for g := range probDay {
				if probDay[g].Probability == 0 {
					break
				}
				pd = append(pd, Day{
					Day:         probDay[g].Day,
					Probability: probDay[g].Probability * 100,
				})
			}
			entry.Statistics.Day = pd

			if clr[i].Name == "weekly" {
				probWeekday := pr.RunWeekday()
				sort.Sort(probWeekday)
				pw := make([]Weekday, 0)
				for g := range probWeekday {
					if probWeekday[g].Probability == 0 {
						break
					}
					pw = append(pw, Weekday{
						Weekday:     probWeekday[g].Weekday.String(),
						Probability: probWeekday[g].Probability * 100,
					})
				}
				entry.Statistics.Weekday = pw
			}

			entries = append(entries, entry)
		}

		res = append(res, entries)
	}

	return res
}
