package customer

import "errors"

var (
	ErrInvalidUniqueCustomerID = errors.New("Invalid Unique Customer ID")
	ErrCustomerNotFound        = errors.New("Customer not found")
)
