package handler

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"todo-api/internal/model"
	"todo-api/internal/service"
)

type UserHandler struct {
	Service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{Service: service}
}

func (h *UserHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.DELETE("/user/:id", h.deleteUser)
	r.GET("/users", h.getAllUsers)
	r.GET("/user/:id", h.getUserInfo)
}

// registerUser регистрирует нового пользователя (открытый маршрут)
func (h *UserHandler) RegisterHandler(c *gin.Context) {
	var input model.RegisterInput
	claims := jwt.ExtractClaims(c)
	_, ok := claims["id"]
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.Service.Register(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// getAllUsers возвращает список всех пользователей
func (h *UserHandler) getAllUsers(c *gin.Context) {
	users, err := h.Service.GetAllUsers()
	claims := jwt.ExtractClaims(c)
	_, ok := claims["id"]
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// getUserInfo возвращает данные пользователя по ID из пути
func (h *UserHandler) getUserInfo(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	claims := jwt.ExtractClaims(c)
	_, ok := claims["id"]
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	user, err := h.Service.GetUser(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// deleteUser удаляет пользователя по ID
func (h *UserHandler) deleteUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	claims := jwt.ExtractClaims(c)
	_, ok := claims["id"]
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	if err := h.Service.DeleteUser(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
