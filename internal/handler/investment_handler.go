package handler

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-playground/validator/v10"
	"net/http"
	"strconv"
	"todo-api/internal/model"
	"todo-api/internal/service"
)

type InvestmentHandler struct {
	service service.InvestmentService
}

func NewInvestmentHandler(service service.InvestmentService) *InvestmentHandler {
	return &InvestmentHandler{service: service}
}

func (h *InvestmentHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.DELETE("/investment/:id", h.deleteAccounts)
	r.GET("/investments", h.GetAllAccounts)
	r.PUT("/investment/update", h.UpdateAccount)
	r.POST("/investment", h.CreateAccount)
	r.GET("/investment/:id", h.GetAccountById)
}

var validate = validator.New()

// CreateTask создает новую задачу для авторизованного пользователя
func (h *InvestmentHandler) CreateAccount(c *gin.Context) {
	var input model.BankInput
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
	input.UserID = uint(uid.(float64))

	if err := validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdAccount, err := h.service.CreateAccount(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, createdAccount)
}

func (h *InvestmentHandler) UpdateAccount(c *gin.Context) {
	var input model.BankInput
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
	input.UserID = uint(uid.(float64))

	if err := validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdAccount, err := h.service.UpdateAccount(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, createdAccount)
}

func (h *InvestmentHandler) GetAllAccounts(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	uid, ok := claims["id"]
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := uint(uid.(float64))

	accounts, err := h.service.GetAllAccount(userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, accounts)

}
func (h *InvestmentHandler) GetAccountById(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	uid, ok := claims["id"]
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := uint(uid.(float64))
	idParam := c.Param("id")
	if idParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing account id"})
	}

	// 2. Преобразуем id в uint
	accID64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account id"})
		return
	}
	accID := uint(accID64)
	accounts, err := h.service.GetAccountById(userID, accID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, accounts)

}

func (h *InvestmentHandler) deleteAccounts(c *gin.Context) {

	claims := jwt.ExtractClaims(c)
	uid, ok := claims["id"]
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := uint(uid.(float64))

	idParam := c.Param("id")
	if idParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing account id"})
	}

	// 2. Преобразуем id в uint
	accID64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account id"})
		return
	}
	accID := uint(accID64)

	err = h.service.DeleteAccount(accID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nil)
}
