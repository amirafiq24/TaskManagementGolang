package main

import (
	"log/slog"
	"os"
	"task_api/config"
	"task_api/internal/database"
	"task_api/internal/observability"
	"task_api/internal/server"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	cfg, err := config.GetConfig()
	if err != nil {
		slog.Error("config error", "error", err)
		os.Exit(1)
	}

	var level slog.Level
	_ = level.UnmarshalText([]byte(cfg.Log.Level))
	slog.SetDefault(slog.New(observability.ContextHandler{
		Handler: slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: &level}),
	}))

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
