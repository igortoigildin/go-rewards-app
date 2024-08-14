package storage

import (
	"database/sql"

	order "github.com/igortoigildin/go-rewards-app/internal/storage/postgres/orders"
	"github.com/igortoigildin/go-rewards-app/internal/storage/postgres/token"
	"github.com/igortoigildin/go-rewards-app/internal/storage/postgres/user"
	withdrawal "github.com/igortoigildin/go-rewards-app/internal/storage/postgres/withdrawals"
)

type Repository struct {
	User       *user.UserRepository
	Token      *token.TokenRepository
	Order      *order.OrderRepository
	Withdrawal *withdrawal.WithdrawalRepository
}

func NewRepository(DB *sql.DB) *Repository {
	return &Repository{
		User:       user.NewUserRepository(DB),
		Token:      token.NewTokenRepository(DB),
		Order:      order.NewOrderRepository(DB),
		Withdrawal: withdrawal.NewWithdrawalRepository(DB),
	}
}
