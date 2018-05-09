package repositories

import (
	"database/sql"
	"sort"
	"time"

	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/models"
	dbx "github.com/go-ozzo/ozzo-dbx"
)

type Transactions struct{}

func (t *Transactions) Create(customerID int, transactions models.Transactions) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	runningBalance := 0.0
	lastTransaction, err := t.GetLastTransaction(customerID)
	if err != nil && err != sql.ErrNoRows {
		return err
	} else {
		runningBalance = lastTransaction.RunningBalance
	}

	sort.Sort(transactions)

	transaction := new(models.Transaction)
	for i := range transactions {
		*transaction = transactions[i]
		if transaction.DateTime.Before(lastTransaction.DateTime) || transaction.DateTime.Equal(lastTransaction.DateTime) {
			continue
		}

		runningBalance = runningBalance + transaction.Credit - transaction.Debit
		transaction.RunningBalance = runningBalance

		err = tx.Model(transaction).Insert()
		if err != nil {
			tx.Rollback()
			return err
		}

		if i+1 == len(transactions) {
			customers := new(Customers)
			customer, err := customers.Get(customerID)
			if err != nil {
				tx.Rollback()
				return err
			}

			customer.UpdatedAt = transaction.DateTime
			err = tx.Model(customer).Update("UpdatedAt")
			if err != nil {
				tx.Rollback()
				return err
			}
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

func (t *Transactions) GetLastTransaction(customerID int) (transaction models.Transaction, err error) {
	err = db.Select().
		Where(dbx.HashExp{"customer_id": customerID}).
		OrderBy("id DESC").
		One(&transaction)

	return transaction, err
}
