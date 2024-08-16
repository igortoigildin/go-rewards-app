package withdrawal

import (
	"context"

	model "github.com/igortoigildin/go-rewards-app/internal/entities/withdrawal"
	"github.com/igortoigildin/go-rewards-app/internal/logger"
	"go.uber.org/zap"
)

type WithdrawalService struct {
	WithdrawalRepository WithdrawalRepository
}

func NewWithdrawalService(WithdrawalRepository WithdrawalRepository) *WithdrawalService {
	return &WithdrawalService{
		WithdrawalRepository: WithdrawalRepository,
	}
}

func (w *WithdrawalService) Withdraw(ctx context.Context, order string, sum int, userID int64) error {
	withdrawal := model.Withdrawal{
		Order:  order,
		Sum:    sum,
		UserID: userID,
	}

	err := w.WithdrawalRepository.Create(ctx, &withdrawal)
	if err != nil {
		logger.Log.Info("error during withdrawal", zap.Error(err))
		return err
	}
	return nil
}

func (w *WithdrawalService) WithdrawalsForUser(ctx context.Context, userID int64) ([]model.Withdrawal, error) {
	trans, err := w.WithdrawalRepository.SelectAllForUserID(ctx, userID)
	if err != nil {
		logger.Log.Info("error loading all withdrawals for user", zap.Error(err))
		return nil, err
	}
	return trans, nil
}
