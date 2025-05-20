package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"todo-api/internal/model"
	"todo-api/internal/service"
)

type TaskHandler struct {
	service service.TaskService
}

func NewTaskHandler(s service.TaskService) *TaskHandler {
	return &TaskHandler{service: s}
}

func (h *TaskHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/tasks", h.getTasks)
	r.POST("/tasks", h.createTask)
	r.DELETE("/tasks/:id", h.deleteTask)
}

func (h *TaskHandler) getTasks(c *gin.Context) {
	tasks, err := h.service.GetAllTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch tasks"})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) createTask(c *gin.Context) {
	var input model.CreateTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	created, err := h.service.CreateTask(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create task"})
		return
	}
	c.JSON(http.StatusCreated, created)
}

func (h *TaskHandler) deleteTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.service.DeleteTask(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete task"})
		return
	}
	c.Status(http.StatusNoContent)
}
