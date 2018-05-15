package api

import (
	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/models"

	dbx "github.com/go-ozzo/ozzo-dbx"
)

var DB *dbx.DB

type APIService struct{}

func New() *APIService {
	return &APIService{}
}

func (as *APIService) GetLoginKey() (entry *models.APIKey, err error) {
	entry = new(models.APIKey)

	err = DB.Select().
		Where(dbx.HashExp{"key": "LOGIN"}).
		One(entry)
	if err != nil {
		return nil, err
	}

	return
}

func (as *APIService) GetRegistrationKey() (entry *models.APIKey, err error) {
	entry = new(models.APIKey)

	err = DB.Select().
		Where(dbx.HashExp{"customer_id": 0}).
		AndWhere(dbx.NewExp("`key`!={:key}", dbx.Params{"key": "LOGIN"})).
		One(entry)
	if err != nil {
		return nil, err
	}

	return
}

func (as *APIService) GetCustomerKey(key string) (entry *models.APIKey, err error) {
	entry, err = as.GetKey(key)
	if err != nil {
		return nil, err
	}

	return
}

func (as *APIService) GetKey(key string) (entry *models.APIKey, err error) {
	entry = new(models.APIKey)

	err = DB.Select().
		Where(dbx.HashExp{"key": key}).
		One(entry)

	return
}

func (as *APIService) GetKeyByCustomerID(customerID int) (entry *models.APIKey, err error) {
	entry = new(models.APIKey)

	err = DB.Select().
		Where(dbx.HashExp{"customer_id": customerID}).
		One(entry)

	return
}

func (as *APIService) GetCustomerByCustomerUniqueID(customerUniqueID string) (customer *models.Customer, err error) {
	customer = new(models.Customer)

	err = DB.Select().
		Where(dbx.HashExp{"cust_unique_id": customerUniqueID}).
		One(customer)
	if err != nil {
		return nil, err
	}

	return
}

func (as *APIService) GetCustomerAccessKey(programID int, programCustomerID int, programCustomerMobile string) (k Key, err error) {
	customerKey, err := NewKey(programID, programCustomerID, programCustomerMobile)
	if err != nil {
		return
	}

	k, err = customerKey.GetCustomerKey()
	if err != nil {
		return
	}

	return
}

func (as *APIService) GenerateCustomerAccessKey(programID int, programCustomerID int, programCustomerMobile string) (k Key, err error) {
	customerKey, err := NewKey(programID, programCustomerID, programCustomerMobile)
	if err != nil {
		return
	}

	k, err = customerKey.GenerateCustomerKey()
	if err != nil {
		return
	}

	return
}
