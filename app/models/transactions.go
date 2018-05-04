package models

import (
	"time"

	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/helpers"

	validation "github.com/go-ozzo/ozzo-validation"
)

type Transactions []Transaction

type Transaction struct {
	ID          int          `json:"id"`
	CustomerID  int          `json:"customer_id" binding:"required"`
	Description string       `json:"description" binding:"required"`
	Credit      float64      `json:"credit" binding:"required"`
	Debit       float64      `json:"debit" binding:"required"`
	Date        helpers.Time `json:"date" binding:"required" db:"-"`
	DateTime    time.Time    `json:"-" db:"datetime"`
}

func (t Transaction) TableName() string {
	return "transactions"
}

func (t Transaction) Validate() error {
	return validation.ValidateStruct(&t,
		validation.Field(&t.CustomerID, validation.Required),
		// validation.Field(&t.Credit, validation.Required),
		// validation.Field(&t.Debit, validation.Required),
		validation.Field(&t.Date, validation.Required),
	)
}
