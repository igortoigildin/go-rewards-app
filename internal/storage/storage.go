package storage

import (
	"context"
	"errors"

	orderEntity "github.com/igortoigildin/go-rewards-app/internal/entities/order"
	userEntity "github.com/igortoigildin/go-rewards-app/internal/entities/user"
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
}

// Repository storage of all repositories.
type Repository struct {
	User  UserRepository
	Token TokenRepository
	Order OrderRepository
}
