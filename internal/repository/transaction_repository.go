package repository

import (
	"fmt"
	"gorm.io/gorm"
	"todo-api/internal/model"
)

type TransactionRepository interface {
	Create(tx *model.Transaction) (*model.Transaction, error)
	GetByID(id uint) (*model.Transaction, error)
	GetAllByUser(userID uint) ([]model.Transaction, error)
	Delete(id uint) error
}

type transactionRepo struct {
	db *gorm.DB
}

func NewTransactionRepo(db *gorm.DB) TransactionRepository {
	return &transactionRepo{db: db}
}

func (r *transactionRepo) Create(tx *model.Transaction) (*model.Transaction, error) {
	if err := r.db.Create(tx).Error; err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}
	return tx, nil
}

func (r *transactionRepo) GetByID(id uint) (*model.Transaction, error) {
	var tx model.Transaction
	if err := r.db.First(&tx, id).Error; err != nil {
		return nil, fmt.Errorf("failed to get transaction by id: %w", err)
	}
	return &tx, nil
}

func (r *transactionRepo) GetAllByUser(userID uint) ([]model.Transaction, error) {
	var txs []model.Transaction
	if err := r.db.Where("user_id = ?", userID).Order("date desc").Find(&txs).Error; err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}
	return txs, nil
}

func (r *transactionRepo) Delete(id uint) error {
	if err := r.db.Delete(&model.Transaction{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete transaction: %w", err)
	}
	return nil
}
