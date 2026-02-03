package main

import (
	"context"
	"flag"
	"log"
	"os"

	"visitor/internal/hash"
	"visitor/internal/server"
	"visitor/internal/storage"
)

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	addr := flag.String("addr", envOrDefault("ADDR", ":8080"), "HTTTP listen address")
	password := flag.String("password", envOrDefault("PASSWORD", ""), "Dashboard password (empty = no auth)")
	databaseURL := flag.String("database-url", envOrDefault("DATABASE_URL", "postgres://visitor:visitor@localhost:5432/visitor?sslmode=disable"), "PostgreSQL connection string")
	flag.Parse()

	ctx := context.Background()

	db, err := storage.New(ctx, *databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	
	defer db.Close()

	hasher := hash.NewManager(db.Pool())

	srv := server.New(*addr, db, hasher, *password)

	log.Printf("Listening on %s", *addr)
	if err := srv.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
