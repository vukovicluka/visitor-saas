package main

import (
	"context"
	"flag"
	"log"

	"visitor/internal/storage"
)

func main() {
	databaseURL := flag.String("database-url", "postgres://visitor:visitor@localhost:5432/visitor?sslmode=disable", "PostgreSQL connection string")
	flag.Parse()

	ctx := context.Background()

	db, err := storage.New(ctx, *databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	
	defer db.Close()

	log.Println("Database connected and migrated successfully")
}
