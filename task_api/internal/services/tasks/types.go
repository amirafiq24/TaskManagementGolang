package tasks

import (
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type TaskService struct {
	DB    *gorm.DB
	Redis *redis.Client
	TTL   time.Duration
}

func NewTaskService(db *gorm.DB, rdb *redis.Client, ttl time.Duration) TaskService {
	return TaskService{DB: db, Redis: rdb, TTL: ttl}
}
