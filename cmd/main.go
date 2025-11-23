package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/Joskmo/avito-trainee-assignment-api/internal/env"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	// Logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Loading config from env
	if err := godotenv.Load(); err != nil {
		slog.Warn("the .env file wasn't read -> using default data", "warning", err)
	}
	cfg := config{
		addr: ":8080",
		db: dbConfig{
			dsn: env.GetString(
				"DATABASE_URL",
				env.GetString("GOOSE_DBSTRING", "host=localhost user=trainee password=trainee_password dbname=trainee_db sslmode=disable"),
			),
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.db.dsn)
	if err != nil {
		slog.Error("failed to connect to the database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		slog.Error("failed to ping the database", "error", err)
		os.Exit(1)
	}

	slog.Info("database connection pool ready")

	// Application
	app := application{
		config: cfg,
		db:     pool,
	}
	if err := app.run(app.mount()); err != nil {
		slog.Error("server failed to starts", "error", err)
		os.Exit(1)
	}
}
