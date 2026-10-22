package server

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"task_api/config"
	taskhandler "task_api/internal/handlers/tasks"
	taskservice "task_api/internal/services/tasks"
)

func Start(cfg *config.Config, db *gorm.DB, rdb *redis.Client) {
	ttl := time.Duration(cfg.Redis.TTL) * time.Second
	service := taskservice.NewTaskService(db, rdb, ttl)
	handler := taskhandler.NewTaskHandler(service)

	e := echo.New()

	e.Use(middleware.CORS("*"))

	e.GET("/health", handler.CheckHealth)
	e.POST("/tasks", handler.CreateTask)
	e.GET("/tasks", handler.GetTasks)
	e.GET("/tasks/:id", handler.GetTask)
	e.PATCH("/tasks/:id", handler.UpdateTask)
	e.DELETE("/tasks/:id", handler.DeleteTask)

	e.Start(fmt.Sprintf(":%d", cfg.Server.Port))
}
