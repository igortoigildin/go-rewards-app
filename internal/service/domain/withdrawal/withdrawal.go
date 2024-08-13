package withdrawal

import (
	"context"
	"time"

	withdrawalEntity "github.com/igortoigildin/go-rewards-app/internal/entities/withdrawal"
	"github.com/igortoigildin/go-rewards-app/internal/logger"
	"github.com/igortoigildin/go-rewards-app/internal/service"
	"github.com/igortoigildin/go-rewards-app/internal/storage"
	"go.uber.org/zap"
)

type WithdrawalService struct {
	WithdrawalRepository storage.WithdrawalRepository
	OrderService         service.OrderService
}

func (w *WithdrawalService) Withdraw(ctx context.Context, order string, sum int, userID int64) error {
	balance, err := w.OrderService.RequestBalance(ctx, userID)
	if err != nil {
		logger.Log.Info("error while obtaining user balance:", zap.Error(err))
	}

	if balance < sum {
		return service.ErrNotEnoughFunds
	}

	t, err := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	if err != nil {
		logger.Log.Info("time parsing error", zap.Error(err))
	}

	withdrawal := withdrawalEntity.Withdrawal{
		Order:  order,
		Sum:    sum,
		UserID: userID,
		Date:   t,
	}
	// add check if number of order is valid - 422 error

	err = w.WithdrawalRepository.Create(ctx, &withdrawal)
	if err != nil {
		logger.Log.Info("error while recoding withdrawal", zap.Error(err))
		return err
	}
	return nil
}

// NewWithdrawalService returns a new instance of user service.
func NewWithdrawalService(WithdrawalRepository storage.WithdrawalRepository) *WithdrawalService {
	return &WithdrawalService{
		WithdrawalRepository: WithdrawalRepository,
	}
}
