package dataaccess

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"task_api/internal/models"
)

func CacheGetTask(ctx context.Context, rdb *redis.Client, id int) (*models.Task, error) {
	key := fmt.Sprintf(`task:%d`, id)
	data, err := rdb.Get(ctx, key).Bytes()

	// 1. cache miss -> whatt is err?
	// 2. actual err
	// 3. cache hitt -> err nil, data is not empty

	switch {
	case err == nil: // cache hit
		var task models.Task

		if err := json.Unmarshal(data, &task); err != nil {
			return nil, err
		}

		return &task, nil

	case errors.Is(err, redis.Nil): //cache miss
		return nil, nil

	default: // redis error
		return nil, err
	}
}

func CacheSetTask(ctx context.Context, rdb *redis.Client, task *models.Task, ttl time.Duration) error {
	// add item to cache
	key := fmt.Sprintf(`task:%d`, task.ID)
	value, err := json.Marshal(task)

	if err != nil {
		return err
	}

	return rdb.Set(ctx, key, value, ttl).Err()
}

func CacheDeleteTask(ctx context.Context, rdb *redis.Client, id int) error {
	key := fmt.Sprintf(`task:%d`, id)
	return rdb.Del(ctx, key).Err()
}
