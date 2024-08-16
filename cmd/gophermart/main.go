package main

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/igortoigildin/go-rewards-app/config"
	"github.com/igortoigildin/go-rewards-app/internal/api"
	"github.com/igortoigildin/go-rewards-app/internal/logger"
	"github.com/igortoigildin/go-rewards-app/internal/service"
	"github.com/igortoigildin/go-rewards-app/internal/storage"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

func main() {
	cfg := config.LoadConfig()
	ctx := context.Background()

	if err := logger.Initialize(cfg.FlagLogLevel); err != nil {
		logger.Log.Info("error while initializing logger", zap.Error(err))
	}

	db, err := sql.Open("pgx", cfg.FlagDBURI)
	if err != nil {
		logger.Log.Fatal("error while connecting to DB", zap.Error(err))
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Log.Fatal("error while closing db connection", zap.Error(err))
		}
	}()

	logger.Log.Info("database connection pool established")

	instance, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		logger.Log.Fatal("migration error", zap.Error(err))
	}

	migrator, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", instance)
	if err != nil {
		logger.Log.Fatal("migration error", zap.Error(err))
	}
	// migrate.ErrNoChange
	err = migrator.Up()
	if err != nil || errors.Is(err, migrate.ErrNoChange) {
		logger.Log.Fatal("migration error", zap.Error(err))
	}

	logger.Log.Info("database connection pool established")

	repository := storage.NewRepository(db)
	services := service.NewService(repository)

	go api.RunAccrualUpdates(ctx, cfg, services)

	err = http.ListenAndServe(cfg.FlagRunAddr, api.Router(services))
	if err != nil {
		logger.Log.Fatal("database migrations applied", zap.Error(err))
	}
}
