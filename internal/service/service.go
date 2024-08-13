package service

import (
	"context"
	"errors"
	"time"

	"github.com/igortoigildin/go-rewards-app/config"
	orderEntity "github.com/igortoigildin/go-rewards-app/internal/entities/order"
	userEntity "github.com/igortoigildin/go-rewards-app/internal/entities/user"
)

var ErrNotEnoughFunds = errors.New("insufficient funds in the account")

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
	SelectAllByUser(ctx context.Context, userID int64) ([]orderEntity.Order, error)
	ValidateOrder(number string) (bool, error)
	UpdateAccruals(cfg *config.Config)
	RequestBalance(ctx context.Context, userID int64) (int, error)
}

type WithdrawalService interface {
	Withdraw(ctx context.Context, order string, sum int, userID int64) error
}

// Service storage of all services.
type Service struct {
	UserService       UserService
	TokenService      TokenService
	OrderService      OrderService
	WithdrawalService WithdrawalService
}
