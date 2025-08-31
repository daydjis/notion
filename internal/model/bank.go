package model

import (
	"gorm.io/gorm"
	"time"
)

type Bank struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	UserID           uint           `json:"user_id"` // Foreign key
	Name             string         `json:"name"`
	BankName         string         `json:"bank_name"`
	InvoiceAmount    float32        `json:"invoice_amount"`
	PercentageIncome float32        `json:"percentage_income"`
	Percentage       float32        `json:"percentage"`
	Type             string         `json:"type"`
	Currency         string         `json:"currency"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

type BankInput struct {
	UserID        uint    `json:"user_id"` // Foreign key
	Name          string  `json:"name" validate:"required"`
	BankName      string  `json:"bank_name" validate:"required"`
	Currency      string  `json:"currency" validate:"required"`
	InvoiceAmount float32 `json:"invoice_amount" validate:"required"`
	Percentage    float32 `json:"percentage" validate:"required,gte=0,lte=100"`
	Type          string  `json:"type" validate:"required"`
}

// DTO
type AccountsSummary struct {
	Accounts    []Bank  `json:"accounts"`
	TotalIncome float32 `json:"totalIncome"`
}
