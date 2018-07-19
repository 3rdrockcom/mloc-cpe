package customer

import (
	dbx "github.com/go-ozzo/ozzo-dbx"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/juju/errors"
)

// DB is the database handler
var DB *dbx.DB

// CustomerService is a service that manages a customer
type CustomerService struct {
	CustomerID   int
	info         *Info
	transactions *Transactions
}

// Validate checks if the values in the struct are valid
func (cs CustomerService) Validate() error {
	return validation.ValidateStruct(&cs,
		validation.Field(&cs.CustomerID, validation.Required),
	)
}

// New creates an instance of the customer service
func New(customerID int) (cs *CustomerService, err error) {
	cs = new(CustomerService)
	cs.CustomerID = customerID
	err = cs.Validate()
	if err != nil {
		err = errors.Trace(err)
		return
	}

	cs.info = &Info{
		cs: cs,
	}
	cs.transactions = &Transactions{
		cs: cs,
	}

	return
}

// Info gets customer info methods
func (cs *CustomerService) Info() *Info {
	return cs.info
}

// Transactions gets customer transaction methods
func (cs *CustomerService) Transactions() *Transactions {
	return cs.transactions
}
