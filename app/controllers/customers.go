package controllers

import (
	"net/http"
	"time"

	"github.com/epointpayment/mloc-cpe/app/helpers"
	"github.com/epointpayment/mloc-cpe/app/models"
	Customer "github.com/epointpayment/mloc-cpe/app/services/customer"
	Profiler "github.com/epointpayment/mloc-cpe/app/services/profiler"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/jinzhu/now"
	"github.com/juju/errors"
	"github.com/labstack/echo"
	"github.com/shopspring/decimal"
)

// GetCustomer displays detailed customer information
func (co *Controllers) GetCustomer(c echo.Context) (err error) {
	// Get customer ID
	customerID := c.Get("customerID").(int)

	// Initialize customer service
	sc, err := Customer.New(customerID)
	if err != nil {
		err = errors.Trace(err)
		return
	}

	// Get detailed customer information
	customerInfo, err := sc.Info().GetDetails()
	if err != nil {
		err = errors.Trace(err)
		return
	}

	return SendResponse(c, http.StatusOK, customerInfo)
}

// CustomerRequest contains information about a customer
type CustomerRequest struct {
	FirstName    string `json:"first_name" form:"first_name"`
	LastName     string `json:"last_name" form:"last_name"`
	Email        string `json:"email" form:"email"`
	MobileNumber string `json:"mobile_number" form:"mobile_number"`
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
func (co Controllers) PostAddCustomer(c echo.Context) (err error) {
	// Get customer ID
	customerID := c.Get("customerID").(int)

	// Initialize customer service
	sc, err := Customer.New(customerID)
	if err != nil {
		err = errors.Trace(err)
		return
	}

	cr := new(CustomerRequest)
	customer := new(models.Customer)
	customer.ID = customerID

	// Bind data to struct
	if err = c.Bind(cr); err != nil {
		err = errors.Wrap(err, Customer.ErrProblemOccurred)
		return
	}

	// Validate struct
	if err = cr.Validate(); err != nil {
		err = errors.Wrap(err, Customer.ErrCustomerIncompleteInfo)
		return
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
	if err = sc.Info().UpdateCustomerBasic(customer, fields...); err != nil {
		err = errors.Wrap(err, Customer.ErrProblemOccurred)
		return
	}

	// Send response
	return SendOKResponse(c, Customer.MsgInfoUpdated)
}

// CustomerTransactionsResponse is a wrapper for transaction data
type CustomerTransactionsResponse struct {
	Transactions []CustomerTransactionResponse `json:"transactions"`
}

// CustomerTransactionResponse contains information about a transaction
type CustomerTransactionResponse struct {
	Type           string `json:"type"`
	Description    string `json:"description"`
	Value          string `json:"amount"`
	Balance        string `json:"balance"`
	RunningBalance string `json:"running_balance"`
	Date           string `json:"date"`
}

// GetCustomerProfile generates a profile of a customer
func (co Controllers) GetCustomerTransactions(c echo.Context) (err error) {
	// Get customer ID
	customerID := c.Get("customerID").(int)

	// Get transaction start date
	startDate, err := time.ParseInLocation(
		"20060102",
		c.QueryParam("startDate"), time.UTC)
	if err != nil {
		err = errors.Wrap(err, Customer.ErrInvalidDate)
		return
	}

	// Get transaction end date
	endDate, err := time.ParseInLocation(
		"20060102",
		c.QueryParam("endDate"), time.UTC)
	if err != nil {
		err = errors.Wrap(err, Customer.ErrInvalidDate)
		return
	}

	// Initialize customer service
	sc, err := Customer.New(customerID)
	if err != nil {
		err = errors.Trace(err)
		return
	}

	// Get all customer transaction within the specified range
	transactions := models.Transactions{}
	if transactions, err = sc.Transactions().GetAllByDateRange(startDate, now.New(endDate).EndOfDay()); err != nil {
		err = errors.Wrap(err, Customer.ErrProblemOccurred)
		return
	}
	if len(transactions) == 0 {
		err = errors.Wrap(err, Customer.ErrNoTransactionsAvailable)
		return
	}

	res := CustomerTransactionsResponse{}
	for i := range transactions {
		t := CustomerTransactionResponse{
			Date:           transactions[i].DateTime.Format("2006-01-02 15:04:05"),
			Description:    transactions[i].Description,
			Balance:        transactions[i].Balance.StringFixed(helpers.DefaultCurrencyPrecision),
			RunningBalance: transactions[i].RunningBalance.StringFixed(helpers.DefaultCurrencyPrecision),
		}

		switch {
		case !transactions[i].Credit.IsZero():
			t.Type = "credit"
			t.Value = transactions[i].Credit.StringFixed(helpers.DefaultCurrencyPrecision)
		case !transactions[i].Debit.IsZero():
			t.Type = "debit"
			t.Value = transactions[i].Debit.StringFixed(helpers.DefaultCurrencyPrecision)
		}

		res.Transactions = append(res.Transactions, t)
	}

	return SendResponse(c, http.StatusOK, res)
}

// CustomerTransactionsRequest is a wrapper for transaction data
type CustomerTransactionsRequest struct {
	Transactions []CustomerTransactionRequest `json:"transactions" form:"transactions"`
}

// CustomerTransactionRequest contains information about a transaction
type CustomerTransactionRequest struct {
	Description string          `json:"description" form:"description"`
	Type        string          `json:"type" form:"type"`
	Value       decimal.Decimal `json:"amount" form:"amount"`
	Balance     decimal.Decimal `json:"balance" form:"balance"`
	Date        string          `json:"date" form:"date"`
}

// Validate checks if the values in the struct are valid
func (t CustomerTransactionRequest) Validate() error {
	err := validation.ValidateStruct(&t,
		validation.Field(&t.Description, validation.Required),
		validation.Field(&t.Type, validation.Required, validation.In("credit", "debit")),
		validation.Field(&t.Value, validation.Required, validation.By(helpers.ValidateCurrencyAmount)),
		validation.Field(&t.Date, validation.Required, validation.Date("2006-01-02 15:04:05")),
	)
	if err != nil {
		return err
	}

	switch {
	case t.Type == "credit" && t.Value.LessThan(decimal.Zero):
		return Customer.ErrCreditNonPositiveValue
	}

	return nil
}

// PostAddCustomerTransactions appends transactions to transaction list
func (co Controllers) PostAddCustomerTransactions(c echo.Context) (err error) {
	// Get customer ID
	customerID := c.Get("customerID").(int)

	// Bind data to struct
	ctr := CustomerTransactionsRequest{}
	if err := c.Bind(&ctr); err != nil {
		err = errors.Wrap(err, Customer.ErrProblemOccurred)
		return err
	}

	// Validate struct
	if err := validation.Validate(ctr.Transactions); err != nil {
		err = errors.Wrap(err, Customer.ErrInvalidData)
		return err
	}

	// Assign values to struct
	transactions := models.Transactions{}
	for i := range ctr.Transactions {
		transaction := models.Transaction{
			CustomerID:  customerID,
			Description: ctr.Transactions[i].Description,
			Balance:     ctr.Transactions[i].Balance,
		}

		transaction.DateTime, _ = time.Parse("2006-01-02 15:04:05", ctr.Transactions[i].Date)

		switch ctr.Transactions[i].Type {
		case "credit":
			transaction.Credit = ctr.Transactions[i].Value
		case "debit":
			transaction.Debit = ctr.Transactions[i].Value
		}

		transactions = append(transactions, transaction)
	}

	// Initialize customer service
	sc, err := Customer.New(customerID)
	if err != nil {
		err = errors.Trace(err)
		return
	}

	// Insert new transactions
	if err = sc.Transactions().Create(transactions); err != nil {
		err = errors.Wrap(err, Customer.ErrProblemOccurred)
		return
	}

	// Send response
	return SendOKResponse(c, Customer.MsgTransactionsInserted)
}

// GetCustomerProfile generates a profile of a customer
func (co Controllers) GetCustomerProfile(c echo.Context) (err error) {
	// Get customer ID
	customerID := c.Get("customerID").(int)

	// Get transaction start date
	startDate, err := time.ParseInLocation(
		"20060102",
		c.QueryParam("startDate"), time.UTC)
	if err != nil {
		err = errors.Wrap(err, Customer.ErrInvalidDate)
		return
	}

	// Get transaction end date
	endDate, err := time.ParseInLocation(
		"20060102",
		c.QueryParam("endDate"), time.UTC)
	if err != nil {
		err = errors.Wrap(err, Customer.ErrInvalidDate)
		return
	}

	// Initialize customer service
	sc, err := Customer.New(customerID)
	if err != nil {
		err = errors.Trace(err)
		return
	}

	// Get all customer transaction within the specified range
	transactions := models.Transactions{}
	if transactions, err = sc.Transactions().GetAllByDateRange(startDate, now.New(endDate).EndOfDay()); err != nil {
		err = errors.Wrap(err, Customer.ErrProblemOccurred)
		return
	}
	if len(transactions) == 0 {
		err = errors.Wrap(err, Customer.ErrNoTransactionsAvailable)
		return
	}

	// Initialize profiler service and run analysis
	p := Profiler.New(transactions, Profiler.DefaultPartitions)
	res, err := p.Run()
	if err != nil {
		err = errors.Wrap(err, Customer.ErrProblemOccurred)
		return
	}

	return SendResponse(c, http.StatusOK, res)
}
