package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // драйвер Postgres!
	_ "github.com/golang-migrate/migrate/v4/source/file"       // источник миграций из файлов
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
	"time"
	auth "todo-api/internal"
	"todo-api/internal/handler"
	"todo-api/internal/repository"
	"todo-api/internal/service"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ No .env file found, using environment variables only")
	}
}

// runMigrations применяет миграции базы данных
func runMigrations(dbURL string) error {
	m, err := migrate.New(
		"file://migrations",
		dbURL,
	)
	if err != nil {
		return fmt.Errorf("failed to initialize migrations: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	log.Println("✅ Database migrations applied successfully")
	return nil
}

// waitForDB ждёт готовности БД с повторными попытками
func waitForDB(dsn string, retries int, delay time.Duration) (*sql.DB, error) {
	var db *sql.DB
	var err error
	for i := 0; i < retries; i++ {
		db, err = sql.Open("postgres", dsn)
		if err == nil {
			err = db.Ping()
		}
		if err == nil {
			return db, nil
		}
		log.Printf("⏳ Database not ready (attempt %d/%d): %v", i+1, retries, err)
		time.Sleep(delay)
	}
	return nil, fmt.Errorf("database not reachable after %d attempts: %w", retries, err)
}

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("❌ DATABASE_URL environment variable is required")
	}

	// Запуск миграций
	if err := runMigrations(dbURL); err != nil {
		log.Fatalf("❌ Migration error: %v", err)
	}

	// Ожидание готовности БД
	db, err := waitForDB(dbURL, 10, 2*time.Second)
	if err != nil {
		log.Fatalf("❌ Could not connect to database: %v", err)
	}
	defer db.Close()

	// Репозитории и сервисы
	userRepo := repository.NewPostgresUserRepository(db)
	userSvc := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc)

	taskRepo := repository.NewTaskRepository(db)
	taskSvc := service.NewTaskService(taskRepo)
	taskHandler := handler.NewTaskHandler(taskSvc)

	// JWT Middleware
	authMiddleware, err := auth.JwtMiddleware(userSvc)
	if err != nil {
		log.Fatalf("❌ Failed to initialize JWT middleware: %v", err)
	}

	// Инициализация роутера Gin
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Открытые маршруты
	router.POST("/register", userHandler.RegisterHandler)
	router.POST("/login", authMiddleware.LoginHandler)

	// Защищённые маршруты
	api := router.Group("/")
	api.Use(authMiddleware.MiddlewareFunc())
	{
		userHandler.RegisterRoutes(api)
		taskHandler.RegisterRoutes(api)
	}

	// Запуск сервера
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("🚀 Server running at http://localhost:%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("❌ Failed to run server: %v", err)
	}
}
