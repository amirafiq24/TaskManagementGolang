package main

import (
	"log/slog"
	"os"
	"task_api/internal/database"
	"task_api/internal/server"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	db, err := database.ConnectDB()
	if err != nil {
		panic(err)
	}

	rdb, err := database.ConnectRedis()

	if err != nil {
		panic(err)
	}
	defer rdb.Close()

	server.Start(db, rdb)
}
