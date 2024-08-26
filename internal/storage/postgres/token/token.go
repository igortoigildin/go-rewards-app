package token

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	userEntity "github.com/igortoigildin/go-rewards-app/internal/entities/user"
)

type TokenRepository struct {
	db *sql.DB
}

// NewTokenRepository returns a new instance of the repository.
func NewTokenRepository(db *sql.DB) *TokenRepository {
	return &TokenRepository{
		db: db,
	}
}

// Adds the data for a specific token to the tokens table.
func (rep *TokenRepository) Insert(ctx context.Context, token *userEntity.Token) error {
	const op = "storage.postgres.token.Insert"

	_, err := rep.db.ExecContext(ctx, "INSERT INTO tokens (hash, user_id, expiry)"+
		"VALUES ($1, $2, $3)", token.Hash, token.UserID, token.Expiry)
	if err != nil {
		return fmt.Errorf("%s: user not found: %w", op, err)
	}
	return nil
}

func (rep *TokenRepository) FindUserByToken(ctx context.Context, tokenHash []byte) (*userEntity.User, error) {
	const op = "storage.postgres.token.FindUserByToken"

	query := `
	SELECT users.user_id, users.login, users.password_hash 
	FROM users
	INNER JOIN tokens
	ON users.user_id = tokens.user_id
	WHERE tokens.hash = $1 
	AND tokens.expiry > $2`
	args := []any{tokenHash, time.Now()}
	var user userEntity.User
	err := rep.db.QueryRowContext(ctx, query, args...).Scan(
		&user.UserID,
		&user.Login,
		&user.Password.Hash,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: user not found: %w", op, err)
	}
	return &user, nil
}
