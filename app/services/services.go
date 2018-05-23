package services

import (
	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/database"
	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/services/api"
	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/services/customer"

	dbx "github.com/go-ozzo/ozzo-dbx"
)

// db is the database handler
var db *dbx.DB

// Services boots application-specific services
type Services struct{}

// New starts the service setup process
func New(DB *database.Database) error {
	db = DB.GetInstance()

	// Attach the database handler to service
	api.DB = db
	customer.DB = db

	return nil
}
