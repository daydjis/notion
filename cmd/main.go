package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // драйвер Postgres!
	_ "github.com/golang-migrate/migrate/v4/source/file"       // источник миграций из файлов
	_ "github.com/lib/pq"
	"log"
	"time"
	auth "todo-api/internal"
	"todo-api/internal/handler"
	"todo-api/internal/repository"
	"todo-api/internal/service"
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
func runMigrations(dbURL string) {
	m, err := migrate.New(
		"file://./migrations",
		dbURL,
	)
	if err != nil {
		log.Fatalf("migrate.New: %v", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("m.Up: %v", err)
	}
	//if err := m.Down(); err != nil && err != migrate.ErrNoChange {
	//	log.Fatalf("m.Up: %v", err)
	//}
}

func main() {
	//err := godotenv.Load()
	//
	//if err != nil {
	//	log.Fatal("Error loading .env")
	//}

	dbURL := "postgres://myuser1:mypass@localhost:5432/mydb1?sslmode=disable"
	runMigrations(dbURL)

	//dsn := os.Getenv("DATABASE_URL")
	db, err := waitForDB(dbURL, 10)
	if err != nil {
		log.Fatal(err)
	}

	taskRepo := repository.NewTaskRepository(db)
	taskService := service.NewTaskService(taskRepo)
	taskHandler := handler.NewTaskHandler(taskService)

	userRepo := repository.NewPostgresUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	authMiddleware, err := auth.JwtMiddleware(userService)
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	taskHandler.RegisterRoutes(r)
	userHandler.RegisterRoutes(r)
	r.POST("/login", authMiddleware.LoginHandler) // логин через middleware

	err = db.Ping()
	if err != nil {

		log.Fatal("Failed to ping DB:", err)
	}

	fmt.Println("Server is running at http://localhost:8080")
	r.Run(":8080")
}
