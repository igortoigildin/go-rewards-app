package storage

import (
	"context"
	"errors"

	orderEntity "github.com/igortoigildin/go-rewards-app/internal/entities/order"
	userEntity "github.com/igortoigildin/go-rewards-app/internal/entities/user"
	withdrawalEntity "github.com/igortoigildin/go-rewards-app/internal/entities/withdrawal"
)

var (
	ErrDuplicateLogin = errors.New("duplicate login")
	ErrRecordNotFound = errors.New("no records found")
)

type UserRepository interface {
	Create(ctx context.Context, user *userEntity.User) error
	Find(ctx context.Context, login string) (*userEntity.User, error)
}

type TokenRepository interface {
	Insert(ctx context.Context, token *userEntity.Token) error
	FindUserByToken(tokenHash []byte) (*userEntity.User, error)
}

type OrderRepository interface {
	InsertOrder(ctx context.Context, order *orderEntity.Order) (int64, error)
	SelectAllByUser(ctx context.Context, user int64) ([]orderEntity.Order, error)
	SelectForAccrualCalc() ([]int64, error)
	Update(order *orderEntity.Order) error
}

type WithdrawalRepository interface {
	Create(ctx context.Context, withdrawal *withdrawalEntity.Withdrawal) error
	SelectAllForUserID(ctx context.Context, userID int64) ([]withdrawalEntity.Withdrawal, error)
}

// Repository storage of all repositories.
type Repository struct {
	User       UserRepository
	Token      TokenRepository
	Order      OrderRepository
	Withdrawal WithdrawalRepository
}
