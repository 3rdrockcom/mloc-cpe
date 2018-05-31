package customer

import "errors"

var (
	// ErrMissingUniqueCustomerID is an error shown when a request is missing the customer unique ID
	ErrMissingUniqueCustomerID = errors.New("Missing Unique Customer ID")

	// ErrInvalidUniqueCustomerID is an error shown when customer unique ID is not a valid
	ErrInvalidUniqueCustomerID = errors.New("Invalid Unique Customer ID")

	// ErrCustomerNotFound is an error for a non-existent customer
	ErrCustomerNotFound = errors.New("Customer not found")
)
