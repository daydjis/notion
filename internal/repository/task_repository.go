package repository

import (
	"database/sql"
	"todo-api/internal/model"
)

type TaskRepository interface {
	GetAll() ([]model.Task, error)
	Create(task model.Task) (model.Task, error)
	Delete(id int) error
}

type taskRepo struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) TaskRepository {
	return &taskRepo{db: db}
}

func (r *taskRepo) GetAll() ([]model.Task, error) {
	rows, err := r.db.Query("SELECT id, name FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []model.Task
	for rows.Next() {
		var t model.Task
		if err := rows.Scan(&t.ID, &t.Name); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (r *taskRepo) Create(task model.Task) (model.Task, error) {
	err := r.db.QueryRow("INSERT INTO tasks (name) VALUES ($1) RETURNING id", task.Name).Scan(&task.ID)
	return task, err
}

func (r *taskRepo) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM tasks WHERE id = $1", id)
	return err
}
