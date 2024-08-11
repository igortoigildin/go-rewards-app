package order

import (
	"context"
	"strconv"
	"time"

	orderEntity "github.com/igortoigildin/go-rewards-app/internal/entities/order"
	"github.com/igortoigildin/go-rewards-app/internal/logger"
	"github.com/igortoigildin/go-rewards-app/internal/storage"
	"go.uber.org/zap"
)

const newOrder = "NEW"

type OrderService struct {
	OrderRepository storage.OrderRepository
}

// Returns -1 in case of success or returns user id who already added this order.
func (o *OrderService) InsertOrder(ctx context.Context, number string, userID int64) (int64, error) {
	t, err := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	if err != nil {
		logger.Log.Info("time parsing error", zap.Error(err))
		return 0, err
	}

	order := orderEntity.Order{
		Number:      number,
		Status:      newOrder,
		Uploaded_at: t,
		UserID:      userID,
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

func (o *OrderService) ValidateOrder(number string) (bool, error) {
	res, err := strconv.Atoi(number)
	if err != nil {
		logger.Log.Info("error while converting number", zap.Error(err))
		return false, err
	}
	return Valid(res), nil
}

// Valid check number is valid or not based on Luhn algorithm
func Valid(number int) bool {
	return (number%10+checksum(number/10))%10 == 0
}

func checksum(number int) int {
	var luhn int

	for i := 0; number > 0; i++ {
		cur := number % 10

		if i%2 == 0 { // even
			cur = cur * 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}

		luhn += cur
		number = number / 10
	}
	return luhn % 10
}

// NewOrderService returns a new instance of order service.
func NewOrderService(OrderRepository storage.OrderRepository) *OrderService {
	return &OrderService{
		OrderRepository: OrderRepository,
	}
}
