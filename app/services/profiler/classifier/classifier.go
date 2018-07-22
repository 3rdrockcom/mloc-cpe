package classifier

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/epointpayment/mloc-cpe/app/models"

	"github.com/jinzhu/now"
	"github.com/juju/errors"
	"github.com/montanaflynn/stats"
	"github.com/shopspring/decimal"
)

// Debug displays additional information
var Debug bool

// Precision is the default numerical precision for results
var Precision int32 = 3

// Classifier is a service that classifies a list of transactions
type Classifier struct {
	Transactions models.Transactions
}

// New creates an instance of the classifier service
func New(t models.Transactions) (*Classifier, error) {
	if len(t) == 0 {
		return nil, errors.New("transactions required")
	}

	// Sort transactions by date
	sort.Sort(t)

	c := &Classifier{
		Transactions: t,
	}
	return c, nil
}

// Process processes transactions to create a list of classifications
func (c *Classifier) Process() (res Results, err error) {
	var listRank Ranks
	var rank Rank
	var list = make(map[string]Credits)

	// Process transactions and append results to list
	buckets := []string{"monthly", "bimonthly", "weekly"}
	for i := range buckets {
		name := buckets[i]
		if Debug {
			fmt.Println(fmt.Sprintf("::: Class: %s :::\n", name))
		}

		switch name {
		case "monthly":
			list[name] = c.processMonthly()
			rank = NewRank(name, c.calcRankValue(list[name]), 10)
		case "bimonthly":
			list[name] = c.processBiMonthly()
			rank = NewRank(name, c.calcRankValue(list[name]), 20)
		case "weekly":
			list[name] = c.processWeekly()
			rank = NewRank(name, c.calcRankValue(list[name]), 30)
		}

		if Debug {
			fmt.Println(fmt.Sprintf("Score: %.6f\n", rank.Value))
		}
		listRank = append(listRank, rank)
	}

	// Sort rank in descending order (high to low)
	sort.Sort(sort.Reverse(listRank))

	// Prepare results
	for i := range listRank {
		entry := Result{
			Name:  listRank[i].Name,
			Score: listRank[i].Value,
			List:  list[listRank[i].Name],
		}
		res = append(res, entry)
	}

	return
}

// processMonthly is a classification that processes transactions for a monthly interval
func (c *Classifier) processMonthly() Credits {
	t := c.Transactions

	// Determine the date range
	dateMin, dateMax := c.getDateRange()
	dateRangeMin := now.New(dateMin).BeginningOfMonth()
	dateRangeMax := now.New(dateMax).EndOfMonth()

	list := make(Credits)

	// Aggregate transactions by interval
	for d := dateRangeMin; d.Before(dateRangeMax); d = d.AddDate(0, 1, 0) {
		k, _ := strconv.Atoi(d.Format("20060102"))
		list[k] = []Credit{}

		for i := 0; i < len(t); i++ {
			if (t[i].DateTime.After(d) || t[i].DateTime.Equal(d)) && t[i].DateTime.Before(d.AddDate(0, 1, 0)) {
				list[k] = append(list[k], Credit{
					Amount: c.Transactions[i].Credit,
					Date:   c.Transactions[i].DateTime,
				})
			}
		}
	}

	return list
}

// processBiMonthly is a classification that processes transactions for a bimonthly interval
func (c *Classifier) processBiMonthly() Credits {
	t := c.Transactions

	// Determine the date range
	dateMin, dateMax := c.getDateRange()
	dateRangeMin := now.New(dateMin).BeginningOfWeek()
	dateRangeMax := now.New(dateMax).EndOfWeek()

	list := make(Credits)

	// Aggregate transactions by interval
	for d := dateRangeMin; d.Before(dateRangeMax); d = d.AddDate(0, 0, 15) {
		k, _ := strconv.Atoi(d.Format("20060102"))
		list[k] = []Credit{}

		for i := 0; i < len(t); i++ {
			if (t[i].DateTime.After(d) || t[i].DateTime.Equal(d)) && t[i].DateTime.Before(d.AddDate(0, 0, 15)) {
				list[k] = append(list[k], Credit{
					Amount: c.Transactions[i].Credit,
					Date:   c.Transactions[i].DateTime,
				})
			}
		}
	}

	return list
}

// processWeekly is a classification that processes transactions for a weekly interval
func (c *Classifier) processWeekly() Credits {
	t := c.Transactions

	// Determine the date range
	dateMin, dateMax := c.getDateRange()
	dateRangeMin := now.New(dateMin).BeginningOfWeek()
	dateRangeMax := now.New(dateMax).EndOfWeek()

	list := make(Credits)

	// Aggregate transactions by interval
	for d := dateRangeMin; d.Before(dateRangeMax); d = d.AddDate(0, 0, 7) {
		k, _ := strconv.Atoi(d.Format("20060102"))
		list[k] = []Credit{}

		for i := 0; i < len(c.Transactions); i++ {
			if (t[i].DateTime.After(d) || t[i].DateTime.Equal(d)) && t[i].DateTime.Before(d.AddDate(0, 0, 7)) {
				list[k] = append(list[k], Credit{
					Amount: c.Transactions[i].Credit,
					Date:   c.Transactions[i].DateTime,
				})
			}
		}
	}

	return list
}

// calcRankValue calculates the score for a classification
func (c *Classifier) calcRankValue(list Credits) float64 {
	data := []decimal.Decimal{}

	var keys []int
	for k := range list {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	if Debug {
		fmt.Println("Date Range(s):")
	}

	rankValue := 0.0
	total := decimal.Zero
	for k := range keys {
		i := keys[k]

		sum := decimal.Zero
		for j := range list[i] {
			total = total.Add(list[i][j].Amount)
			sum = sum.Add(list[i][j].Amount)
		}

		data = append(data, sum)

		v := -1.0
		if len(list[i]) == 1 {
			v = 1
		}
		if len(list[i]) > 1 {
			v = -.5 * (float64(len(list[i])) - 1)
		}
		rankValue += v

		if Debug {
			fmt.Println(fmt.Sprintf("%v: %10v %5v [%v / %v]", i, sum, v, rankValue, k+1))
		}
	}

	rankValue = rankValue / float64(len(list))

	if Debug {
		mean, sd, _ := c.getStatistics(data)

		fmt.Println()
		fmt.Println(fmt.Sprintf("Statistics: %.2f ± %.2f", mean, sd))
	}

	return rankValue
}

// getDateRange gets the date range for a list of transactions
func (c *Classifier) getDateRange() (time.Time, time.Time) {
	dateMin := c.Transactions[0].DateTime
	dateMax := c.Transactions[len(c.Transactions)-1].DateTime

	return dateMin, dateMax
}

// getStatistics calculate the mean and standard deviation from a list
func (c *Classifier) getStatistics(input []decimal.Decimal) (float64, float64, error) {
	var err error

	// Create a list of values
	data := []float64{}
	for i := range input {
		f, _ := input[i].Float64()
		data = append(data, f)
	}

	// Calculate mean
	mean, err := stats.Mean(data)
	if err != nil {
		err = errors.Trace(err)
		return 0, 0, err
	}

	// Calculate standard deviation
	sd, err := stats.StandardDeviation(data)
	if err != nil {
		err = errors.Trace(err)
		return 0, 0, err
	}

	return mean, sd, nil
}
