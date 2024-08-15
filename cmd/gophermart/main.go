package main

import (
	"context"
	"database/sql"
	"net/http"

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

	conn, err := sql.Open("pgx", cfg.FlagDBURI)
	if err != nil {
		logger.Log.Info("error while connecting to DB", zap.Error(err))
	}
	defer conn.Close()
	logger.Log.Info("database connection pool established")

	repository := storage.NewRepository(conn)
	services := service.NewService(repository)

	go api.RunAccrualUpdates(ctx, cfg, services)

	err = http.ListenAndServe(cfg.FlagRunAddr, api.Router(services))
	if err != nil {
		logger.Log.Fatal("cannot start server", zap.Error(err))
	}
}
