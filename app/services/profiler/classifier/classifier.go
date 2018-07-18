package classifier

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/epointpayment/mloc-cpe/app/models"
	"github.com/epointpayment/mloc-cpe/app/services/profiler/ranks"

	"github.com/jinzhu/now"
	"github.com/montanaflynn/stats"
	"github.com/shopspring/decimal"
)

var Debug bool

type Classifier struct {
	Transactions models.Transactions
}

func NewClassifier(t models.Transactions) (*Classifier, error) {
	if len(t) == 0 {
		return nil, errors.New("transactions required")
	}

	sort.Sort(t)

	c := &Classifier{
		Transactions: t,
	}
	return c, nil
}

func (c *Classifier) Process() Results {
	var listRank ranks.Ranks
	var rank ranks.Rank
	var list = make(map[string]Credits)

	buckets := []string{"monthly", "bimonthly", "weekly"}
	for i := range buckets {
		name := buckets[i]
		if Debug {
			fmt.Println(fmt.Sprintf("::: Class: %s :::\n", name))
		}

		switch name {
		case "monthly":
			list[name] = c.processMonthly()
			rank = ranks.NewRank(name, c.calcRankValue(list[name]), 10)
		case "bimonthly":
			list[name] = c.processBiMonthly()
			rank = ranks.NewRank(name, c.calcRankValue(list[name]), 20)
		case "weekly":
			list[name] = c.processWeekly()
			rank = ranks.NewRank(name, c.calcRankValue(list[name]), 30)
		}

		if Debug {
			fmt.Println(fmt.Sprintf("Score: %.6f\n", rank.Value))
		}
		listRank = append(listRank, rank)
	}

	sort.Sort(sort.Reverse(listRank))

	res := Results{}
	for i := range listRank {
		entry := Result{
			Name:  listRank[i].Name,
			Score: listRank[i].Value,
			List:  list[listRank[i].Name],
		}
		res = append(res, entry)
	}

	return res
}

func (c *Classifier) processMonthly() Credits {
	t := c.Transactions

	dateMin, dateMax := c.getDateRange()
	dateRangeMin := now.New(dateMin).BeginningOfMonth()
	dateRangeMax := now.New(dateMax).EndOfMonth()

	list := make(Credits)

	//
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

func (c *Classifier) processBiMonthly() Credits {
	t := c.Transactions

	dateMin, dateMax := c.getDateRange()
	dateRangeMin := now.New(dateMin).BeginningOfWeek()
	dateRangeMax := now.New(dateMax).EndOfWeek()

	list := make(Credits)

	//
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

func (c *Classifier) processWeekly() Credits {
	t := c.Transactions

	dateMin, dateMax := c.getDateRange()
	dateRangeMin := now.New(dateMin).BeginningOfWeek()
	dateRangeMax := now.New(dateMax).EndOfWeek()

	list := make(Credits)

	//
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

func (c *Classifier) getDateRange() (time.Time, time.Time) {
	dateMin := c.Transactions[0].DateTime
	dateMax := c.Transactions[len(c.Transactions)-1].DateTime

	return dateMin, dateMax
}

func (c *Classifier) getStatistics(input []decimal.Decimal) (float64, float64, error) {
	var err error

	data := []float64{}
	for i := range input {
		f, _ := input[i].Float64()
		data = append(data, f)
	}

	mean, err := stats.Mean(data)
	if err != nil {
		return 0, 0, err
	}

	sd, err := stats.StandardDeviation(data)
	if err != nil {
		return 0, 0, err
	}

	return mean, sd, nil
}
