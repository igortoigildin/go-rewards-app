package user

import (
	"context"
	"database/sql"
	"errors"

	userEntity "github.com/igortoigildin/go-rewards-app/internal/entities/user"
)

var (
	ErrDuplicateLogin = errors.New("duplicate login")
	ErrRecordNotFound = errors.New("no records found")
)

type UserRepository struct {
	DB *sql.DB
}

// NewUserRepository returns a new instance of the repository.
func NewUserRepository(DB *sql.DB) *UserRepository {
	return &UserRepository{
		DB: DB,
	}
}

func (rep *UserRepository) Create(ctx context.Context, user *userEntity.User) error {
	err := rep.DB.QueryRowContext(ctx, "INSERT INTO users (login, password_hash)"+
		"VALUES ($1, $2) RETURNING user_id", user.Login, user.Password.Hash).Scan(&user.ID)
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

func (rep *UserRepository) Find(ctx context.Context, login string) (*userEntity.User, error) {
	var user userEntity.User

	err := rep.DB.QueryRowContext(ctx, "SELECT user_id, login, password_hash FROM users WHERE login = $1", login).Scan(
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
