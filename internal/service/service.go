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

// type OrderService interface {
// 	InsertOrder(ctx context.Context, number string, userID int64) (int64, error)
// 	SelectAllByUser(ctx context.Context, userID int64) ([]orderEntity.Order, error)
// 	ValidateOrder(number string) (bool, error)
// 	UpdateAccruals(cfg *config.Config)
// 	RequestBalance(ctx context.Context, userID int64) (int, error)
// }

// type UserService interface {
// 	Find(ctx context.Context, login string) (*userEntity.User, error)
// 	Create(ctx context.Context, user *userEntity.User) error
// }

// type TokenService interface {
// 	NewToken(ctx context.Context, userID int64, ttl time.Duration) (*userEntity.Token, error)
// 	FindUserByToken(tokenHash []byte) (*userEntity.User, error)
// }

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
