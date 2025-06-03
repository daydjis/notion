package repository

import (
	"fmt"
	"gorm.io/gorm"
	"todo-api/internal/model"
)

// TaskRepository определяет методы доступа к задачам пользователя
type TaskRepository interface {
	GetAllByUser(userID uint) ([]model.Task, error)
	Create(task model.Task) (model.Task, error)
	Delete(id uint, userID uint) error
}

type gormTaskRepo struct {
	db *gorm.DB
}

func NewTaskRepo(db *gorm.DB) TaskRepository {
	return &gormTaskRepo{db: db}
}

// GetAllByUser возвращает все задачи для конкретного пользователя
func (r *gormTaskRepo) GetAllByUser(userID uint) ([]model.Task, error) {
	var tasks []model.Task
	if err := r.db.Where("user_id = ?", userID).Find(&tasks).Error; err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}
	return tasks, nil
}

// Create создаёт новую задачу
func (r *gormTaskRepo) Create(task model.Task) (model.Task, error) {
	if err := r.db.Create(&task).Error; err != nil {
		return task, fmt.Errorf("failed to create task: %w", err)
	}
	return task, nil
}

// Delete удаляет задачу только если она принадлежит пользователю
func (r *gormTaskRepo) Delete(id uint, userID uint) error {
	res := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Task{})
	if res.Error != nil {
		return fmt.Errorf("failed to delete task: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("task not found or does not belong to user")
	}
	return nil
}
