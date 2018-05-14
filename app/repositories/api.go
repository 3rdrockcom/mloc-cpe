package repositories

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"strconv"
	"time"

	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/models"
	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/labstack/gommon/random"
)

type API struct{}

func (a API) Create(apiKey *models.APIKey) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	apiKey.DateCreated = time.Now().UTC()
	err = tx.Model(apiKey).Insert()
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

type CustomerKey struct {
	CustomerUniqueID string `json:"cust_unique_id"`
	ApiKey           string `json:"api_key"`
}

func (a API) GetCustomerKey(programID int, programCustomerID int, programCustomerMobile string) (customerKey CustomerKey, err error) {
	customerUniqueID := a.GenerateCustomerUniqueID(programID, programCustomerID, programCustomerMobile)

	customers := new(Customers)
	customer, err := customers.GetByCustomerUniqueID(customerUniqueID)
	if err != nil {
		if err == sql.ErrNoRows {
			regKey, err := a.GetAPIRegistrationKey()
			customerKey.ApiKey = regKey.Key
			return customerKey, err
		}
		return customerKey, err
	}

	apiKey, err := a.GetAPIKeyByCustomerID(customer.ID)
	if err != nil {
		return customerKey, err
	}

	customerKey.ApiKey = apiKey.Key
	customerKey.CustomerUniqueID = customer.CustomerUniqueID

	return customerKey, nil
}

func (a API) GenerateCustomerKey(programID int, programCustomerID int, programCustomerMobile string) (customerKey CustomerKey, err error) {
	customerUniqueID := a.GenerateCustomerUniqueID(programID, programCustomerID, programCustomerMobile)

	customers := new(Customers)
	_, err = customers.GetByCustomerUniqueID(customerUniqueID)
	if err != nil && err != sql.ErrNoRows {
		return customerKey, err
	}

	if err == nil {
		err = errors.New("Customer already exists")
		return customerKey, err
	}

	customer := new(models.Customer)
	customer.ProgramID = programID
	customer.ProgramCustomerID = programCustomerID
	customer.ProgramCustomerMobile = programCustomerMobile
	customer.CustomerUniqueID = customerUniqueID

	err = customers.Create(customer)
	if err != nil {
		return customerKey, err
	}

	apiKey := new(models.APIKey)
	apiKey.CustomerID = customer.ID
	apiKey.Key = random.String(32)
	err = a.Create(apiKey)
	if err != nil {
		return customerKey, err
	}

	customerKey.ApiKey = apiKey.Key
	customerKey.CustomerUniqueID = customer.CustomerUniqueID

	return customerKey, err
}

func (a API) GenerateCustomerUniqueID(programID int, programCustomerID int, programCustomerMobile string) string {
	str := strconv.Itoa(programID) + "_" + strconv.Itoa(programCustomerID) + "_" + programCustomerMobile

	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func (a API) GetAPIKeyByCustomerID(customerID int) (apiKey *models.APIKey, err error) {
	apiKey = new(models.APIKey)

	err = db.Select().
		Where(dbx.HashExp{"customer_id": customerID}).
		One(apiKey)

	return apiKey, err
}

func (a API) GetAPICustomerKey(key string) (apiKey *models.APIKey, err error) {
	apiKey = new(models.APIKey)

	err = db.Select().
		Where(dbx.HashExp{"key": key}).
		One(apiKey)

	return apiKey, err
}

func (a API) GetAPIRegistrationKey() (apiKey *models.APIKey, err error) {
	apiKey = new(models.APIKey)

	err = db.Select().
		Where(dbx.HashExp{"customer_id": 0}).
		AndWhere(dbx.NewExp("`key`!={:key}", dbx.Params{"key": "LOGIN"})).
		One(apiKey)

	return apiKey, err
}

func (a API) GetAPILoginKey() (apiKey *models.APIKey, err error) {
	apiKey = new(models.APIKey)

	err = db.Select().
		Where(dbx.HashExp{"key": "LOGIN"}).
		One(apiKey)

	return apiKey, err
}
