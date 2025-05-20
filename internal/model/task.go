package model

// Task представляет структуру задачи
type Task struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// CreateTaskInput структура для создания задачи
type CreateTaskInput struct {
	Name string `json:"name" binding:"required"`
}

// UpdateTaskInput структура для обновления задачи
type UpdateTaskInput struct {
	Name *string `json:"name,omitempty"`
}
