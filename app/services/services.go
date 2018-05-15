package services

import (
	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/database"
	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/services/api"
	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/services/customer"

	dbx "github.com/go-ozzo/ozzo-dbx"
)

var db *dbx.DB

type Services struct{}

func New(DB *database.Database) error {
	db = DB.GetInstance()

	api.DB = db
	customer.DB = db

	return nil
}
