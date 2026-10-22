package server

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"task_api/config"
	taskhandler "task_api/internal/handlers/tasks"
	"task_api/internal/observability"
	taskservice "task_api/internal/services/tasks"
)

func Start(cfg *config.Config, db *gorm.DB, rdb *redis.Client) {
	ttl := time.Duration(cfg.Redis.TTL) * time.Second
	service := taskservice.NewTaskService(db, rdb, ttl)
	handler := taskhandler.NewTaskHandler(service)

	e := echo.New()

	e.Use(middleware.CORS("*"))

	e.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		RequestIDHandler: func(c *echo.Context, id string) {
			ctx := observability.WithRequestID(c.Request().Context(), id)
			c.SetRequest(c.Request().WithContext(ctx))
		},
	}))

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		Skipper:   func(c *echo.Context) bool { return c.Path() == "/health" },
		LogMethod: true, LogURI: true, LogStatus: true, LogLatency: true,
		LogValuesFunc: func(c *echo.Context, v middleware.RequestLoggerValues) error {
			slog.InfoContext(c.Request().Context(), "request",
				"method", v.Method, "uri", v.URI, "status", v.Status,
				"latency", v.Latency, "error", v.Error)
			return nil
		},
	}))

	e.Use(middleware.Recover())
	e.Use(middleware.ContextTimeout(time.Duration(cfg.Timeout.RequestSeconds) * time.Second))

	e.GET("/health", handler.CheckHealth)
	e.POST("/tasks", handler.CreateTask)
	e.GET("/tasks", handler.GetTasks)
	e.GET("/tasks/:id", handler.GetTask)
	e.PATCH("/tasks/:id", handler.UpdateTask)
	e.DELETE("/tasks/:id", handler.DeleteTask)

	e.Start(fmt.Sprintf(":%d", cfg.Server.Port))
}
