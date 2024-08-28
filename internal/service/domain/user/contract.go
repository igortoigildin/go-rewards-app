package user

import (
	"context"

	userEntity "github.com/igortoigildin/go-rewards-app/internal/entities/user"
)

//go:generate mockgen -package mocks -destination=../../../../mocks/userRepository.go github.com/igortoigildin/go-rewards-app/internal/service/domain/user UserRepository
type UserRepository interface {
	Create(ctx context.Context, user *userEntity.User) error
	Find(ctx context.Context, login string) (*userEntity.User, error)
	Balance(ctx context.Context, UserID int64) (float64, error)
}
