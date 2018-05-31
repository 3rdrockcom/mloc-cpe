package api

import "errors"

var (
	// ErrMissingAPIKey is an error given when the request is missing an API key
	ErrMissingAPIKey = errors.New("Missing API Key")

	// ErrInvalidAPIKey is an error given when the requestor uses an invalid API key
	ErrInvalidAPIKey = errors.New("Invalid API Key")

	// ErrCustomerExists is an error given when the customer was already created
	ErrCustomerExists = errors.New("Customer already has an existing id and key")
)
