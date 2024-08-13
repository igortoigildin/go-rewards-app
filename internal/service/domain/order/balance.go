package order

import (
	"context"
	"database/sql"
	"errors"

	"github.com/igortoigildin/go-rewards-app/internal/entities/order"
	"github.com/igortoigildin/go-rewards-app/internal/logger"
	"go.uber.org/zap"
)

func (o *OrderService) RequestBalance(ctx context.Context, userID int64) (int, error) {
	var orders []order.Order
	var totalAccruals int
	orders, err := o.OrderRepository.SelectAllByUser(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return 0, nil
		default:
			logger.Log.Info("error all orders for userID:", zap.Error(err))
			return 0, err
		}
	}

	for _, order := range orders {
		if order.Accrual != nil {
			totalAccruals += *order.Accrual
		}
	}

	var totalWithdrawals int
	withdrawals, err := o.WithdrawalRepository.SelectAllForUserID(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			totalWithdrawals = 0
		default:
			logger.Log.Info("error all orders for userID:", zap.Error(err))
		}
	}

	for _, val := range withdrawals {
		totalWithdrawals += val.Sum
	}

	return totalAccruals - totalWithdrawals, nil
}
