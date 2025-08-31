package repository

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"todo-api/internal/model"
)

type InvestmentRepository interface {
	Create(account *model.Bank) (*model.Bank, error)
	GetByID(id uint) (*model.Bank, error)
	GetByName(name string) (*model.Bank, error)
	GetAll(userId uint) ([]model.Bank, error)
	Update(account *model.Bank) (*model.Bank, error)
	Delete(id uint) error
}
type gormInvestmentRepo struct {
	db *gorm.DB
}

func NewInvestmentRepo(db *gorm.DB) InvestmentRepository {
	return &gormInvestmentRepo{db: db}
}

func (g gormInvestmentRepo) Create(account *model.Bank) (*model.Bank, error) {
	if err := g.db.Create(&account).Error; err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}
	return account, nil
}

func (g gormInvestmentRepo) GetByID(id uint) (*model.Bank, error) {
	var account model.Bank
	if err := g.db.First(&account, id).Error; err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}
	return &account, nil
}

func (g gormInvestmentRepo) GetByName(name string) (*model.Bank, error) {
	var account model.Bank
	if err := g.db.Where("name = ?", name).First(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("account with name '%s' not found", name)
		}
		return nil, fmt.Errorf("failed to get account by name: %w", err)
	}
	return &account, nil
}

func (g gormInvestmentRepo) GetAll(userId uint) ([]model.Bank, error) {
	var accounts []model.Bank
	if err := g.db.Where("user_id = ?", userId).Find(&accounts).Error; err != nil {
		return nil, fmt.Errorf("failed to get accounts: %w", err)
	}
	return accounts, nil
}

func (g gormInvestmentRepo) Update(accountInput *model.Bank) (*model.Bank, error) {
	var account model.Bank
	if err := g.db.First(&account, account.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("account with id %d not found", account.ID)
		}
		return nil, fmt.Errorf("failed to get account with id %d: %w", account.ID, err)
	}

	return &account, nil
}

func (g gormInvestmentRepo) Delete(id uint) error {
	var account model.Bank
	if err := g.db.First(&account, id).Error; err != nil {
		return fmt.Errorf("failed to get account with id %d: %w", id, err)
	}
	return g.db.Delete(&account).Error
}
