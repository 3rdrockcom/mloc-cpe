package models

import (
	"time"

	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/helpers"

	validation "github.com/go-ozzo/ozzo-validation"
)

// Transactions is an array of Transaction entries
type Transactions []Transaction

// Sort a list transactions ascending order
func (t Transactions) Len() int           { return len(t) }
func (t Transactions) Less(i, j int) bool { return t[i].Date.Time.Before(t[j].Date.Time) }
func (t Transactions) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }

// Separator splits a list of transactions into p partitions
func (t Transactions) Separator(p float64) []Transactions {
	MaxTransaction := 0.0

	for i := 0; i < len(t); i++ {
		if t[i].Credit > MaxTransaction {
			MaxTransaction = t[i].Credit
		}
	}

	threshold := MaxTransaction * p

	res := make([]Transactions, 2)
	for i := 0; i < len(t); i++ {
		if t[i].Credit >= threshold {
			k := 0
			res[k] = append(res[k], t[i])
		} else {
			k := 1
			res[k] = append(res[k], t[i])
		}
	}

	return res
}

// Transaction contains information about a transaction
type Transaction struct {
	ID             int          `json:"id"`
	CustomerID     int          `json:"customer_id" binding:"required"`
	Description    string       `json:"description" binding:"required"`
	Credit         float64      `json:"credit" binding:"required"`
	Debit          float64      `json:"debit" binding:"required"`
	RunningBalance float64      `json:"running_balance"`
	Balance        float64      `json:"balance"`
	Date           helpers.Time `json:"date" binding:"required" db:"-"`
	DateTime       time.Time    `json:"-" db:"datetime"`
}

// TableName gets the name of the database table
func (t Transaction) TableName() string {
	return "transactions"
}

// Validate checks if the values in the struct are valid
func (t Transaction) Validate() error {
	return validation.ValidateStruct(&t,
		validation.Field(&t.CustomerID, validation.Required),
		// validation.Field(&t.Credit, validation.Required),
		// validation.Field(&t.Debit, validation.Required),
		validation.Field(&t.Date, validation.Required),
	)
}
