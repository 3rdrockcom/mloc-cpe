package database

import (
	"log"

	dbx "github.com/go-ozzo/ozzo-dbx"
	_ "github.com/go-sql-driver/mysql"
)

var db *dbx.DB

type Database struct {
	driverName string
	dsn        string
	DB         *dbx.DB
}

func NewDatabase(driverName, dsn string) *Database {
	return &Database{
		driverName: driverName,
		dsn:        dsn,
	}
}

func (d Database) Open() error {
	var err error

	db, err = dbx.Open(d.driverName, d.dsn)
	if err != nil {
		return err
	}
	db.LogFunc = log.Printf

	return nil
}

func (d Database) Close() error {
	return d.DB.Close()
}

func (d Database) GetInstance() *dbx.DB {
	return db
}

func Get() *dbx.DB {
	return db
}
