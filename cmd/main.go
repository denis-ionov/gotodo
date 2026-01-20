package main

import (
	"log"

	"todo-api/internal/handlers"
	"todo-api/internal/middleware"
	"todo-api/internal/repository"
	"todo-api/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	taskRepository := repository.NewTaskRepository()
	userRepository := repository.NewUserRepository()

	jwtSecret := "secret_api_key"

	authService := service.NewAuthService(userRepository, jwtSecret)
	taskService := service.NewTaskService(taskRepository)
	userService := service.NewUserService(userRepository)

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
	protectedRoute.Use(middleware.AuthMiddleware(jwtSecret))
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

	port := ":8080"

	if err := r.Run(port); err != nil {
		log.Fatal(err)
	}
}
