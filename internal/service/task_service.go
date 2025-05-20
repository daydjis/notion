package service

import (
	"todo-api/internal/model"
	"todo-api/internal/repository"
)

type TaskService interface {
	GetAllTasks() ([]model.Task, error)
	CreateTask(input model.CreateTaskInput) (model.Task, error)
	DeleteTask(id uint) error
}

type taskService struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) TaskService {
	return &taskService{repo: repo}
}

func (s *taskService) GetAllTasks() ([]model.Task, error) {
	return s.repo.GetAll()
}

func (s *taskService) CreateTask(input model.CreateTaskInput) (model.Task, error) {
	task := model.Task{
		Name: input.Name,
	}
	return s.repo.Create(task)
}

func (s *taskService) DeleteTask(id uint) error {
	return s.repo.Delete(id)
}
