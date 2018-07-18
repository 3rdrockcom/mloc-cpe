package controllers

import (
	"net/http"
	"strings"
	"time"

	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/models"
	Customer "github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/services/customer"
	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/services/profiler"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/jinzhu/now"
	"github.com/labstack/echo"
	"github.com/shopspring/decimal"
)

// GetCustomer gets customer information
func (co Controllers) GetCustomer(c echo.Context) error {
	// Get customer ID
	customerID := c.Get("customerID").(int)

	// Initialize customer service
	sc, err := Customer.New(customerID)
	if err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	// Get customer information
	customer, err := sc.Info().GetDetails()
	if err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	// Send response
	return SendResponse(c, http.StatusOK, customer)
}

// CustomerRequest contains information about a customer
type CustomerRequest struct {
	FirstName    string `json:"first_name" form:"first_name" binding:"required"`
	LastName     string `json:"last_name" form:"last_name" binding:"required"`
	Email        string `json:"email" form:"email" binding:"required"`
	MobileNumber string `json:"mobile_number" form:"mobile_number" binding:"required"`
}

// Validate checks if the values in the struct are valid
func (c CustomerRequest) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.FirstName, validation.Required),
		validation.Field(&c.LastName, validation.Required),
		validation.Field(&c.Email, validation.Required, is.Email),
		validation.Field(&c.MobileNumber, validation.Required),
	)
}

// PostAddCustomer updates customer information
func (co Controllers) PostAddCustomer(c echo.Context) error {
	// Get customer ID
	customerID := c.Get("customerID").(int)

	// Initialize customer service
	sc, err := Customer.New(customerID)
	if err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	cr := new(CustomerRequest)
	customer := new(models.Customer)
	customer.ID = customerID

	// Bind data to struct
	if err = c.Bind(cr); err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	// Validate struct
	if err = cr.Validate(); err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	// Prepare customer information
	fields := []string{"FirstName", "LastName", "Email", "MobileNumber"}
	for i := range fields {
		switch fields[i] {
		case "FirstName":
			customer.FirstName = cr.FirstName
		case "LastName":
			customer.LastName = cr.LastName
		case "Email":
			customer.Email = cr.Email
		case "MobileNumber":
			customer.MobileNumber = cr.MobileNumber
		}
	}

	// Update information
	if err = sc.Info().Update(customer, fields...); err != nil {
		return SendErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	// Send response
	return SendOKResponse(c, Customer.MsgInfoUpdated)
}

// CustomerTransactionsRequest is a wrapper for transaction data
type CustomerTransactionsRequest struct {
	Transactions []CustomerTransactionRequest `json:"transactions" binding:"required"`
}

// CustomerTransactionRequest contains information about a transaction
type CustomerTransactionRequest struct {
	Description string          `json:"description" binding:"required"`
	Type        string          `json:"type" binding:"required"`
	Value       decimal.Decimal `json:"amount" binding:"required"`
	Balance     decimal.Decimal `json:"balance"`
	Timestamp   int64           `json:"timestamp" binding:"required"`
}

// Validate checks if the values in the struct are valid
func (t CustomerTransactionRequest) Validate() error {
	switch {
	case t.Type == "credit" && t.Value.LessThan(decimal.Zero):
		return Customer.ErrCreditNonPositiveValue

	case t.Timestamp <= 0:
		return Customer.ErrInvalidTimestamp

	//
	// case t.Credit.Equal(decimal.Zero) && t.Debit.Equal(decimal.Zero):
	// 	return Customer.ErrCreditDebitRequired
	// case !t.Credit.Equal(decimal.Zero) && !t.Debit.Equal(decimal.Zero):
	// 	return Customer.ErrCreditDebitRequired
	// case t.Credit.LessThan(decimal.Zero):
	// 	return Customer.ErrCreditNonPositiveValue
	case t.Timestamp <= 0:
		return Customer.ErrInvalidTimestamp
	}

	return nil
}

// PostAddCustomerTransactions appends transactions to transaction list
func (co Controllers) PostAddCustomerTransactions(c echo.Context) error {
	// Get customer ID
	customerID := c.Get("customerID").(int)

	// Bind data to struct
	ctr := CustomerTransactionsRequest{}
	if err := c.Bind(&ctr); err != nil {
		err = Customer.ErrInvalidData
		return err
	}

	// Validate struct
	if err := validation.Validate(ctr.Transactions); err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	// Assign values to struct
	transactions := models.Transactions{}
	for i := range ctr.Transactions {
		transaction := models.Transaction{
			CustomerID:  customerID,
			Description: ctr.Transactions[i].Description,
			Balance:     ctr.Transactions[i].Balance,
			DateTime:    time.Unix(ctr.Transactions[i].Timestamp, 0),
		}

		switch strings.ToUpper(ctr.Transactions[i].Type) {
		case "CREDIT":
			transaction.Credit = ctr.Transactions[i].Value
		case "DEBIT":
			transaction.Debit = ctr.Transactions[i].Value
		}

		transactions = append(transactions, transaction)
	}

	// Initialize customer service
	sc, err := Customer.New(customerID)
	if err != nil {
		return err
	}

	// Insert new transactions
	if err = sc.Transactions().Create(transactions); err != nil {
		return SendErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	// Send response
	return SendOKResponse(c, Customer.MsgTransactionsInserted)
}

// GetCustomerProfile generates a profile of a customer
func (co Controllers) GetCustomerProfile(c echo.Context) error {
	// Get customer ID
	customerID := c.Get("customerID").(int)

	// Get transaction start date
	startDate, err := time.ParseInLocation(
		"20060102",
		c.QueryParam("startDate"), time.UTC)
	if err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	// Get transaction end date
	endDate, err := time.ParseInLocation(
		"20060102",
		c.QueryParam("endDate"), time.UTC)
	if err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	// Initialize customer service
	sc, err := Customer.New(customerID)
	if err != nil {
		return err
	}

	// Get all customer transaction within the specified range
	transactions := models.Transactions{}
	if transactions, err = sc.Transactions().GetAllByDateRange(startDate, now.New(endDate).EndOfDay()); err != nil {
		return SendErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	// Initialize profiler service and run analysis
	p := profiler.New(transactions, 2)
	res := p.Run()

	return SendResponse(c, http.StatusOK, res)
}
