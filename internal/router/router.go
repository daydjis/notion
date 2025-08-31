package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	auth "todo-api/internal"
	"todo-api/internal/handler"
	"todo-api/internal/service"
)

func NewRouter(
	userHandler *handler.UserHandler,
	taskHandler *handler.TaskHandler,
	transactionHandler *handler.TransactionHandler,
	investmentHandler *handler.InvestmentHandler,
	userSvc service.UserService,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// JWT Middleware
	authMiddleware, _ := auth.JwtMiddleware(userSvc)

	router.POST("/register", userHandler.RegisterHandler)
	router.POST("/login", authMiddleware.LoginHandler)

	api := router.Group("/")
	api.Use(authMiddleware.MiddlewareFunc())
	{
		investmentHandler.RegisterRoutes(api)
		userHandler.RegisterRoutes(api)
		taskHandler.RegisterRoutes(api)
		transactionHandler.RegisterRoutes(api)
	}

	return router
}
