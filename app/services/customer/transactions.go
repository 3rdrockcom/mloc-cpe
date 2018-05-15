package customer

import (
	"sort"
	"time"

	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/models"

	dbx "github.com/go-ozzo/ozzo-dbx"
)

type Transactions struct {
	cs *CustomerService
}

func (t *Transactions) Create(transactions models.Transactions) (err error) {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	customers := t.cs.Info()
	customer, err := customers.Get()
	if err != nil {
		tx.Rollback()
		return err
	}

	runningBalance := 0.0
	lastTransaction := models.Transaction{}
	if customer.LastTransactionID != 0 {
		lastTransaction, err = t.Get(customer.LastTransactionID)
		if err != nil {
			return err
		}
		runningBalance = lastTransaction.RunningBalance
	}

	sort.Sort(transactions)

	transaction := new(models.Transaction)
	for i := range transactions {
		*transaction = transactions[i]

		if customer.LastTransactionID != 0 {
			if transaction.DateTime.Before(lastTransaction.DateTime) || transaction.DateTime.Equal(lastTransaction.DateTime) {
				continue
			}
		}

		runningBalance = runningBalance + transaction.Credit - transaction.Debit
		transaction.RunningBalance = runningBalance

		err = tx.Model(transaction).Insert()
		if err != nil {
			tx.Rollback()
			return err
		}

		if i+1 == len(transactions) {
			customer.LastTransactionID = transaction.ID
			err = tx.Model(customer).Update("LastTransactionID")
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	tx.Commit()
	return nil
}

func (t *Transactions) Get(transactionID int) (transaction models.Transaction, err error) {
	customerID := t.cs.CustomerID
	err = DB.Select().
		Where(dbx.HashExp{"id": transactionID, "customer_id": customerID}).
		One(&transaction)

	return transaction, err
}

func (t *Transactions) GetAll() (transactions models.Transactions, err error) {
	customerID := t.cs.CustomerID
	err = DB.Select().
		Where(dbx.HashExp{"customer_id": customerID}).
		All(&transactions)

	return transactions, err
}

func (t *Transactions) GetAllByDateRange(startDate, endDate time.Time) (transactions models.Transactions, err error) {
	customerID := t.cs.CustomerID
	err = DB.Select().
		Where(dbx.HashExp{"customer_id": customerID}).
		AndWhere(dbx.Between("datetime", startDate, endDate)).
		All(&transactions)

	return transactions, err
}
