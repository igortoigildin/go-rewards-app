package token

import (
	"context"
	"database/sql"
	"errors"
	"time"

	userEntity "github.com/igortoigildin/go-rewards-app/internal/entities/user"
	"github.com/igortoigildin/go-rewards-app/internal/storage"
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

func (rep *TokenRepository) FindUserByToken(tokenHash []byte) (*userEntity.User, error) {
	query := `
	SELECT users.id, users.login, users.password_hash 
	FROM users
	INNER JOIN tokens
	ON users.id = tokens.user_id
	WHERE tokens.hash = $1 
	AND tokens.expiry > $2`

	args := []any{tokenHash, time.Now()}
	var user userEntity.User

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := rep.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.Login,
		&user.Password.Hash,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, storage.ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}
