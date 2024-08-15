package withdrawal

import "context"

type WithdrawalService struct {
	WithdrawalRepository WithdrawalRepository
}

func NewWithdrawalService(WithdrawalRepository WithdrawalRepository) *WithdrawalService {
	return &WithdrawalService{
		WithdrawalRepository: WithdrawalRepository,
	}
}

func (w *WithdrawalService) Withdraw(ctx context.Context, order string, sum int, userID int64) error {
	var order 
	
	
	
	w.WithdrawalRepository.Create(ctx, )








}
