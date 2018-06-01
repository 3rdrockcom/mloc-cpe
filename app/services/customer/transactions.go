package customer

import (
	"sort"
	"time"

	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/models"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/shopspring/decimal"
)

// Transactions manages customer transaction information
type Transactions struct {
	cs *CustomerService
}

// Create inserts an array of transactions into the database
func (t *Transactions) Create(transactions models.Transactions) (err error) {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	// Initialize customer info service
	customers := t.cs.Info()

	// Get customer information
	customer, err := customers.Get()
	if err != nil {
		tx.Rollback()
		return err
	}

	// Get current running balance
	runningBalance := decimal.Zero
	lastTransaction := models.Transaction{}
	if customer.LastTransactionID != 0 {
		lastTransaction, err = t.Get(customer.LastTransactionID)
		if err != nil {
			return err
		}
		runningBalance = lastTransaction.RunningBalance
	}

	// Sort transactions in ascending order
	sort.Sort(transactions)

	transaction := new(models.Transaction)
	for i := range transactions {
		*transaction = transactions[i]

		// Skip previously inserted transactions
		if customer.LastTransactionID != 0 {
			if transaction.DateTime.Before(lastTransaction.DateTime) || transaction.DateTime.Equal(lastTransaction.DateTime) {
				continue
			}
		}

		// Calculate running balance
		runningBalance = runningBalance.Add(transaction.Credit.Sub(transaction.Debit))
		transaction.RunningBalance = runningBalance

		// Insert into database
		err = tx.Model(transaction).Insert()
		if err != nil {
			tx.Rollback()
			return err
		}

		// Update last transaction ID marker
		if i+1 == len(transactions) {
			customer.LastTransactionID = transaction.ID
			err = tx.Model(customer).Update("LastTransactionID")
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return
	}

	return
}

// Get gets a specific transaction by ID
func (t *Transactions) Get(transactionID int) (transaction models.Transaction, err error) {
	customerID := t.cs.CustomerID

	err = DB.Select().
		Where(dbx.HashExp{"id": transactionID, "customer_id": customerID}).
		One(&transaction)

	return transaction, err
}

// GetAll gets all stored transactions
func (t *Transactions) GetAll() (transactions models.Transactions, err error) {
	customerID := t.cs.CustomerID

	err = DB.Select().
		Where(dbx.HashExp{"customer_id": customerID}).
		All(&transactions)

	return transactions, err
}

// GetAllByDateRange gets all transactions from a specified date range
func (t *Transactions) GetAllByDateRange(startDate, endDate time.Time) (transactions models.Transactions, err error) {
	customerID := t.cs.CustomerID

	err = DB.Select().
		Where(dbx.HashExp{"customer_id": customerID}).
		AndWhere(dbx.Between("datetime", startDate, endDate)).
		All(&transactions)

	return transactions, err
}
