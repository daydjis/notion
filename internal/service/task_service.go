package service

import (
	"todo-api/internal/model"
	"todo-api/internal/repository"
)

type TaskService interface {
	GetAllTasks(userID uint) ([]model.Task, error)
	CreateTask(task model.Task) (model.Task, error)
	DeleteTask(id uint, userID uint) error
}

type taskService struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) TaskService {
	return &taskService{repo: repo}
}

// Получить все задачи только для пользователя
func (s *taskService) GetAllTasks(userID uint) ([]model.Task, error) {
	return s.repo.GetAllByUser(userID)
}

// Создать задачу для пользователя (userID прокидывается из handler)
func (s *taskService) CreateTask(task model.Task) (model.Task, error) {
	return s.repo.Create(task)
}

// Удалить задачу только если она принадлежит userID
func (s *taskService) DeleteTask(id uint, userID uint) error {
	return s.repo.Delete(id, userID)
}
