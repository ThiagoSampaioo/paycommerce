package models

import (
	"time"
	"gorm.io/gorm"
)

type Payment struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CompanyID   uint           `json:"company_id"`
	Company     Company        `gorm:"foreignKey:CompanyID" json:"-"`
	Amount      float64        `json:"amount"`
	Status      string         `gorm:"default:PENDING" json:"status"` // Use constantes
	CallbackURL string         `json:"callback_url"`
	ReferenceID string         `gorm:"index" json:"reference_id"`     // indexável ou único se necessário
	ApiKey       string    `json:"api_key"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
