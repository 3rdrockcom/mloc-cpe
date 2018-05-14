package repositories

import (
	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/models"

	dbx "github.com/go-ozzo/ozzo-dbx"
)

type Customers struct{}

func (c *Customers) Create(customer *models.Customer) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	err = tx.Model(customer).Insert()
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (c *Customers) Get(customerID int) (customer *models.Customer, err error) {
	customer = new(models.Customer)

	err = db.Select().
		Where(dbx.HashExp{"id": customerID}).
		One(customer)

	return customer, err
}

func (c *Customers) Update(customer *models.Customer) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	err = tx.Model(customer).Update()
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (c *Customers) GetByCustomerUniqueID(customerUniqueID string) (customer *models.Customer, err error) {
	customer = new(models.Customer)

	err = db.Select().
		Where(dbx.HashExp{"cust_unique_id": customerUniqueID}).
		One(customer)

	return customer, err
}
