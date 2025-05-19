package service

import (
	"todo-api/internal/model"
	"todo-api/internal/repository"
)

type TaskService interface {
	GetAll() ([]model.Task, error)
	Create(task model.Task) (model.Task, error)
	Delete(id int) error
}

type taskService struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) TaskService {
	return &taskService{repo: repo}
}

func (s *taskService) GetAll() ([]model.Task, error) {
	return s.repo.GetAll()
}

func (s *taskService) Create(task model.Task) (model.Task, error) {
	return s.repo.Create(task)
}

func (s *taskService) Delete(id int) error {
	return s.repo.Delete(id)
}
