package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
	"todo-api/internal/handler"
	"todo-api/internal/repository"
	"todo-api/internal/service"

	_ "github.com/lib/pq"
)

func waitForDB(dsn string, retries int) (*sql.DB, error) {
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

		log.Printf("DB not ready (attempt %d/%d): %s", i+1, retries, err)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("database not reachable after %d attempts: %w", retries, err)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env")
	}
	fmt.Println(os.Getenv("DATABASE_URL"), "asdasdada")

	dsn := os.Getenv("DATABASE_URL")
	db, err := waitForDB(dsn, 10)
	if err != nil {
		log.Fatal(err)
	}

	taskRepo := repository.NewTaskRepository(db)
	taskService := service.NewTaskService(taskRepo)
	taskHandler := handler.NewTaskHandler(taskService)
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	taskHandler.RegisterRoutes(r)
	err = db.Ping()
	if err != nil {

		log.Fatal("Failed to ping DB:", err)
	}

	fmt.Println("Server is running at http://localhost:8080")
	r.Run(":8080")
}
