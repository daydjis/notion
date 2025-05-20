package repository

import (
	"database/sql"
	"fmt"
	"todo-api/internal/model"
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

// PostgresUserRepository реализует интерфейс UserRepository
type PostgresUserRepository struct {
	DB *sql.DB
}

// NewPostgresUserRepository создаёт новый экземпляр репозитория
func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{DB: db}
}

// Create создаёт нового пользователя
func (r *PostgresUserRepository) Create(user *model.User) (*model.User, error) {
	query := `INSERT INTO users (name, second_name, age, password) VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.DB.QueryRow(query, user.Name, user.SecondName, user.Age, user.Password).Scan(&user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return user, nil
}

// GetByID возвращает пользователя по ID
func (r *PostgresUserRepository) GetByID(id uint) (*model.User, error) {
	var user model.User
	query := `SELECT id, name, second_name, age, password FROM users WHERE id = $1`
	row := r.DB.QueryRow(query, id)
	err := row.Scan(&user.ID, &user.Name, &user.SecondName, &user.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return &user, nil
}

// GetByName возвращает пользователя по имени
func (r *PostgresUserRepository) GetByName(name string) (*model.User, error) {
	var user model.User
	query := `SELECT id, name, second_name, age, password FROM users WHERE name = $1`
	row := r.DB.QueryRow(query, name)
	err := row.Scan(&user.ID, &user.Name, &user.SecondName, &user.Age, &user.Password)
	fmt.Println(&row, err, &user)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with name '%s' not found", name)
		}
		return nil, fmt.Errorf("failed to get user by name: %w", err)
	}
	return &user, nil
}

// GetAll возвращает список всех пользователей
func (r *PostgresUserRepository) GetAll() ([]model.User, error) {
	query := `SELECT id, name, second_name, age FROM users`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID, &user.Name, &user.SecondName, &user.Age); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}
	return users, nil
}

// Update обновляет данные пользователя
func (r *PostgresUserRepository) Update(user *model.User) error {
	query := `UPDATE users SET name = $1, second_name = $2, age = $3, password = $4 WHERE id = $5`
	_, err := r.DB.Exec(query, user.Name, user.SecondName, user.Age, user.Password, user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// Delete удаляет пользователя по ID
func (r *PostgresUserRepository) Delete(id uint) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}
