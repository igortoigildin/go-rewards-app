package withdrawal

type WithdrawalService struct {
	WithdrawalRepository WithdrawalRepository
}

func NewWithdrawalService(WithdrawalRepository WithdrawalRepository) *WithdrawalService {
	return &WithdrawalService{
		WithdrawalRepository: WithdrawalRepository,
	}
}
