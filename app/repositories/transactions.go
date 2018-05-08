package repositories

import (
	"time"

	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/models"
	dbx "github.com/go-ozzo/ozzo-dbx"
)

type Transactions struct{}

func (t *Transactions) Create(transactions models.Transactions) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	transaction := new(models.Transaction)
	for _, *transaction = range transactions {
		err = tx.Model(transaction).Insert()
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

func (t *Transactions) Get(customerID int) (transactions models.Transactions, err error) {
	err = db.Select().
		Where(dbx.HashExp{"customer_id": customerID}).
		All(&transactions)

	return transactions, err
}

func (t *Transactions) GetByDateRange(customerID int, startDate, endDate time.Time) (transactions models.Transactions, err error) {
	err = db.Select().
		Where(dbx.HashExp{"customer_id": customerID}).
		AndWhere(dbx.Between("datetime", startDate, endDate)).
		All(&transactions)

	return transactions, err
}
