package storage

import (
	"context"
	"database/sql"

	"github.com/igortoigildin/go-rewards-app/config"
	"github.com/igortoigildin/go-rewards-app/internal/logger"
	"go.uber.org/zap"
)

type Repository struct {
	conn *sql.DB
}

func NewRepository(conn *sql.DB) *Repository {
	return &Repository{
		conn: conn,
	}
}

func InitPostgresRepo(ctx context.Context, cfg *config.Config) *Repository {
	DBURI := cfg.FlagDBURI
	conn, err := sql.Open("pgx", DBURI)
	if err != nil {
		logger.Log.Info("error while connecting to DB", zap.Error(err))
	}
	rep := NewRepository(conn)
	// TODO Create tables as needed
	return rep
}