package order

import (
	"context"

	"github.com/igortoigildin/go-rewards-app/internal/entities/order"
	orderEntity "github.com/igortoigildin/go-rewards-app/internal/entities/order"
)

//go:generate mockgen -package mocks -destination=../../../../mocks/orderRepository.go github.com/igortoigildin/go-rewards-app/internal/service/domain/order OrderRepository
type OrderRepository interface {
	InsertOrder(ctx context.Context, order *orderEntity.Order) (int64, error)
	SelectAllByUser(ctx context.Context, user int64) ([]orderEntity.Order, error)
	SelectForAccrualCalc() ([]order.Order, error) 
	UpdateOrderAndBalance(ctx context.Context, order *orderEntity.Order) error
}
