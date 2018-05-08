package repositories

import (
	"github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/database"

	dbx "github.com/go-ozzo/ozzo-dbx"
)

var db *dbx.DB

type Repositories struct{}

func New(DB *database.Database) (*Repositories, error) {
	db = DB.GetInstance()

	r := new(Repositories)
	return r, nil
}
