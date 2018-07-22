package classifier

import (
	"sort"

	"github.com/shopspring/decimal"
)

// Results is a list of classifier results
type Results []Result

// Result contains information about classifier entry
type Result struct {
	Name        string
	Score       float64
	Probability decimal.Decimal
	List        Credits
}

// GetClassification gets a classification result from the result set
func (r Results) GetClassification(id int) Result {
	classification := r[id]
	return classification
}

// GetProbability calculates the probability for a result entry
func (r Results) GetProbability(id int) decimal.Decimal {
	data := make([]float64, len(r))

	// Normalize scores, zero being the lowest value
	minScore := 0.0
	for i := range r {
		score := r[i].Score

		if score < minScore {
			minScore = score
		}

		data[i] = score
	}

	// Get the sum total of all the scores
	sum := 0.0
	for i := range data {
		data[i] = data[i] - minScore
		sum += data[i]
	}

	// Calculate the probability
	for i := range data {
		// Prevent division by zero error
		if sum == 0 {
			data[i] = 0
			continue
		}
		data[i] = (data[i] / sum)
	}

	return decimal.NewFromFloat(data[id]).Round(Precision)
}

// GetAveragePerInterval calculates the average result amount per interval
func (r Results) GetAveragePerInterval(id int) decimal.Decimal {
	data := []decimal.Decimal{}

	classification := r[id]
	list := classification.List

	// Get list of intervals sorted by time
	var keys []int
	for k := range list {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	// Interate through each interval
	for k := range keys {
		i := keys[k]
		sum := decimal.Zero

		// Sum all the amounts in a particular interval
		for j := range list[i] {
			sum = sum.Add(list[i][j].Amount)
		}

		data = append(data, sum)
	}

	// Calculate average
	avg, _ := calcAvg(data)
	return avg
}

// GetAverage calculates the average result amount
func (r Results) GetAverage(id int) decimal.Decimal {
	data := []decimal.Decimal{}

	classification := r[id]
	list := classification.List

	// Create a list of all the amounts
	for i := range list {
		for j := range list[i] {
			data = append(data, list[i][j].Amount)
		}
	}

	// Calculate average
	avg, _ := calcAvg(data)
	return avg
}
