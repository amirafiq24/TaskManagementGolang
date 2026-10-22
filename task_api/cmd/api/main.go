package main

import (
	"log/slog"
	"os"
	"task_api/config"
	"task_api/internal/database"
	"task_api/internal/server"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	cfg, err := config.GetConfig()
	if err != nil {
		slog.Error("config setup")
	}

	db, err := database.ConnectDB(&cfg.DB)
	if err != nil {
		panic(err)
	}

	rdb, err := database.ConnectRedis(&cfg.Redis)

	if err != nil {
		panic(err)
	}
	defer rdb.Close()

	server.Start(cfg, db, rdb)
}
