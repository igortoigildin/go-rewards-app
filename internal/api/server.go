package api

import (
	"context"
	"net/http"

	"github.com/igortoigildin/go-rewards-app/config"
	"github.com/igortoigildin/go-rewards-app/internal/logger"
	"go.uber.org/zap"
)

func RunServer() {
	ctx := context.Background()
	cfg := config.LoadConfig()

	if err := logger.Initialize(cfg.FlagLogLevel); err != nil {
		logger.Log.Info("error while initializing logger", zap.Error(err))
	}

	err := http.ListenAndServe(cfg.FlagRunAddr, router(ctx, cfg))
	if err != nil {
		logger.Log.Fatal("cannot start server", zap.Error(err))
	}
}

func registerUserHandler(rw http.ResponseWriter, r *http.Request) {
	_, cancel := context.WithCancel(r.Context())
	defer cancel()

}
