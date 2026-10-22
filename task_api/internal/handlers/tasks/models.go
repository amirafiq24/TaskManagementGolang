package tasks

import (
	services "task_api/internal/services/tasks"

	"github.com/labstack/echo/v5"
)

type ITasksHandler interface {
	CheckHealth(c *echo.Context) error
	CreateTask(c *echo.Context) error
	GetTask(c *echo.Context) error
	GetTasks(c *echo.Context) error
	UpdateTask(c *echo.Context) error
	DeleteTask(c *echo.Context) error
}

type TaskHandler struct {
	Services services.TaskService
}

func NewTaskHandler(service services.TaskService) ITasksHandler {
	return &TaskHandler{
		Services: service,
	}
}
