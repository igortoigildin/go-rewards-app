package withdrawal

import (
	"context"

	withdrawalEntity "github.com/igortoigildin/go-rewards-app/internal/entities/withdrawal"
)

//go:generate mockgen -package mocks -destination=../../../../mocks/withdrawalRepository.go github.com/igortoigildin/go-rewards-app/internal/service/domain/withdrawal WithdrawalRepository
type WithdrawalRepository interface {
	Create(ctx context.Context, withdrawal *withdrawalEntity.Withdrawal) error
	SelectAllForUserID(ctx context.Context, userID int64) ([]withdrawalEntity.Withdrawal, error)
}
