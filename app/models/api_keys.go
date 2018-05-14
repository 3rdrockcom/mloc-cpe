package models

import "time"

type APIKeys []APIKey

type APIKey struct {
	ID          int
	CustomerID  int
	Key         string
	DateCreated time.Time
}

func (a APIKey) TableName() string {
	return "api_keys"
}
