package controllers

import "github.com/epointpayment/customerprofilingengine-demo-classifier-api/app/database"

type Controllers struct {
	DB *database.Database
}

func NewControllers(db *database.Database) *Controllers {
	c := &Controllers{
		DB: db,
	}

	return c
}
