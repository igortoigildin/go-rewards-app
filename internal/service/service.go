package service

import (
	"errors"

	"github.com/igortoigildin/go-rewards-app/internal/service/domain/order"
	"github.com/igortoigildin/go-rewards-app/internal/service/domain/token"
	"github.com/igortoigildin/go-rewards-app/internal/service/domain/user"
	"github.com/igortoigildin/go-rewards-app/internal/service/domain/withdrawal"
	"github.com/igortoigildin/go-rewards-app/internal/storage"
)

var ErrNotEnoughFunds = errors.New("insufficient funds in the account")

// Service storage of all services.
type Service struct {
	UserService       *user.UserService
	TokenService      *token.TokenService
	OrderService      *order.OrderService
	WithdrawalService *withdrawal.WithdrawalService
}

// NewService implementation for storage of all services.
func NewService(repositories *storage.Repository) *Service {
	return &Service{
		UserService:       user.NewUserService(repositories.User),
		TokenService:      token.NewTokenService(repositories.Token),
		OrderService:      order.NewOrderService(repositories.Order),
		WithdrawalService: withdrawal.NewWithdrawalService(repositories.Withdrawal),
	}
}
