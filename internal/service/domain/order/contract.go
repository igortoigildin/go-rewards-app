package order

import (
	"context"

	orderEntity "github.com/igortoigildin/go-rewards-app/internal/entities/order"
)

//go:generate mockgen -package mocks -destination=../../../../mocks/orderRepository.go github.com/igortoigildin/go-rewards-app/internal/service/domain/order OrderRepository
type OrderRepository interface {
	InsertOrder(ctx context.Context, order *orderEntity.Order) (int64, error)
	SelectAllByUser(ctx context.Context, user int64) ([]orderEntity.Order, error)
	SelectForAccrualCalc(ctx context.Context) ([]orderEntity.Order, error)
	UpdateOrderAndBalance(ctx context.Context, order *orderEntity.Order) error
}
