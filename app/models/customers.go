package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type Customers []Customer

type Customer struct {
	ID                    int    `json:"id"`
	FirstName             string `json:"first_name" binding:"required"`
	LastName              string `json:"last_name" binding:"required"`
	Email                 string `json:"email" binding:"required"`
	ProgramID             int    `json:"-"`
	ProgramCustomerID     int    `json:"-"`
	ProgramCustomerMobile string `json:"-"`
	CustomerUniqueID      string `json:"-" db:"cust_unique_id"`
	LastTransactionID     int    `json:"-"`
}

func (c Customer) TableName() string {
	return "customers"
}

func (c Customer) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.FirstName, validation.Required),
		validation.Field(&c.LastName, validation.Required),
		validation.Field(&c.Email, validation.Required, is.Email),
	)
}
