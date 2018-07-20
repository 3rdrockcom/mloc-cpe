package models

// Customers is an array of Customer entries
type Customers []Customer

// Customer contains information about a customer
type Customer struct {
	ID                    int    `json:"id"`
	FirstName             string `json:"first_name" form:"first_name"`
	LastName              string `json:"last_name" form:"last_name"`
	Email                 string `json:"email" form:"email"`
	MobileNumber          string `json:"mobile_number" form:"mobile_number"`
	ProgramID             int    `json:"-"`
	ProgramCustomerID     int    `json:"-"`
	ProgramCustomerMobile string `json:"-"`
	CustomerUniqueID      string `json:"-" db:"cust_unique_id"`
	LastTransactionID     int    `json:"-"`
}

// TableName gets the name of the database table
func (c Customer) TableName() string {
	return "customers"
}
