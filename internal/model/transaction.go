package model

import (
	"gorm.io/gorm"
	"time"
)

// Transaction — финансовая операция (расход, доход, перевод и т.д.)
type Transaction struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `json:"user_id"` // Foreign key
	AccountID uint           `json:"account_id"`
	Amount    float64        `json:"amount"`
	Type      string         `json:"type"`     // "income", "expense", "transfer", "invest" и т.д.
	Category  string         `json:"category"` // или CategoryID, если у тебя категории в отдельной таблице
	Date      time.Time      `json:"date"`
	Comment   *string        `json:"comment,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// input для создания транзакции
type TransactionInput struct {
	AccountID uint    `json:"account_id"`
	Amount    float64 `json:"amount"`
	Type      string  `json:"type"`     // "income", "expense", etc.
	Category  string  `json:"category"` // или CategoryID, если нужна связь
	Date      string  `json:"date"`     // ISO формат, парси к time.Time
	Comment   *string `json:"comment,omitempty"`
}
