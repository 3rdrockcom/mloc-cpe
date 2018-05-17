package api

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"strconv"
	"time"

	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/models"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/labstack/gommon/random"
)

type Key struct {
	programID             int
	programCustomerID     int
	programCustomerMobile string
	CustomerUniqueID      string `json:"cust_unique_id"`
	ApiKey                string `json:"api_key"`
}

func (k Key) Validate() error {
	return validation.ValidateStruct(&k,
		validation.Field(&k.programID, validation.Required),
		validation.Field(&k.programCustomerID, validation.Required),
		validation.Field(&k.programCustomerMobile, validation.Required),
	)
}

func NewKey(programID int, programCustomerID int, programCustomerMobile string) (k *Key, err error) {
	k = &Key{
		programID:             programID,
		programCustomerID:     programCustomerID,
		programCustomerMobile: programCustomerMobile,
	}

	err = k.Validate()
	return
}

func (k *Key) GetCustomerKey() (customerKey Key, err error) {
	customerUniqueID := k.generateCustomerUniqueID()

	as := New()

	customer, err := as.GetCustomerByCustomerUniqueID(customerUniqueID)
	if err == sql.ErrNoRows {
		entry, err := as.GetRegistrationKey()
		if err != nil {
			return customerKey, err
		}
		customerKey.ApiKey = entry.Key

		return customerKey, nil
	} else if err != nil {
		return
	}

	entry, err := as.GetKeyByCustomerID(customer.ID)
	if err != nil {
		return
	}

	customerKey.CustomerUniqueID = k.generateCustomerUniqueID()
	customerKey.ApiKey = entry.Key

	return
}

func (k *Key) GenerateCustomerKey() (customerKey Key, err error) {
	customerUniqueID := k.generateCustomerUniqueID()

	as := New()

	_, err = as.GetCustomerByCustomerUniqueID(customerUniqueID)
	if err != nil && err != sql.ErrNoRows {
		return
	}

	if err == nil {
		err = ErrCustomerExists
		return
	}

	tx, err := DB.Begin()
	if err != nil {
		return
	}

	customer := &models.Customer{
		ProgramID:             k.programID,
		ProgramCustomerID:     k.programCustomerID,
		ProgramCustomerMobile: k.programCustomerMobile,
		CustomerUniqueID:      customerUniqueID,
	}
	err = tx.Model(customer).Insert()
	if err != nil {
		tx.Rollback()
		return
	}

	entry := &models.APIKey{
		CustomerID:  customer.ID,
		Key:         k.generateAPIKey(),
		DateCreated: time.Now().UTC(),
	}
	err = tx.Model(entry).Insert()
	if err != nil {
		tx.Rollback()
		return
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return
	}

	customerKey.ApiKey = entry.Key
	customerKey.CustomerUniqueID = customer.CustomerUniqueID

	return
}

func (k *Key) generateCustomerUniqueID() string {
	str := strconv.Itoa(k.programID) + "_" + strconv.Itoa(k.programCustomerID) + "_" + k.programCustomerMobile

	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func (k *Key) generateAPIKey() string {
	return random.String(32)
}
