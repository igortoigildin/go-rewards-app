package order

import (
	"context"

	orderEntity "github.com/igortoigildin/go-rewards-app/internal/entities/order"
)

type OrderRepository interface {
	InsertOrder(ctx context.Context, order *orderEntity.Order) (int64, error)
	SelectAllByUser(ctx context.Context, user int64) ([]orderEntity.Order, error)
	SelectForAccrualCalc() ([]int64, error)
	Update(ctx context.Context, order *orderEntity.Order) error
}
