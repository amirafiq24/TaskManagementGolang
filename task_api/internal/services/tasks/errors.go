package tasks

import "errors"

var (
	ErrTaskNotFound  = errors.New("task not found")
	ErrTaskTitle     = errors.New("invalid title")
	ErrTaskStatus    = errors.New("invalid status")
	ErrInvalidID     = errors.New("invalid ID")
	ErrDatabaseError = errors.New("database error")
	ErrRedisError    = errors.New("redis error")
)
