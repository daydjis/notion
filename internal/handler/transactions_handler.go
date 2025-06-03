package handler

import (
	"net/http"
	"strconv"
	"todo-api/internal/model"
	"todo-api/internal/service"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	service service.TransactionService
}

func NewTransactionHandler(s service.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: s}
}

// RegisterRoutes добавляет маршруты для транзакций
func (h *TransactionHandler) RegisterRoutes(rg *gin.RouterGroup) {
	tx := rg.Group("/transactions")
	{
		tx.POST("/", h.CreateTransaction)
		tx.GET("/", h.GetAllTransactions)
		tx.GET("/:id", h.GetTransaction)
		tx.DELETE("/:id", h.DeleteTransaction)
	}
}

// CreateTransaction — создать новую транзакцию
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var input model.TransactionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Получаем userID из JWT claims
	claims := jwt.ExtractClaims(c)
	uid, ok := claims["id"]
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := uint(uid.(float64))

	tx, err := h.service.CreateTransaction(input, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, tx)
}

// GetAllTransactions — все транзакции пользователя
func (h *TransactionHandler) GetAllTransactions(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	uid, ok := claims["id"]
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := uint(uid.(float64))

	txs, err := h.service.GetAllTransactions(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, txs)
}

// GetTransaction — получить транзакцию по id (и проверить принадлежит ли она userID)
func (h *TransactionHandler) GetTransaction(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	tx, err := h.service.GetTransaction(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tx)
}

// DeleteTransaction — удалить транзакцию по id (только свою)
func (h *TransactionHandler) DeleteTransaction(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
