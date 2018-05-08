package repositories

import (
	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/models"
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
