package main

import (
	"task_api/internal/database"
	"task_api/internal/server"
)

func main() {
	db, err := database.ConnectDB()

	if err != nil {
		panic("noo")
	}

	server.Start(db)
}