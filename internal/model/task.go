package model

type Task struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	UserID uint   `json:"user_id"` // Foreign key
	Name   string `json:"name"`
	// другие поля
}

// CreateTaskInput структура для создания задачи
type CreateTaskInput struct {
	Name string `json:"name" binding:"required"`
}

// UpdateTaskInput структура для обновления задачи
type UpdateTaskInput struct {
	Name *string `json:"name,omitempty"`
}
