package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"encoding/json"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/labstack/echo/v5"
)

func TestHealthHandler(t *testing.T) {
	dsn := "host=localhost user=taskuser password=taskpass dbname=taskdb_test port=5433 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		t.Fatal(err)
	}

	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	h := Handler{
		DB: db,
	}

	err = h.healthHandler(c)

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

	result := db.Exec(
		"TRUNCATE tasks RESTART IDENTITY",
	)

	if result.Error != nil {
		t.Fatal(err)
	}

	h := Handler{
		DB: db,
	}

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

	err = h.createTasksHandler(c)
	if err != nil {
		t.Fatal(err)
	}

	if rec.Code != http.StatusCreated {
		t.Errorf("Expected 201, got %d", rec.Code)
	}

	var task Task

	err = json.	Unmarshal(rec.Body.Bytes(), &task)
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
