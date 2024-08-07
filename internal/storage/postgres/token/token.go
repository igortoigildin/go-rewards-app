package token

import (
	"context"
	"database/sql"

	userEntity "github.com/igortoigildin/go-rewards-app/internal/entities/user"
)

type TokenRepository struct {
	DB *sql.DB
}

// NewTokenRepository returns a new instance of the repository.
func NewTokenRepository(DB *sql.DB) *TokenRepository {
	return &TokenRepository{
		DB: DB,
	}
}

// Adds the data for a specific token to the tokens table.
func (rep *TokenRepository) Insert(ctx context.Context, token *userEntity.Token) error {
	_, err := rep.DB.ExecContext(ctx, "INSERT INTO tokens (hash, user_id, expiry)"+
		"VALUES ($1, $2, $3)", token.Hash, token.UserID, token.Expiry)
	return err
}
