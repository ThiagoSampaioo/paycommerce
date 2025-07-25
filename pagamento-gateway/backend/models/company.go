package models

import (
	"time"

	"gorm.io/gorm"
)

type Company struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `json:"name"`
	Email     string         `gorm:"unique" json:"email"`
	Password  string         `json:"password"`
	ApiKey    string         `gorm:"unique" json:"api_key"`
	IsActive  bool           `json:"is_active"`
	Balance   float64        `json:"balance" gorm:"default:0"`
	CallbackURL string       `json:"callback_url"`

	// Dados bancários
	BankCode      string `json:"bank_code"`
	BankName      string `json:"bank_name"`
	AgencyAccount string `json:"agency_account"`
	NumberAccount string `json:"number_account"`
	TypeAccount   string `json:"type_account"` // "CONTA_CORRENTE", "POUPANCA"

	// Pix
	TypeKeyPix string `json:"type_key_pix"`   // "CPF", "CNPJ", "EMAIL", "TELEFONE", "CHAVE_ALEATORIA"
	KeyPix     string `json:"key_pix"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relação com pagamentos
	Payments []Payment `gorm:"foreignKey:CompanyID" json:"-"`
}

