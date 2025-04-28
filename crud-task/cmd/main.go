package main

import (
	"log"

	"github.com/berezovskiydeval/crud-task/internal/domain"
	"github.com/berezovskiydeval/crud-task/internal/server"
	"github.com/berezovskiydeval/crud-task/pkg/database"
)

func main() {
	db := database.NewPostgresDB()

	if err := db.AutoMigrate(&domain.Task{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	r := server.NewServer(db)

	r.Run(":8080")
}
