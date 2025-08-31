// cmd/main.go
package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"todo-api/internal/config"
	"todo-api/internal/handler"
	"todo-api/internal/model"
	"todo-api/internal/repository"
	"todo-api/internal/router"
	"todo-api/internal/service"
)

func main() {
	// Загрузка переменных окружения
	cfg := config.LoadConfig()

	// Подключение к БД через GORM
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Could not connect to database: %v", err)
	}

	// Миграция всех моделей (сразу)
	err = db.AutoMigrate(&model.User{}, &model.Task{}, &model.Transaction{}, &model.Bank{})
	if err != nil {
		log.Fatalf("❌ AutoMigrate error: %v", err)
	}
	log.Println("✅ Database schema up to date (GORM AutoMigrate)")

	// Создание зависимостей
	userRepo := repository.NewUserRepo(db)
	userSvc := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc)

	taskRepo := repository.NewTaskRepo(db)
	taskSvc := service.NewTaskService(taskRepo)
	taskHandler := handler.NewTaskHandler(taskSvc)

	investmentRepo := repository.NewInvestmentRepo(db)
	investmentSvc := service.NewInvestmentService(investmentRepo)
	investmentHandler := handler.NewInvestmentHandler(investmentSvc)

	transactionRepo := repository.NewTransactionRepo(db)
	transactionSvc := service.NewTransactionService(transactionRepo)
	transactionHandler := handler.NewTransactionHandler(transactionSvc)

	// Инициализация роутера и всех маршрутов
	r := router.NewRouter(userHandler, taskHandler, transactionHandler, investmentHandler, userSvc)

	port := cfg.Port
	if port == "" {
		port = "8080"
	}
	log.Printf("🚀 Server running at http://localhost:%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("❌ Failed to run server: %v", err)
	}
}
