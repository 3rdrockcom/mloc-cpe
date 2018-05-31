package customer

import (
	"database/sql"

	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/models"

	dbx "github.com/go-ozzo/ozzo-dbx"
)

// Info manages customer information
type Info struct {
	cs *CustomerService
}

// Get gets basic customer information
func (i *Info) Get() (customer *models.Customer, err error) {
	customer = new(models.Customer)

	err = DB.Select().
		Where(dbx.HashExp{"id": i.cs.CustomerID}).
		One(customer)
	if err != nil {
		if err == sql.ErrNoRows {
			err = ErrCustomerNotFound
		}
		return nil, err
	}

	return
}

// CustomerDetails holds detailed customer information
type CustomerDetails struct {
	ID                    int    `json:"id"`
	FirstName             string `json:"first_name"`
	LastName              string `json:"last_name"`
	Email                 string `json:"email"`
	MobileNumber          string `json:"mobile_number"`
	ProgramID             int    `json:"program_id"`
	ProgramCustomerID     int    `json:"program_customer_id"`
	ProgramCustomerMobile string `json:"program_customer_mobile"`
	CustomerUniqueID      string `json:"cust_unique_id" db:"cust_unique_id"`
	Key                   string `json:"key" db:"key"`
}

// GetDetails gets detailed customer information
func (i *Info) GetDetails() (customerDetails *CustomerDetails, err error) {
	customerDetails = new(CustomerDetails)

	err = DB.Select("customers.*", "api_keys.key").
		From("customers").
		LeftJoin("api_keys", dbx.NewExp("api_keys.customer_id = customers.id")).
		Where(dbx.HashExp{"customers.id": i.cs.CustomerID}).
		One(customerDetails)
	if err != nil {
		if err == sql.ErrNoRows {
			err = ErrCustomerNotFound
		}
		return nil, err
	}

	return
}

// Update updates customer information
func (i *Info) Update(customer *models.Customer) (err error) {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	customer.ID = i.cs.CustomerID
	err = tx.Model(customer).Update()
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
