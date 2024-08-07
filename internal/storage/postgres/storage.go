package pg

import (
	"database/sql"

	"github.com/igortoigildin/go-rewards-app/internal/storage"
	"github.com/igortoigildin/go-rewards-app/internal/storage/postgres/token"
	"github.com/igortoigildin/go-rewards-app/internal/storage/postgres/user"
)

func NewRepository(DB *sql.DB) *storage.Repository {
	return &storage.Repository{
		User:  user.NewUserRepository(DB),
		Token: token.NewTokenRepository(DB),
	}
}
