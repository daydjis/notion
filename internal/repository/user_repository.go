package repository

import (
	"fmt"
	"todo-api/internal/model"

	"gorm.io/gorm"
)

// UserRepository описывает методы работы с пользователями в БД
type UserRepository interface {
	Create(user *model.User) (*model.User, error)
	GetByID(id uint) (*model.User, error)
	GetByName(name string) (*model.User, error)
	GetAll() ([]model.User, error)
	Update(user *model.User) error
	Delete(id uint) error
}

// GormUserRepository реализует интерфейс UserRepository
type GormUserRepository struct {
	DB *gorm.DB
}

// NewUserRepo создаёт новый экземпляр репозитория
func NewUserRepo(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{DB: db}
}

// Create создаёт нового пользователя
func (r *GormUserRepository) Create(user *model.User) (*model.User, error) {
	if err := r.DB.Create(user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return user, nil
}

// GetByID возвращает пользователя по ID
func (r *GormUserRepository) GetByID(id uint) (*model.User, error) {
	var user model.User
	if err := r.DB.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return &user, nil
}

// GetByName возвращает пользователя по имени
func (r *GormUserRepository) GetByName(name string) (*model.User, error) {
	var user model.User
	if err := r.DB.Where("name = ?", name).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user with name '%s' not found", name)
		}
		return nil, fmt.Errorf("failed to get user by name: %w", err)
	}
	return &user, nil
}

// GetAll возвращает список всех пользователей
func (r *GormUserRepository) GetAll() ([]model.User, error) {
	var users []model.User
	if err := r.DB.Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	return users, nil
}

// Update обновляет данные пользователя
func (r *GormUserRepository) Update(user *model.User) error {
	if err := r.DB.Save(user).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// Delete удаляет пользователя по ID
func (r *GormUserRepository) Delete(id uint) error {
	if err := r.DB.Delete(&model.User{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}
