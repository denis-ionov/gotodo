package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	"todo-api/internal/config"
	"todo-api/internal/handlers"
	"todo-api/internal/middleware"
	"todo-api/internal/repository"
	"todo-api/internal/repository/memory"
	"todo-api/internal/repository/postgres"
	"todo-api/internal/service"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func runMigrations(db *sql.DB) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(db, "./migrations"); err != nil {
		return err
	}

	return nil
}

func main() {
	cfg := config.LoadConfig()

	var db *sql.DB
	var useInMemory = false

	db, err := sql.Open("postgres", cfg.GetDBConnectionString())
	if err != nil {
		log.Printf("Failed to open database connection: %v", err)
		useInMemory = true
	} else {
		db.SetMaxOpenConns(25)
		db.SetMaxIdleConns(25)
		db.SetConnMaxLifetime(5 * time.Minute)

		connected := make(chan bool, 1)
		go func() {
			err := db.Ping()
			connected <- (err == nil)
		}()

		select {
		case success := <-connected:
			if !success {
				log.Println("Failed to connect to PostgreSQL, falling back to in-memory storage")
				db.Close()
				useInMemory = true
			} else {
				log.Println("Successfully connected to PostgreSQL")

				if _, err := os.Stat("./migrations"); err == nil {
					if err := runMigrations(db); err != nil {
						log.Printf("Warning: Failed to run migrations: %v", err)
					} else {
						log.Println("Database migrations completed successfully")
					}
				} else {
					log.Println("Migrations directory not found, skipping migrations")
				}
			}
		case <-time.After(5 * time.Second):
			log.Println("Database connection timeout, falling back to in-memory storage")
			db.Close()
			useInMemory = true
		}
	}

	var repo *repository.Repository
	if db != nil && !useInMemory {
		repo = &repository.Repository{
			Task: postgres.NewTaskRepository(db),
			User: postgres.NewUserRepository(db),
		}
	} else {
		repo = &repository.Repository{
			Task: memory.NewTaskRepository(),
			User: memory.NewUserRepository(),
		}
	}

	authService := service.NewAuthService(repo.User, cfg.JWTSecret)
	taskService := service.NewTaskService(repo.Task)
	userService := service.NewUserService(repo.User)

	authHandler := handlers.NewAuthHandler(authService)
	taskHandler := handlers.NewTaskHandler(taskService)
	userHandler := handlers.NewUserHandler(userService)

	r := gin.Default()

	publicRoute := r.Group("/api")
	{
		publicRoute.POST("/login", authHandler.Login)
		publicRoute.POST("/register", authHandler.Register)
	}

	protectedRoute := r.Group("/api")
	protectedRoute.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		protectedRoute.GET("/tasks", taskHandler.GetTasks)
		protectedRoute.GET("/tasks/:id", taskHandler.GetTask)
		protectedRoute.POST("/tasks", taskHandler.CreateTask)
		protectedRoute.PUT("/tasks/:id", taskHandler.UpdateTask)
		protectedRoute.DELETE("/tasks/:id", taskHandler.DeleteTask)

		protectedRoute.GET("/users", userHandler.GetUsers)
		protectedRoute.GET("/users/:id", userHandler.GetUser)
		protectedRoute.POST("/users", userHandler.CreateUser)
		protectedRoute.PUT("/users/:id", userHandler.UpdateUser)
		protectedRoute.DELETE("/users/:id", userHandler.DeleteUser)

		protectedRoute.GET("/profile", authHandler.GetProfile)
	}

	port := ":" + cfg.Port
	log.Printf("Server is running on port %s", cfg.Port)
	log.Printf("Database mode: %s", map[bool]string{true: "in-memory", false: "PostgreSQL"}[useInMemory])

	if err := r.Run(port); err != nil {
		log.Fatal(err)
	}
}
