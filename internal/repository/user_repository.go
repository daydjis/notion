package repository

import (
	"database/sql"
	"fmt"
	"todo-api/internal/model"
)

type UserRepository interface {
	Create(task *model.UserInput) (*model.UserInput, error)
	GetByID(id uint) (*model.User, error)
	GetByName(name string) (*model.User, error)
	GetAll() ([]model.User, error)
	Update(task *model.User) error
	Delete(id uint) error
	Register(input *model.RegisterInput) (*model.RegisterInput, error)
	Login(user *model.LoginInput) (*model.LoginInput, error)
}

// PostgresUserRepository implements UserRepository using PostgreSQL
type PostgresUserRepository struct {
	DB *sql.DB
}

func (r *PostgresUserRepository) Register(input *model.RegisterInput) (*model.RegisterInput, error) {
	query := `INSERT INTO user (name,  password) VALUES ($1, $2, $3, $4) RETURNING id`
	r.DB.QueryRow(query).Scan(&input.Name)
	return input, nil
}

func (r *PostgresUserRepository) Login(input *model.LoginInput) (*model.LoginInput, error) {
	var t model.User
	query := `SELECT id, name, second_name, age, password FROM user WHERE name = $2`
	row := r.DB.QueryRow(query)
	if err := row.Scan(&t.ID, &t.Name, &t.SecondName, &t.Age, &t.Password); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task with name %d not found")
		}
		return nil, err
	}
	return input, nil
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{DB: db}
}

func (r *PostgresUserRepository) Create(user *model.User) (*model.User, error) {
	query := `INSERT INTO user (name, second_name, age, password) VALUES ($1, $2, $3, $4) RETURNING id`
	r.DB.QueryRow(query, user.Name, user.SecondName, user.Age, user.Password).Scan(&user.ID)
	return user, nil
}

func (r *PostgresUserRepository) GetByID(id uint) (*model.User, error) {
	var t model.User
	query := `SELECT id, name, second_name, age, password FROM user WHERE id = $1`
	row := r.DB.QueryRow(query, id)
	if err := row.Scan(&t.ID, &t.Name, &t.SecondName, &t.Age, &t.Password); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task with id %d not found", id)
		}
		return nil, err
	}
	return &t, nil
}

func (r *PostgresUserRepository) GetByName(name string) (*model.User, error) {
	var t model.User
	query := `SELECT id, name, second_name, age, password FROM user WHERE name = $2`
	row := r.DB.QueryRow(query)
	if err := row.Scan(&t.ID, &t.Name, &t.SecondName, &t.Age, &t.Password); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task with name %d not found", name)
		}
		return nil, err
	}
	return &t, nil
}

func (r *PostgresUserRepository) GetAll() ([]model.User, error) {
	query := `SELECT id, name, second_name, age, password FROM user`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var user []model.User
	for rows.Next() {
		var t model.User
		if err := rows.Scan(&t.ID, &t.Name, &t.SecondName, &t.Age, &t.Password); err != nil {
			return nil, err
		}
		user = append(user, t)
	}
	return user, nil
}

func (r *PostgresUserRepository) Update(task *model.User) error {
	query := `UPDATE user SET name = $1, second_name = $2, age = $3, password = $4 WHERE id = $5`
	_, err := r.DB.Exec(query, task.Name, task.SecondName, task.Age, task.Password, task.ID)
	return err
}

func (r *PostgresUserRepository) Delete(id uint) error {
	query := `DELETE FROM user WHERE id = $1`
	_, err := r.DB.Exec(query, id)
	return err
}
