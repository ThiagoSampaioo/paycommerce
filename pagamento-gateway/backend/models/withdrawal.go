package models

import "time"

type Withdrawal struct {
	ID        uint      `json:"id"`
	CompanyID uint      `json:"company_id"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}
