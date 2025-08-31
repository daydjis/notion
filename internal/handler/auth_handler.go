package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"todo-api/internal/model"
)

// registerUser регистрирует нового пользователя (открытый маршрут)
func (h *UserHandler) RegisterHandler(c *gin.Context) {
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
