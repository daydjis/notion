package router

import (
	"github.com/gin-gonic/gin"
	auth "todo-api/internal"
	"todo-api/internal/handler"
	"todo-api/internal/service"
)

func NewRouter(
	userHandler *handler.UserHandler,
	taskHandler *handler.TaskHandler,
	transactionHandler *handler.TransactionHandler,
	userSvc service.UserService,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// JWT Middleware
	authMiddleware, _ := auth.JwtMiddleware(userSvc)

	router.POST("/register", userHandler.RegisterHandler)
	router.POST("/login", authMiddleware.LoginHandler)

	api := router.Group("/")
	api.Use(authMiddleware.MiddlewareFunc())
	{
		userHandler.RegisterRoutes(api)
		taskHandler.RegisterRoutes(api)
		transactionHandler.RegisterRoutes(api)
	}

	return router
}
