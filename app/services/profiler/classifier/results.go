package classifier

import (
	"sort"

	"github.com/shopspring/decimal"
)

type Results []Result

type Result struct {
	Name        string
	Score       float64
	Probability decimal.Decimal
	List        Credits
}

func (r Results) GetClassification(e int) Result {
	classification := r[e]
	return classification
}

func (r Results) GetProbability(e int) decimal.Decimal {
	data := make([]float64, len(r))

	minScore := 0.0
	for i := range r {
		score := r[i].Score

		if score < minScore {
			minScore = score
		}

		data[i] = score
	}

	sum := 0.0
	for i := range data {
		data[i] = data[i] - minScore
		sum += data[i]
	}

	for i := range data {
		data[i] = (data[i] / sum)
	}

	return decimal.NewFromFloat(data[e]).Round(3)
}

func (r Results) GetAveragePerInterval(e int) decimal.Decimal {
	data := []decimal.Decimal{}

	classification := r[e]
	list := classification.List

	var keys []int
	for k := range list {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for k := range keys {
		i := keys[k]
		sum := decimal.Zero

		for j := range list[i] {
			sum = sum.Add(list[i][j].Amount)
		}

		data = append(data, sum)
	}

	mean, _ := calcAvg(data)
	return mean
}

func (r Results) GetAverage(e int) decimal.Decimal {
	data := []decimal.Decimal{}

	classification := r[e]
	list := classification.List

	for i := range list {
		for j := range list[i] {
			data = append(data, list[i][j].Amount)
		}
	}

	mean, _ := calcAvg(data)
	return mean
}
