package controllers

import (
	"net/http"
	"time"

	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/helpers"
	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/models"
	Customer "github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/services/customer"
	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/services/profiler"

	validation "github.com/go-ozzo/ozzo-validation"
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

// PostAddCustomer updates customer information
func (co Controllers) PostAddCustomer(c echo.Context) error {
	// Get customer ID
	customerID := c.Get("customerID").(int)

	// Initialize customer service
	sc, err := Customer.New(customerID)
	if err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	// Get customer information
	customer, err := sc.Info().Get()
	if err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	// Bind data to struct
	if err = c.Bind(customer); err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	// Validate struct
	if err = customer.Validate(); err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	// Update information
	if err = sc.Info().Update(customer); err != nil {
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
	Credit      decimal.Decimal `json:"credit" binding:"required"`
	Debit       decimal.Decimal `json:"debit" binding:"required"`
	Balance     decimal.Decimal `json:"balance"`
	DateTime    helpers.Time    `json:"date" binding:"required"`
}

// PostAddCustomerTransactions appends transactions to transaction list
func (co Controllers) PostAddCustomerTransactions(c echo.Context) error {
	// Get customer ID
	customerID := c.Get("customerID").(int)

	// Bind data to struct
	ctr := CustomerTransactionsRequest{}
	if err := c.Bind(&ctr); err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	// Assign values to struct
	transactions := models.Transactions{}
	for i := range ctr.Transactions {
		transactions = append(transactions, models.Transaction{
			CustomerID:  customerID,
			Description: ctr.Transactions[i].Description,
			Credit:      ctr.Transactions[i].Credit,
			Debit:       ctr.Transactions[i].Debit,
			Balance:     ctr.Transactions[i].Balance,
			DateTime:    ctr.Transactions[i].DateTime.Time,
		})
	}

	// Validate struct
	if err := validation.Validate(transactions); err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err.Error())
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
