package server

import (
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"gorm.io/gorm"

	taskhandler "task_api/internal/handlers/tasks"
)

func Start(db *gorm.DB) {
	handler := taskhandler.Handler{
		DB: db,	
	}

	e := echo.New()

	e.Use(middleware.CORS("*"))

	e.GET("/health", handler.HealthHandler)
	e.POST("/tasks", handler.CreateTasksHandler)
	e.GET("/tasks", handler.GetAllTaskHandler)
	e.GET("/tasks/:id", handler.GetTaskHandler)
	e.PATCH("/tasks/:id", handler.UpdateTaskHandler)
	e.DELETE("/tasks/:id", handler.DeleteTaskHandler)

	e.Start(":8000")
}