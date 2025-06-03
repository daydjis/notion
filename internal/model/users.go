package model

import (
	"time"

	"gorm.io/gorm"
)

// User — основная структура пользователя для хранения в БД
type User struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	Name       string         `json:"name"`
	SecondName *string        `json:"second_name,omitempty"`
	Age        *int           `json:"age,omitempty"`
	Password   string         `json:"-"` // не отдаём пароль в json
	Tasks      []Task         `json:"tasks,omitempty" gorm:"foreignKey:UserID"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

// DTO для ответа пользователю (без пароля)
type UserResponse struct {
	ID         uint    `json:"id"`
	Name       string  `json:"name"`
	SecondName *string `json:"second_name,omitempty"`
	Age        *int    `json:"age,omitempty"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:         u.ID,
		Name:       u.Name,
		SecondName: u.SecondName,
		Age:        u.Age,
	}
}

// Category — пример категории задачи (необязательно, но часто удобно)
type Category struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Title string `json:"title"`
	Tasks []Task `json:"tasks,omitempty"`
}

// LoginInput — для логина
type LoginInput struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterInput — для регистрации
type RegisterInput struct {
	Name       string  `json:"name" binding:"required"`
	SecondName *string `json:"second_name,omitempty"`
	Age        *int    `json:"age,omitempty"`
	Password   string  `json:"password" binding:"required"`
}

// UserInput — для обновления данных пользователя
type UserInput struct {
	Name       string  `json:"name" binding:"required"`
	SecondName *string `json:"second_name,omitempty"`
	Age        *int    `json:"age,omitempty"`
}
