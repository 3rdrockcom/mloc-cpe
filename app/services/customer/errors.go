package customer

import (
	"net/http"

	"github.com/epointpayment/mloc-cpe/app/codes"
)

var (
	// ErrInvalidUniqueCustomerID is an error shown when customer unique ID is not a valid
	ErrInvalidUniqueCustomerID = codes.New("ERROR_CUSTOMER_UNIQUE_ID_INVALID").
					WithMessage("Invalid Unique Customer ID").
					WithStatusCode(http.StatusForbidden).
					InGroup("CUSTOMER").
					RegisterError()

	// ErrCustomerNotFound is an error for a non-existent customer
	ErrCustomerNotFound = codes.New("ERROR_CUSTOMER_NOT_FOUND").
				WithMessage("Customer not found").
				WithStatusCode(http.StatusNotFound).
				InGroup("CUSTOMER").
				RegisterError()

	// ErrInvalidData is an error for invalid payload data
	ErrInvalidData = codes.New("ERROR_CUSTOMER_INVALID_DATA").
			WithMessage("Payload contains invalid data").
			WithStatusCode(http.StatusBadRequest).
			InGroup("CUSTOMER").
			RegisterError()

	// ErrCreditDebitRequired is an error given when credit/debit is not set for a transaction
	ErrCreditDebitRequired = codes.New("ERROR_CUSTOMER_CREDIT_DEBIT_REQUIRED").
				WithMessage("Credit or debit is required").
				WithStatusCode(http.StatusBadRequest).
				InGroup("CUSTOMER").
				RegisterError()

	// ErrCreditNonPositiveValue is an error given when a transaction credit amount is not positive
	ErrCreditNonPositiveValue = codes.New("ERROR_CUSTOMER_CREDIT_NONPOSITIVE_VALUE").
					WithMessage("Credit must be greater than 0").
					WithStatusCode(http.StatusBadRequest).
					InGroup("CUSTOMER").
					RegisterError()

	// ErrInvalidTimestamp is an error for invalid transaction timestamps
	ErrInvalidTimestamp = codes.New("ERROR_CUSTOMER_TIMESTAMP_INVALID").
				WithMessage("Invalid timestamp").
				WithStatusCode(http.StatusBadRequest).
				InGroup("CUSTOMER").
				RegisterError()
)
