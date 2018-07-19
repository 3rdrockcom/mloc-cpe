package api

import (
	"net/http"

	"github.com/epointpayment/mloc-cpe/app/codes"
)

var (
	// ErrInvalidAPIKey is an error given when the requestor uses an invalid API key
	ErrInvalidAPIKey = codes.New("ERROR_API_KEY_INVALID").
				WithMessage("Invalid API Key").
				WithStatusCode(http.StatusForbidden).
				InGroup("API").
				RegisterError()

	// ErrInvalidProgramID is an error shown when the program ID is not valid
	ErrInvalidProgramID = codes.New("ERROR_PROGRAM_ID_INVALID").
				WithMessage("Invalid Program ID").
				WithStatusCode(http.StatusBadRequest).
				InGroup("API").
				RegisterError()

	// ErrInvalidProgramCustomerID is an error shown when customer ID is not valid
	ErrInvalidProgramCustomerID = codes.New("ERROR_CUSTOMER_ID_INVALID").
					WithMessage("Invalid Customer ID").
					WithStatusCode(http.StatusBadRequest).
					InGroup("API").
					RegisterError()

	// ErrCustomerExists is an error given when the customer was already created
	ErrCustomerExists = codes.New("ERROR_CUSTOMER_EXISTS").
				WithMessage("Customer already has an existing id and key").
				WithStatusCode(http.StatusBadRequest).
				InGroup("API").
				RegisterError()
)
