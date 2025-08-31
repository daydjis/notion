package service

import (
	"errors"
	"todo-api/internal/model"
	"todo-api/internal/repository"
)

type InvestmentService interface {
	GetAllAccount(userID uint) (*model.AccountsSummary, error)
	GetAccountById(userID uint, id uint) (*model.Bank, error)
	CreateAccount(task model.BankInput) (*model.Bank, error)
	UpdateAccount(task model.BankInput) (*model.Bank, error)
	DeleteAccount(id uint, userID uint) error
}

type investmentService struct {
	repo repository.InvestmentRepository
}

func NewInvestmentService(repo repository.InvestmentRepository) InvestmentService {
	return &investmentService{repo: repo}
}

func (i investmentService) GetAllAccount(userID uint) (*model.AccountsSummary, error) {
	accounts, err := i.repo.GetAll(userID)
	if err != nil {
		return nil, err
	}
	var totalIncome float32
	for i := range accounts {

		income := accounts[i].InvoiceAmount * accounts[i].Percentage / 100
		accounts[i].PercentageIncome = income
		totalIncome += income
	}

	accountsSummary := &model.AccountsSummary{
		TotalIncome: totalIncome,
		Accounts:    accounts,
	}

	return accountsSummary, nil
}
func (i investmentService) GetAccountById(userID uint, accId uint) (*model.Bank, error) {
	account, err := i.repo.GetByID(accId)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (i investmentService) CreateAccount(account model.BankInput) (*model.Bank, error) {

	tx := &model.Bank{
		UserID:        account.UserID,
		Name:          account.Name,
		Currency:      account.Currency,
		BankName:      account.BankName,
		InvoiceAmount: account.InvoiceAmount,
		Percentage:    account.Percentage,
		Type:          account.Type,
	}
	return i.repo.Create(tx)
}

func (i investmentService) UpdateAccount(account model.BankInput) (*model.Bank, error) {
	tx := &model.Bank{
		UserID:        account.UserID,
		Name:          account.Name,
		Currency:      account.Currency,
		BankName:      account.BankName,
		InvoiceAmount: account.InvoiceAmount,
		Percentage:    account.Percentage,
		Type:          account.Type,
	}
	return i.repo.Update(tx)
}

func (i investmentService) DeleteAccount(id uint, userID uint) error {

	acc, err := i.repo.GetByID(id)
	if err != nil {
		return errors.New("forbidden: not found user")
	}
	if acc.UserID != userID {
		errors.New("forbidden: investment user ID not match")
	}

	err = i.repo.Delete(id)
	if err != nil {
		return errors.New("forbidden: investment ID not match")
	}

	return nil
}
