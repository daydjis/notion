package service

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"todo-api/internal/model"
	"todo-api/internal/repository"
)

type UserService interface {
	CreateUser(input *model.RegisterInput) (*model.User, error)
	GetUser(id uint) (*model.User, error)
	GetAllUsers() ([]model.User, error)
	UpdateUser(user *model.User) (*model.User, error)
	DeleteUser(id uint) error
	Register(input model.RegisterInput) (*model.User, error)
	AuthenticateUser(name, password string) (*model.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

// CreateUser создаёт пользователя (используется при регистрации)
func (s *userService) CreateUser(input *model.RegisterInput) (*model.User, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Name:       input.Name,
		SecondName: input.SecondName,
		Age:        input.Age,
		Password:   string(hashedPass),
	}

	return s.repo.Create(user)
}

// GetUser возвращает пользователя по ID
func (s *userService) GetUser(id uint) (*model.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// GetAllUsers возвращает всех пользователей
func (s *userService) GetAllUsers() ([]model.User, error) {
	users, err := s.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	return users, nil
}

// UpdateUser обновляет данные пользователя
func (s *userService) UpdateUser(user *model.User) (*model.User, error) {
	if err := s.repo.Update(user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	return user, nil
}

// DeleteUser удаляет пользователя по ID
func (s *userService) DeleteUser(id uint) error {
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// Register регистрирует нового пользователя
func (s *userService) Register(input model.RegisterInput) (*model.User, error) {
	return s.CreateUser(&input)
}

// AuthenticateUser проверяет имя и пароль пользователя
func (s *userService) AuthenticateUser(name, password string) (*model.User, error) {
	user, err := s.repo.GetByName(name)
	fmt.Println(name, password, user)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}
	return user, nil
}
