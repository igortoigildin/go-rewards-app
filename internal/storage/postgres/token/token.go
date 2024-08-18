package token

import (
	"context"
	"database/sql"
	"errors"
	"time"

	userEntity "github.com/igortoigildin/go-rewards-app/internal/entities/user"
)

var (
	ErrDuplicateLogin = errors.New("duplicate login")
	ErrRecordNotFound = errors.New("no records found")
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
	_, err := rep.db.ExecContext(ctx, "INSERT INTO tokens (hash, user_id, expiry)"+
		"VALUES ($1, $2, $3)", token.Hash, token.UserID, token.Expiry)
	return err
}

func (rep *TokenRepository) FindUserByToken(ctx context.Context, tokenHash []byte) (*userEntity.User, error) {
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
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}
