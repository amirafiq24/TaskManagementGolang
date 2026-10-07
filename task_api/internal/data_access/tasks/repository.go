package tasks

import (
	"context"
	"time"
	"gorm.io/gorm"

	"task_api/internal/models"
)

func HealthCheck(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		2*time.Second,
	)

	defer cancel()

	return sqlDB.PingContext(ctx)
}

func GetTasks(db *gorm.DB) ([]models.Task, error) {
	var tasks []models.Task
	result := db.Order("id ASC").Find(&tasks)

	if result.Error != nil {
		return nil, result.Error
	}

	return tasks, nil
}

func GetTask(db *gorm.DB, id int) (models.Task, error) {
	var task models.Task
	result := db.First(&task, id)

	if result.Error != nil {
		return models.Task{}, result.Error
	}

	return task, nil
}

func CreateTask(db *gorm.DB, title string, status string) (models.Task, error) {
	task := models.Task {
		Title: title,
		Status: status,
	}

	result := db.Create(&task)

	if result.Error != nil {
		return models.Task{}, result.Error
	}

	return task, nil
}

func UpdateTask(db *gorm.DB, id int, updates map[string]interface{}) (models.Task, error) {
	var task models.Task

	result := db.Model(&task).Where("id = ?", id).Updates(updates)

	if result.Error != nil {
		return models.Task{}, result.Error
	}

	if result.RowsAffected == 0 {
		return models.Task{}, gorm.ErrRecordNotFound
	}

	result = db.First(&task, id)
	if result.Error != nil {
		return models.Task{}, result.Error
	}

	return task, nil
}

func DeleteTask(db *gorm.DB, id int) error {
	result := db.Delete(&models.Task{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}