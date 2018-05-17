package api

import "errors"

var (
	ErrInvalidAPIKey  = errors.New("Invalid API Key")
	ErrCustomerExists = errors.New("Customer already has an existing id and key")
)
