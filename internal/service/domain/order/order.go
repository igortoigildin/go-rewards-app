package order

import (
	"context"

	orderEntity "github.com/igortoigildin/go-rewards-app/internal/entities/order"
)

const (
	statusNew        = "NEW"
	statusRegistered = "REGISTERED"
)

type OrderService struct {
	OrderRepository OrderRepository
}

func NewOrderService(OrderRepository OrderRepository) *OrderService {
	return &OrderService{
		OrderRepository: OrderRepository,
	}
}

// Returns -1 in case of success or returns user id who already added this order.
func (o *OrderService) InsertOrder(ctx context.Context, number string, userID int64) (int64, error) {
	order := orderEntity.Order{
		Number: number,
		Status: statusNew,
		UserID: userID,
	}
	id, err := o.OrderRepository.InsertOrder(ctx, &order)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (o *OrderService) SelectAllByUser(ctx context.Context, userID int64) ([]orderEntity.Order, error) {
	orders, err := o.OrderRepository.SelectAllByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return orders, nil
}
