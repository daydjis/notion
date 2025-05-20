package service

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"todo-api/internal/model"
	"todo-api/internal/repository"
)

type UserService interface {
	CreateUser(input *model.UserInput) (*model.UserInput, error)
	GetUser(id uint) (*model.User, error)
	GetAllUsers() ([]model.User, error)
	UpdateUser(input *model.User) (*model.User, error)
	DeleteUser(id uint) error
	Register(input model.RegisterInput) (*model.RegisterInput, error)
	Login(user model.LoginInput) (*model.User, error)
}

type userService struct {
	Repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{Repo: repo}
}

func (s *userService) CreateUser(input *model.UserInput) (*model.UserInput, error) {
	if err, _ := s.Repo.Create(input); err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}
	return input, nil
}

func (s *userService) GetUser(id uint) (*model.User, error) {
	task, err := s.Repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	return task, nil
}

func (s *userService) GetAllUsers() ([]model.User, error) {
	tasks, err := s.Repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}
	return tasks, nil
}

func (s *userService) Register(input model.RegisterInput) (*model.RegisterInput, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	u := &model.RegisterInput{
		Name:     input.Name,
		Password: string(hashed),
	}
	res, err := s.Repo.Create((*model.UserInput)(u))
	return (*model.RegisterInput)(res), err
}

func (s *userService) Login(input model.LoginInput) (*model.User, error) {

	u, err := s.Repo.GetByName(input.Name)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(input.Password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}
	return u, nil
}

func (s *userService) UpdateUser(input *model.User) (*model.User, error) {
	if err := s.Repo.Update(input); err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}
	return input, nil
}

func (s *userService) DeleteUser(id uint) error {
	if err := s.Repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}
	return nil
}
