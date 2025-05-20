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

func (h *UserHandler) RegisterRoutes(r *gin.Engine) {
	r.POST("/register", h.registerUser)
	r.POST("/login", h.loginUser)
	r.DELETE("/user/:id", h.deleteUser)
	r.GET("/users", h.getAllUsers)
	r.GET("/user/:id", h.getUserInfo)
}

// registerUser handles new user registration
func (h *UserHandler) registerUser(c *gin.Context) {
	var input model.RegisterInput
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

// loginUser handles user authentication and returns a token
func (h *UserHandler) loginUser(c *gin.Context) {
	var input model.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := h.Service.Login(input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// getAllUsers returns a list of all users
func (h *UserHandler) getAllUsers(c *gin.Context) {
	users, err := h.Service.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// getUserInfo returns details of a single user by ID
func (h *UserHandler) getUserInfo(c *gin.Context) {
	claims := jwt.ExtractClaims(c)

	id, err := strconv.Atoi(claims["id"].(string))
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

// deleteUser removes a user by ID
func (h *UserHandler) deleteUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
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
