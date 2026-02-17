package main

import (
	"context"
	"flag"
	"log"
	"os"

	"GitHub/go-chat/backend/internal/config"
	"GitHub/go-chat/backend/internal/infra/postgres"
)

func main() {
	if err := config.LoadDotEnv(); err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	drop := flag.Bool("drop", false, "drop database")
	flag.Parse()

	ctx := context.Background()
	pool, err := postgres.NewDatabaseConnection(ctx, postgres.DbConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Name:     os.Getenv("DB_NAME"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	if *drop {
		log.Println("⚠️  Database dropping is now handled by migrations")
		log.Println("⚠️  Use migration down commands to revert schema changes")
	}
}
