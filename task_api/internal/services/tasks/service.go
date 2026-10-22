package tasks

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	cache "task_api/internal/data_access"
	repo "task_api/internal/data_access/tasks"
	ent "task_api/internal/entities/tasks"
	"task_api/internal/models"

	"gorm.io/gorm"
)

func (h *TaskService) CheckHealth() (models.HealthStatus, error) {
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

	return systemHealth, dbCheck
}

func (h *TaskService) CreateTask(req ent.CreateTask) (*models.Task, error) {
	// input validation
	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		return nil, ErrTaskTitle
	}

	if !slices.Contains([]string{"todo", "inprogress", "blocked", "onhold", "done", "cancelled"}, req.Status) {
		return nil, ErrTaskStatus
	}

	task, err := repo.CreateTask(h.DB, req.Title, req.Status)

	if err != nil {
		return nil, err
	}

	return task, nil
}

func (h *TaskService) GetTasks() ([]*models.Task, error) {
	tasks, err := repo.GetTasks(h.DB)

	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (h *TaskService) GetTask(c context.Context, id int) (*models.Task, error) {
	if id <= 0 {
		return nil, ErrInvalidID
	}

	// try cache first
	task, err := cache.CacheGetTask(c, h.Redis, id)

	// redis error
	if err != nil {
		fmt.Println("redis error")
	}

	// cache hit
	if task != nil {
		return task, nil
	}

	// here, err = nil, task = nil. means cache miss. retrieve from db

	task, err = repo.GetTask(h.DB, id)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrTaskNotFound
	}

	if err != nil {
		return nil, err
	}

	err = cache.CacheSetTask(c, h.Redis, task, h.TTL)

	if err != nil {
		fmt.Println("cache failed")
	} else {
		fmt.Println("cache add success")
	}

	return task, nil
}

func (h *TaskService) UpdateTask(ctx context.Context, id int, req map[string]interface{}) (*models.Task, error) {
	if id <= 0 {
		return nil, ErrInvalidID
	}

	task, err := repo.UpdateTask(h.DB, id, req)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrTaskNotFound
	}

	if err != nil {
		return nil, err
	}

	err = cache.CacheDeleteTask(ctx, h.Redis, id)

	if err != nil {
		slog.ErrorContext(ctx, "cache delete failed", "err", err)
	}

	return task, nil
}

func (h *TaskService) DeleteTask(ctx context.Context, id int) error {
	err := repo.DeleteTask(h.DB, id)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrTaskNotFound
	}

	if err != nil {
		return err
	}

	err = cache.CacheDeleteTask(ctx, h.Redis, id)

	if err != nil {
		fmt.Print("failed to delete cache")
	}

	return nil
}
