package customer

import (
	dbx "github.com/go-ozzo/ozzo-dbx"
	validation "github.com/go-ozzo/ozzo-validation"
)

var DB *dbx.DB

type CustomerService struct {
	CustomerID   int
	info         *Info
	transactions *Transactions
}

func (cs CustomerService) Validate() error {
	return validation.ValidateStruct(&cs,
		validation.Field(&cs.CustomerID, validation.Required),
	)
}

func New(customerID int) (cs *CustomerService, err error) {
	// validate

	cs = new(CustomerService)
	cs.CustomerID = customerID
	err = cs.Validate()
	if err != nil {
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

func (cs *CustomerService) Info() *Info {
	return cs.info
}

func (cs *CustomerService) Transactions() *Transactions {
	return cs.transactions
}
