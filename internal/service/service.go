package service

import (
	"context"
	"time"

	userEntity "github.com/igortoigildin/go-rewards-app/internal/entities/user"
)

type UserService interface {
	Find(ctx context.Context, login string) (*userEntity.User, error)
	Create(ctx context.Context, user *userEntity.User) error
}

type TokenService interface {
	NewToken(ctx context.Context, userID int64, ttl time.Duration) (*userEntity.Token, error)
	FindUserByToken(tokenHash []byte) (*userEntity.User, error)
}

type OrderService interface {
	InsertOrder(ctx context.Context, number string, userID int64) (int64, error)
	ValidateOrder(number string) (bool, error)
}

// Service storage of all services.
type Service struct {
	UserService  UserService
	TokenService TokenService
	OrderService OrderService
}
