package service

import (
	"time"
	"todo-api/internal/model"
	"todo-api/internal/repository"
)

type TransactionService interface {
	CreateTransaction(input model.TransactionInput, userID uint) (*model.Transaction, error)
	GetTransaction(id uint) (*model.Transaction, error)
	GetAllTransactions(userID uint) ([]model.Transaction, error)
	DeleteTransaction(id uint) error
}

type transactionService struct {
	repo repository.TransactionRepository
}

func NewTransactionService(repo repository.TransactionRepository) TransactionService {
	return &transactionService{repo: repo}
}

func (s *transactionService) CreateTransaction(input model.TransactionInput, userID uint) (*model.Transaction, error) {
	date, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		return nil, err
	}
	tx := &model.Transaction{
		UserID:    userID,
		AccountID: input.AccountID,
		Amount:    input.Amount,
		Type:      input.Type,
		Category:  input.Category,
		Date:      date,
		Comment:   input.Comment,
	}
	return s.repo.Create(tx)
}

func (s *transactionService) GetTransaction(id uint) (*model.Transaction, error) {
	return s.repo.GetByID(id)
}

func (s *transactionService) GetAllTransactions(userID uint) ([]model.Transaction, error) {
	return s.repo.GetAllByUser(userID)
}

func (s *transactionService) DeleteTransaction(id uint) error {
	return s.repo.Delete(id)
}
