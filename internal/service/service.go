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
}

// Service storage of all services.
type Service struct {
	UserService  UserService
	TokenService TokenService
}
