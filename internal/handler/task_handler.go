package handler

import (
	"net/http"
	"strconv"
	"todo-api/internal/model"
	"todo-api/internal/service"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	service service.TaskService
}

func NewTaskHandler(s service.TaskService) *TaskHandler {
	return &TaskHandler{service: s}
}

// RegisterRoutes добавляет маршруты для задач
func (h *TaskHandler) RegisterRoutes(rg *gin.RouterGroup) {
	tasks := rg.Group("/tasks")
	{
		tasks.POST("/", h.CreateTask)
		tasks.GET("/", h.GetAllTasks)
		tasks.DELETE("/:id", h.DeleteTask)
	}
}

// CreateTask создает новую задачу для авторизованного пользователя
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var input model.CreateTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims := jwt.ExtractClaims(c)
	uid, ok := claims["id"]
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := uint(uid.(float64))

	task := model.Task{
		UserID: userID,
		Name:   input.Name,
		// добавь остальные поля, если они есть в Task
	}
	createdTask, err := h.service.CreateTask(task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, createdTask)
}

// GetAllTasks возвращает все задачи пользователя
func (h *TaskHandler) GetAllTasks(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	uid, ok := claims["id"]
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := uint(uid.(float64))

	tasks, err := h.service.GetAllTasks(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

// DeleteTask удаляет задачу по id (только свою!)
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	uid, ok := claims["id"]
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := uint(uid.(float64))

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.service.DeleteTask(uint(id), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
