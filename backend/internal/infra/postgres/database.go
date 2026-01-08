package postgres

import (
	"context"
	"fmt"
	"log"

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
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=require",
		conf.Host, conf.Port, conf.User, conf.Name, conf.Password,
	)

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("parse config error: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("create pool error: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping error: %w", err)
	}

	migrationsDir := "./migrations"
	if err := RunMigrations(ctx, pool, migrationsDir); err != nil {
		return nil, fmt.Errorf("migrate error: %w", err)
	}

	log.Printf("💿 Connected to database %s", conf.Name)

	return pool, nil
}

func NewQueries(pool *pgxpool.Pool) *db.Queries {
	return db.New(pool)
}
