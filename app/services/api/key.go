package api

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"strconv"
	"time"

	"github.com/epointpayment/mloc-cpe/app/models"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/juju/errors"
	"github.com/labstack/gommon/random"
)

// Key manages the customer API
type Key struct {
	programID             int
	programCustomerID     int
	programCustomerMobile string
	CustomerUniqueID      string `json:"cust_unique_id"`
	ApiKey                string `json:"api_key"`
}

// Validate checks if the values in the struct are valid
func (k Key) Validate() error {
	return validation.ValidateStruct(&k,
		validation.Field(&k.programID, validation.Required),
		validation.Field(&k.programCustomerID, validation.Required),
		validation.Field(&k.programCustomerMobile, validation.Required),
	)
}

// NewKey creates an instance of the customer key service
func NewKey(programID int, programCustomerID int, programCustomerMobile string) (k *Key, err error) {
	k = &Key{
		programID:             programID,
		programCustomerID:     programCustomerID,
		programCustomerMobile: programCustomerMobile,
	}

	err = k.Validate()
	if err != nil {
		err = errors.Trace(err)
		return
	}

	return
}

// GetCustomerKey gets a customer and associated customer API key
func (k *Key) GetCustomerKey() (customerKey Key, err error) {
	customerUniqueID := k.generateCustomerUniqueID()

	as := New()

	customer, err := as.GetCustomerByCustomerUniqueID(customerUniqueID)
	if errors.Cause(err) == sql.ErrNoRows {
		entry, err := as.GetRegistrationKey()
		if err != nil {
			err = errors.Trace(err)
			return customerKey, err
		}
		customerKey.ApiKey = entry.Key

		return customerKey, nil
	} else if err != nil {
		err = errors.Trace(err)
		return
	}

	entry, err := as.GetKeyByCustomerID(customer.ID)
	if err != nil {
		err = errors.Trace(err)
		return
	}

	customerKey.CustomerUniqueID = k.generateCustomerUniqueID()
	customerKey.ApiKey = entry.Key

	return
}

// GenerateCustomerKey generates a customer and associated customer API key
func (k *Key) GenerateCustomerKey() (customerKey Key, err error) {
	customerUniqueID := k.generateCustomerUniqueID()

	as := New()

	_, err = as.GetCustomerByCustomerUniqueID(customerUniqueID)
	if err != nil && errors.Cause(err) != sql.ErrNoRows {
		err = errors.Trace(err)
		return
	}

	if err == nil {
		err = errors.Wrap(err, ErrCustomerExists)
		return
	}

	tx, err := DB.Begin()
	if err != nil {
		err = errors.Trace(err)
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
		err = errors.Trace(err)
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
		err = errors.Trace(err)
		tx.Rollback()
		return
	}

	err = tx.Commit()
	if err != nil {
		err = errors.Trace(err)
		tx.Rollback()
		return
	}

	customerKey.ApiKey = entry.Key
	customerKey.CustomerUniqueID = customer.CustomerUniqueID

	return
}

// generateCustomerUniqueID generates an MD5 hash from customer program information
func (k *Key) generateCustomerUniqueID() string {
	str := strconv.Itoa(k.programID) + "_" + strconv.Itoa(k.programCustomerID) + "_" + k.programCustomerMobile

	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// generateAPIKey generates a random string of 32 characters
func (k *Key) generateAPIKey() string {
	return random.String(32)
}
