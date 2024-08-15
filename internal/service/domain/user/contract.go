package user

import (
	"context"

	userEntity "github.com/igortoigildin/go-rewards-app/internal/entities/user"
)

type UserRepository interface {
	Create(ctx context.Context, user *userEntity.User) error
	Find(ctx context.Context, login string) (*userEntity.User, error)
	Balance(ctx context.Context, UserID int64) (int, error)
}
