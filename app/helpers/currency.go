package helpers

import (
	"github.com/epointpayment/mloc-cpe/app/log"

	"github.com/juju/errors"
	"github.com/rmg/iso4217"
	"github.com/shopspring/decimal"
)

var (
	// DefaultCurrency is the alphabetic code used to represent the default currency
	DefaultCurrency = "USD"

	// DefaultCurrencyPrecision is the number of digits to keep past the decimal point
	DefaultCurrencyPrecision int32

	// ErrInvalidCurrency is an error given when the default currency is set to an invalid value
	ErrInvalidCurrency = errors.New("Invalid default currency used")

	// ErrInvalidCurrencyAmount is an error given when the default currency is set to an invalid value
	ErrInvalidCurrencyAmount = errors.New("Invalid default currency amount used")
)

func init() {
	var currencyCode, currencyPrecision int

	currencyCode, currencyPrecision = iso4217.ByName(DefaultCurrency)
	if currencyCode == 0 {
		log.Fatal(ErrInvalidCurrency)
	}

	DefaultCurrencyPrecision = int32(currencyPrecision)
}

// ValidateCurrencyAmount is used by the validator to check if value is valid
func ValidateCurrencyAmount(value interface{}) (err error) {
	if dec, ok := value.(decimal.Decimal); ok {
		if isValidateCurrencyAmount(dec) {
			return nil
		}
	}

	return errors.Trace(ErrInvalidCurrencyAmount)
}

func isValidateCurrencyAmount(dec decimal.Decimal) bool {
	if dec.Truncate(DefaultCurrencyPrecision).Equal(dec) {
		// if dec.GreaterThan(decimal.Zero) {
		return true
		// }
	}

	return false
}
