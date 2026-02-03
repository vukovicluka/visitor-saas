package main

import (
	"context"
	"flag"
	"log"

	"visitor/internal/hash"
	"visitor/internal/server"
	"visitor/internal/storage"
)

func main() {
	addr := flag.String("addr", ":8080", "HTTTP listen address")
	password := flag.String("password", "", "Dashboard password (empty = no auth)")
	databaseURL := flag.String("database-url", "postgres://visitor:visitor@localhost:5432/visitor?sslmode=disable", "PostgreSQL connection string")
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
