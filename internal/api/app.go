package api

import (
	"github.com/igortoigildin/go-rewards-app/config"
	"github.com/igortoigildin/go-rewards-app/internal/storage"
)

type app struct {
	storage storage.Storage
	cfg     *config.Config
}

func newApp(s storage.Storage, cfg *config.Config) *app {
	return &app{
		storage: s,
		cfg:     cfg,
	}
}
