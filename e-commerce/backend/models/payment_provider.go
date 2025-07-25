package models

import "gorm.io/gorm"

type PaymentProvider struct {
	gorm.Model
	Name      string `json:"name" gorm:"unique"`
	ApiKey   string `json:"api_key"`
	PaymentURL string `json:"payment_url"`
	CancelationURL string `json:"cancellation_url"`
	Active   bool   `json:"active"`
}
