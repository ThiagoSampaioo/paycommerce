package models

import (
	"time"
	"gorm.io/gorm"
)

type Sale struct {
	gorm.Model
	UserID      uint          `json:"user_id"`
	User        User          `gorm:"foreignKey:UserID"`
	Amount      float64       `json:"amount"`
	Status      string        `json:"status"`
	ApiKey      string        `json:"api_key"`
	ReferenceID string        `json:"reference_id" gorm:"type:varchar(191);uniqueIndex"`
	CallbackURL string        `json:"callback_url"`
	PaidAt      *time.Time    `json:"paid_at"`
	Items       []SaleItem    `json:"items"` // relacionamento
}

type SaleItem struct {
	gorm.Model
	SaleID   uint    `json:"sale_id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}
