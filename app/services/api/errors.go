package api

import "errors"

var (
	ErrInvalidAPIKey  = errors.New("Invalid API Key")
	ErrCustomerExists = errors.New("Customer already exists")
)
