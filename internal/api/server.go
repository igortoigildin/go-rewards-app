package api

import (
	"database/sql"
	"net/http"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/igortoigildin/go-rewards-app/config"
	"github.com/igortoigildin/go-rewards-app/internal/logger"
	"github.com/igortoigildin/go-rewards-app/internal/service"
	"github.com/igortoigildin/go-rewards-app/internal/service/domain"
	pg "github.com/igortoigildin/go-rewards-app/internal/storage/postgres"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func RunServer() {
	cfg := config.LoadConfig()

	if err := logger.Initialize(cfg.FlagLogLevel); err != nil {
		logger.Log.Info("error while initializing logger", zap.Error(err))
	}

	conn, err := sql.Open("pgx", cfg.FlagDBURI)
	if err != nil {
		logger.Log.Info("error while connecting to DB", zap.Error(err))
	}
	defer conn.Close()
	logger.Log.Info("database connection pool established")

	repositories := pg.NewRepository(conn)
	services := domain.NewService(repositories)

	app := newApp(services, cfg)

	go app.services.OrderService.UpdateAccruals(cfg)

	err = http.ListenAndServe(cfg.FlagRunAddr, router(app))
	if err != nil {
		logger.Log.Fatal("cannot start server", zap.Error(err))
	}
}

type app struct {
	services service.Service
	cfg      *config.Config
}

func newApp(service *service.Service, cfg *config.Config) *app {
	return &app{
		services: *service,
		cfg:      cfg,
	}
}
