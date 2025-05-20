package repository

import (
	"database/sql"
	"fmt"
	"todo-api/internal/model"
)

// TaskRepository определяет методы доступа к задачам
type TaskRepository interface {
	GetAll() ([]model.Task, error)
	Create(task model.Task) (model.Task, error)
	Delete(id uint) error
}

// taskRepo реализация TaskRepository
type taskRepo struct {
	db *sql.DB
}

// NewTaskRepository создаёт экземпляр репозитория задач
func NewTaskRepository(db *sql.DB) TaskRepository {
	return &taskRepo{db: db}
}

// GetAll возвращает все задачи из БД
func (r *taskRepo) GetAll() ([]model.Task, error) {
	rows, err := r.db.Query("SELECT id, name FROM tasks")
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}
	defer rows.Close()

	var tasks []model.Task
	for rows.Next() {
		var t model.Task
		if err := rows.Scan(&t.ID, &t.Name); err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

// Create создаёт новую задачу
func (r *taskRepo) Create(task model.Task) (model.Task, error) {
	query := "INSERT INTO tasks (name) VALUES ($1) RETURNING id"
	err := r.db.QueryRow(query, task.Name).Scan(&task.ID)
	if err != nil {
		return task, fmt.Errorf("failed to create task: %w", err)
	}
	return task, nil
}

// Delete удаляет задачу по ID
func (r *taskRepo) Delete(id uint) error {
	query := "DELETE FROM tasks WHERE id = $1"
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task with id %d: %w", id, err)
	}
	return nil
}
