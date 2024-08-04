package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/igortoigildin/go-rewards-app/config"
	"github.com/igortoigildin/go-rewards-app/internal/logger"
	"github.com/igortoigildin/go-rewards-app/internal/models"
	"go.uber.org/zap"
)

var (
	ErrDuplicateLogin = errors.New("duplicate login")
	ErrRecordNotFound = errors.New("no records found")
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
	_, err = rep.conn.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS users (id bigserial PRIMARY KEY, login TEXT UNIQUE NOT NULL,"+
		"password_hash bytea NOT NULL);")
	if err != nil {
		logger.Log.Info("error while creating users table")
	}
	return rep
}

func (rep *Repository) CreateUser(user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := rep.conn.QueryRowContext(ctx, "INSERT INTO users (login, password_hash)"+
		"VALUES ($1, $2) RETURNING id", user.Login, user.Password.Hash).Scan(&user.ID)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_login_key"`:
			return ErrDuplicateLogin
		default:
			return err
		}
	}
	return nil
}

func (rep *Repository) GetUserByLogin(login string) (*models.User, error) {
	var user models.User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := rep.conn.QueryRowContext(ctx, "SELECT id, login, password_hash FROM users WHERE login = $1", login).Scan(
		&user.ID, &user.Login, &user.Password.Hash,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}
