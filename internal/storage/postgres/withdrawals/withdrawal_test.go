package withdrawal

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	model "github.com/igortoigildin/go-rewards-app/internal/entities/withdrawal"
	"github.com/igortoigildin/go-rewards-app/mocks"
)

func TestWithdrawalRepository_Create(t *testing.T) {
	withdrawal := model.Withdrawal{
		Order:  "12345",
		Sum:    float64(50),
		UserID: int64(45),
	}
	blank := model.Withdrawal{}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := mocks.NewMockWithdrawalRepository(ctrl)
	r.EXPECT().Create(context.Background(), &withdrawal).Times(1).Return(nil)
	r.EXPECT().Create(context.Background(), &blank).Times(1).Return(errors.New("withdrawal data not provided"))

	tests := []struct {
		name       string
		ctx        context.Context
		withdrawal *model.Withdrawal
		wantErr    bool
	}{
		{
			name:       "success",
			ctx:        context.Background(),
			withdrawal: &withdrawal,
			wantErr:    false,
		},
		{
			name:       "fail",
			ctx:        context.Background(),
			withdrawal: &blank,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := r.Create(tt.ctx, tt.withdrawal); (err != nil) != tt.wantErr {
				t.Errorf("WithdrawalRepository.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
