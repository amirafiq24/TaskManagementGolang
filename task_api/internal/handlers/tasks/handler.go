package tasks

import (
	"strconv"
	"net/http"
	"github.com/labstack/echo/v5"

	"task_api/internal/models"
	repo "task_api/internal/data_access/tasks"
)

func (h Handler) HealthHandler(c *echo.Context) error {
	var dbHealth models.DependencyStatus
	dbCheck := repo.HealthCheck(h.DB)
	if dbCheck != nil {
		dbHealth.Status = "Down"
		dbHealth.Error = dbCheck.Error()
	} else {
		dbHealth.Status = "Up"
	}

	systemHealth := models.HealthStatus{
		Dependencies: make(map[string]models.DependencyStatus),
	}
	systemHealth.Dependencies["Database"] = dbHealth

	if dbCheck != nil {
		systemHealth.Status = "Unhealthy"
		systemHealth.Uptime = "Placeholder"
	} else {
		systemHealth.Status = "Healthy"
		systemHealth.Uptime = "Placeholder"
	}

	return c.JSON(http.StatusOK, systemHealth)
}

func (h Handler) CreateTasksHandler(c *echo.Context) error {
	type CreateTask struct {
		Title string `json:"title"`
		Status string `json:"status"`
	}

	var req CreateTask

	if err := c.Bind(&req); err != nil	{
		return c.String(http.StatusBadRequest, "Bad request")
	}

	task, err := repo.CreateTask(h.DB, req.Title, req.Status)

	if err != nil {
		return c.String(http.StatusBadRequest, "Failed to create task")
	}

	return c.JSON(http.StatusCreated, task)
}

func (h Handler) GetAllTaskHandler(c *echo.Context) error {
	tasks, err := repo.GetTasks(h.DB)
	if err != nil {
		return c.String(500, "Something wrong idk")
	}
	return c.JSON(http.StatusOK, tasks)
}

func (h Handler) GetTaskHandler(c *echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid ID")
	}

	task, err := repo.GetTask(h.DB, id)

	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid id kot")
	}

	return c.JSON(http.StatusOK, task)
}

func (h Handler) UpdateTaskHandler(c *echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid ID")
	}

	type UpdateTask struct {
		Title *string `json:"title"`
		Status *string `json:"status"`
	}

	var updatedTask UpdateTask
	if err := c.Bind(&updatedTask); err != nil	{
		return c.String(http.StatusBadRequest, "Bad request: Invalid ID")
	}

	updates := map[string]interface{}{}

	if updatedTask.Title != nil {
		updates["title"] = *updatedTask.Title
	}

	if updatedTask.Status != nil {
		updates["status"] = *updatedTask.Status
	}

	task, err := repo.UpdateTask(h.DB, id, updates)
	
	if err != nil {
		return c.String(http.StatusBadRequest, "something wrong la")
	}

	return c.JSON(http.StatusOK, task)
}

func (h Handler) DeleteTaskHandler (c *echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid ID")
	}

  	err = repo.DeleteTask(h.DB, id)

	if err != nil {
		return c.String(http.StatusBadRequest, "something wrong idk")
	}

	return c.NoContent(http.StatusOK)
}