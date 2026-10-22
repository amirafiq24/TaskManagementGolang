package server

import (
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	taskhandler "task_api/internal/handlers/tasks"
	taskservice "task_api/internal/services/tasks"
)

func Start(db *gorm.DB, rdb *redis.Client) {
	service := taskservice.NewTaskService(db, rdb)
	handler := taskhandler.NewTaskHandler(service)

	e := echo.New()

	e.Use(middleware.CORS("*"))

	e.GET("/health", handler.CheckHealth)
	e.POST("/tasks", handler.CreateTask)
	e.GET("/tasks", handler.GetTasks)
	e.GET("/tasks/:id", handler.GetTask)
	e.PATCH("/tasks/:id", handler.UpdateTask)
	e.DELETE("/tasks/:id", handler.DeleteTask)

	e.Start(":8000")
}
