package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"GitHub/go-chat/backend/internal/infra/postgres/db"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DbConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

func NewDatabaseConnection(ctx context.Context, conf DbConfig) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		conf.Host, conf.Port, conf.User, conf.Name, conf.Password,
	)

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("parse config error: %w", err)
	}

	var pool *pgxpool.Pool
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		pool, err = pgxpool.NewWithConfig(ctx, config)
		if err != nil {
			if i < maxRetries-1 {
				log.Printf("Failed to create pool, retrying in 2s... (%d/%d)", i+1, maxRetries)
				time.Sleep(2 * time.Second)
				continue
			}
			return nil, fmt.Errorf("create pool error: %w", err)
		}

		if err := pool.Ping(ctx); err != nil {
			if i < maxRetries-1 {
				log.Printf("Failed to ping, retrying in 2s... (%d/%d)", i+1, maxRetries)
				pool.Close()
				time.Sleep(2 * time.Second)
				continue
			}
			return nil, fmt.Errorf("ping error: %w", err)
		}
		break
	}

	migrationsDir := "internal/infra/postgres/migrations/up"
	if err := RunMigrations(ctx, pool, migrationsDir); err != nil {
		return nil, fmt.Errorf("migrate error: %w", err)
	}

	log.Printf("💿 Connected to database %s", conf.Name)

	return pool, nil
}

func NewQueries(pool *pgxpool.Pool) *db.Queries {
	return db.New(pool)
}
