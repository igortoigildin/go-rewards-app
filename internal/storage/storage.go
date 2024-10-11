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

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		User:       user.NewUserRepository(db),
		Token:      token.NewTokenRepository(db),
		Order:      order.NewOrderRepository(db),
		Withdrawal: withdrawal.NewWithdrawalRepository(db),
	}
}
