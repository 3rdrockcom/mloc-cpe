package customer

import (
	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/models"

	dbx "github.com/go-ozzo/ozzo-dbx"
)

type Info struct {
	cs *CustomerService
}

func (i *Info) Get() (customer *models.Customer, err error) {
	customer = new(models.Customer)

	err = DB.Select().
		Where(dbx.HashExp{"id": i.cs.CustomerID}).
		One(customer)
	if err != nil {
		return nil, err
	}

	return
}

func (i *Info) Update(customer *models.Customer) (err error) {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	customer.ID = i.cs.CustomerID
	err = tx.Model(customer).Update()
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
