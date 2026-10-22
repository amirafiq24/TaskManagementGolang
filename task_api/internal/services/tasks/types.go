package tasks

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type TaskService struct {
	DB    *gorm.DB
	Redis *redis.Client
}

func NewTaskService(db *gorm.DB, rdb *redis.Client) TaskService {
	return TaskService{DB: db, Redis: rdb}
}
