package tasks

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"

	ent "task_api/internal/entities/tasks"
	service "task_api/internal/services/tasks"
)

func (h TaskHandler) CheckHealth(c *echo.Context) error {
	result, err := h.Services.CheckHealth()

	if err != nil {
		return c.String(http.StatusServiceUnavailable, "Unable to check health")
	}

	return c.JSON(http.StatusOK, result)
}

func (h TaskHandler) CreateTask(c *echo.Context) error {
	var req ent.CreateTask

	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "Invalid title/status")
	}

	result, err := h.Services.CreateTask(req)

	if err != nil {
		switch {
		case errors.Is(err, service.ErrTaskTitle):
			return c.String(http.StatusBadRequest, service.ErrTaskTitle.Error())
		case errors.Is(err, service.ErrTaskStatus):
			return c.String(http.StatusBadRequest, service.ErrTaskStatus.Error())
		default:
			slog.ErrorContext(c.Request().Context(), "create task failed", "error", err)
			return c.String(http.StatusInternalServerError, service.ErrDatabaseError.Error())
		}
	}

	return c.JSON(http.StatusCreated, result)
}

func (h TaskHandler) GetTasks(c *echo.Context) error {
	tasks, err := h.Services.GetTasks()

	if err != nil {
		slog.ErrorContext(c.Request().Context(), "get tasks failed", "error", err)
		return c.String(500, service.ErrDatabaseError.Error())
	}

	return c.JSON(http.StatusOK, tasks)
}

func (h TaskHandler) GetTask(c *echo.Context) error {
	id, err := parseTaskID(c)
	if err != nil {
		return err
	}

	task, err := h.Services.GetTask(c.Request().Context(), id)

	if err != nil {
		switch {
		case errors.Is(err, service.ErrTaskNotFound):
			return c.String(http.StatusNotFound, service.ErrTaskNotFound.Error())
		case errors.Is(err, service.ErrInvalidID):
			return c.String(http.StatusBadRequest, service.ErrInvalidID.Error())
		default:
			slog.ErrorContext(c.Request().Context(), "get task failed", "task_id", id, "error", err)
			return c.String(http.StatusInternalServerError, "database connection error")
		}
	}

	return c.JSON(http.StatusOK, task)
}

func (h TaskHandler) UpdateTask(c *echo.Context) error {
	id, err := parseTaskID(c)
	if err != nil {
		return err
	}

	var updatedTask ent.UpdateTask
	if err := c.Bind(&updatedTask); err != nil {
		return c.String(http.StatusBadRequest, "invalid title/status")
	}

	updates := map[string]interface{}{}

	if updatedTask.Title != nil {
		updates["title"] = *updatedTask.Title
	}

	if updatedTask.Status != nil {
		updates["status"] = *updatedTask.Status
	}

	task, err := h.Services.UpdateTask(c.Request().Context(), id, updates)

	if err != nil {
		switch {
		case errors.Is(err, service.ErrTaskNotFound):
			return c.String(http.StatusNotFound, service.ErrTaskNotFound.Error())
		case errors.Is(err, service.ErrInvalidID):
			return c.String(http.StatusBadRequest, service.ErrInvalidID.Error())
		default:
			slog.ErrorContext(c.Request().Context(), "update task failed", "task_id", id, "error", err)
			return c.String(http.StatusInternalServerError, service.ErrDatabaseError.Error())
		}
	}

	return c.JSON(http.StatusOK, task)
}

func (h TaskHandler) DeleteTask(c *echo.Context) error {
	id, err := parseTaskID(c)
	if err != nil {
		return err
	}

	err = h.Services.DeleteTask(c.Request().Context(), id)

	if err != nil {
		switch {
		case errors.Is(err, service.ErrTaskNotFound):
			return c.String(http.StatusNotFound, service.ErrTaskNotFound.Error())
		case errors.Is(err, service.ErrInvalidID):
			return c.String(http.StatusBadRequest, service.ErrInvalidID.Error())
		default:
			slog.ErrorContext(c.Request().Context(), "delete task failed", "task_id", id, "error", err)
			return c.String(http.StatusInternalServerError, service.ErrDatabaseError.Error())
		}
	}

	return c.NoContent(http.StatusOK)
}

func parseTaskID(c *echo.Context) (int, error) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return 0, echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	return id, nil
}
