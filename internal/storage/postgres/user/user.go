package user

import (
	"context"
	"database/sql"
	"fmt"

	userEntity "github.com/igortoigildin/go-rewards-app/internal/entities/user"
)

type UserRepository struct {
	db *sql.DB
}

// NewUserRepository returns a new instance of the repository.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (rep *UserRepository) Create(ctx context.Context, user *userEntity.User) error {
	const op = "storage.postgres.user.Create"

	err := rep.db.QueryRowContext(ctx, `INSERT INTO users (login, password_hash, balance)
		VALUES ($1, $2, $3) RETURNING user_id`, user.Login, user.Password.Hash, user.Balance).Scan(&user.UserID)
	if err != nil {
		return fmt.Errorf("%s: cannot create user: %w", op, err)
	}
	return nil
}

func (rep *UserRepository) Find(ctx context.Context, login string) (*userEntity.User, error) {
	const op = "storage.postgres.user.Find"

	var user userEntity.User
	err := rep.db.QueryRowContext(ctx, "SELECT user_id, login, password_hash FROM users WHERE login = $1", login).Scan(
		&user.UserID, &user.Login, &user.Password.Hash,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: cannot find user: %w", op, err)
	}
	return &user, nil
}

func (rep *UserRepository) Balance(ctx context.Context, UserID int64) (float64, error) {
	const op = "storage.postgres.user.Balance"

	var balance float64
	err := rep.db.QueryRowContext(ctx, "SELECT balance FROM users WHERE user_id = $1", UserID).Scan(
		&balance,
	)
	if err != nil {
		return 0, fmt.Errorf("%s: error while finding user: %w", op, err)
	}
	return balance, nil
}
