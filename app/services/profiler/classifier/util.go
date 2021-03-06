package classifier

import (
	"github.com/juju/errors"
	"github.com/shopspring/decimal"
)

// calcAvg calculates the average from an array of decimal values
func calcAvg(input []decimal.Decimal) (avg decimal.Decimal, err error) {
	avg = decimal.Zero

	// Empty list, cannot calculate average
	if len(input) == 0 {
		err = errors.New("Input must not be empty")
		return
	}

	// No need to calculate average
	if len(input) == 1 {
		avg = input[0]
		return
	}

	// Get the average
	avg = decimal.Avg(input[0], input[1:]...)
	return avg, nil
}
