package main

import (
	"flag"
	"os"

	"GitHub/go-chat/backend/internal/infra/postgres"
)

func main() {
	drop := flag.Bool("drop", false, "drop database")
	flag.Parse()

	db := postgres.NewDatabaseConnection(postgres.DbConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Name:     os.Getenv("DB_NAME"),
		Password: os.Getenv("DB_PASSWORD"),
	})

	if *drop {
		err := postgres.Drop(db)
		if err != nil {
			panic(err)
		}
	}
}
