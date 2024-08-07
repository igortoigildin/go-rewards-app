package api

import (
	"database/sql"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/igortoigildin/go-rewards-app/config"
	"github.com/igortoigildin/go-rewards-app/internal/logger"
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

	runMigration(conn)

	repositories := pg.NewRepository(conn)
	services := domain.NewService(repositories)

	err = http.ListenAndServe(cfg.FlagRunAddr, router(services, cfg))
	if err != nil {
		logger.Log.Fatal("cannot start server", zap.Error(err))
	}
}

func runMigration(conn *sql.DB) {
	migrationDriver, err := postgres.WithInstance(conn, &postgres.Config{})
	if err != nil {
		logger.Log.Fatal("error performing migration", zap.Error(err))
	}

	migrator, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", migrationDriver)
	if err != nil {
		logger.Log.Fatal("error performing migration", zap.Error(err))
	}

	err = migrator.Up()
	if err != nil && err != migrate.ErrNoChange {
		logger.Log.Fatal("error performing migration", zap.Error(err))
	}

	logger.Log.Info("database migrations applied")
}
