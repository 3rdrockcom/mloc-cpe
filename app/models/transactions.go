package models

import (
	"time"

	"github.com/shopspring/decimal"
)

// Transactions is an array of Transaction entries
type Transactions []Transaction

// Sort a list transactions ascending order
func (t Transactions) Len() int           { return len(t) }
func (t Transactions) Less(i, j int) bool { return t[i].DateTime.Before(t[j].DateTime) }
func (t Transactions) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }

// Separator splits a list of transactions into p partitions
func (t Transactions) Separator(p float64) []Transactions {
	MaxTransaction := decimal.Zero
	partitions := decimal.NewFromFloat(p)

	for i := 0; i < len(t); i++ {
		if t[i].Credit.GreaterThan(MaxTransaction) {
			MaxTransaction = t[i].Credit
		}
	}

	threshold := MaxTransaction.Mul(partitions)

	res := make([]Transactions, 2)
	for i := 0; i < len(t); i++ {
		if t[i].Credit.GreaterThanOrEqual(threshold) {
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
	ID             int             `json:"id"`
	CustomerID     int             `json:"customer_id" binding:"required"`
	Description    string          `json:"description" binding:"required"`
	Credit         decimal.Decimal `json:"credit" binding:"required"`
	Debit          decimal.Decimal `json:"debit" binding:"required"`
	RunningBalance decimal.Decimal `json:"running_balance"`
	Balance        decimal.Decimal `json:"balance"`
	DateTime       time.Time       `json:"-" db:"datetime"`
}

// TableName gets the name of the database table
func (t Transaction) TableName() string {
	return "transactions"
}
