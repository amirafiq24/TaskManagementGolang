package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	taskhandler "task_api/internal/handlers/tasks"
	"task_api/internal/models"
	taskservice "task_api/internal/services/tasks"
)

func TestHealthHandler(t *testing.T) {
	dsn := "host=localhost user=taskuser password=taskpass dbname=taskdb_test port=5433 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		t.Fatal(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		Protocol: 2,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		t.Fatal(err)
	}

	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	service := taskservice.NewTaskService(db, rdb)
	handler := taskhandler.NewTaskHandler(service)

	err = handler.CheckHealth(c)

	if err != nil {
		t.Fatal(err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestCreateTask(t *testing.T) {
	dsn := "host=localhost user=taskuser password=taskpass dbname=taskdb_test port=5433 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		t.Fatal(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		Protocol: 2,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		t.Fatal(err)
	}

	result := db.Exec(
		"TRUNCATE tasks RESTART IDENTITY",
	)

	if result.Error != nil {
		t.Fatal(err)
	}

	service := taskservice.NewTaskService(db, rdb)
	handler := taskhandler.NewTaskHandler(service)

	e := echo.New()

	body := `{
		"title": "learn testing",
		"status": "todo"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/tasks",
		strings.NewReader(body),
	)

	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	err = handler.CreateTask(c)
	if err != nil {
		t.Fatal(err)
	}

	if rec.Code != http.StatusCreated {
		t.Errorf("Expected 201, got %d", rec.Code)
	}

	var task models.Task

	err = json.Unmarshal(rec.Body.Bytes(), &task)
	if err != nil {
		t.Fatal(err)
	}

	if task.Title != "learn testing" {
		t.Errorf("expected title learn testing, got %s", task.Title)
	}

	if task.Status != "todo" {
		t.Errorf("expected status todo, got %s", task.Status)
	}

	if task.ID != 1 {
		t.Errorf("expected ID 1, got %d", task.ID)
	}
}
