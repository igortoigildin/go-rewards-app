package server

import (
	"context"
	"net/http"

	"github.com/igortoigildin/go-rewards-app/config"
	"github.com/igortoigildin/go-rewards-app/internal/logger"
	"go.uber.org/zap"
)

func RunServer() {
	_ = initialize()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/user/register", registerUser)
}

// Init logger and load config
func initialize() *config.Config {
	cfg := config.LoadConfig()
	if err := logger.Initialize(cfg.FlagLogLevel); err != nil {
		logger.Log.Info("error while initializing logger", zap.Error(err))
	}
	return cfg
}

func registerUser(rw http.ResponseWriter, r *http.Request) {
	_, cancel := context.WithCancel(r.Context())
	defer cancel()


}
