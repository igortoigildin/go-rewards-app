package order

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	model "github.com/igortoigildin/go-rewards-app/internal/entities/order"
	"github.com/igortoigildin/go-rewards-app/mocks"
)

func TestOrderRepository_InsertOrder(t *testing.T) {
	orderNew := model.Order{
		Number: "123",
		UserID: int64(456),
	}
	orderExists := model.Order{
		Number: "120",
		UserID: int64(456),
	}
	blankOrder := model.Order{}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := mocks.NewMockOrderRepository(ctrl)
	r.EXPECT().InsertOrder(context.Background(), &orderNew).Times(1).Return(int64(0), nil)
	r.EXPECT().InsertOrder(context.Background(), &blankOrder).Times(1).Return(int64(0), errors.New("order not provided"))
	r.EXPECT().InsertOrder(context.Background(), &orderExists).Times(1).Return(int64(1), nil)

	tests := []struct {
		name    string
		ctx     context.Context
		order   *model.Order
		want    int64
		wantErr bool
	}{
		{
			name:    "success, order does not exist",
			ctx:     context.Background(),
			order:   &orderNew,
			want:    int64(0),
			wantErr: false,
		},
		{
			name:    "fail, order exists",
			ctx:     context.Background(),
			order:   &orderExists,
			want:    int64(1),
			wantErr: false,
		},
		{
			name:    "fail, error",
			ctx:     context.Background(),
			order:   &blankOrder,
			want:    int64(0),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.InsertOrder(tt.ctx, tt.order)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderRepository.InsertOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("OrderRepository.InsertOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderRepository_SelectAllByUser(t *testing.T) {
	successReply := []model.Order{
		{
			Number: "123",
			Status: "NEW",
			UserID: int64(5),
		},
		{
			Number: "124",
			Status: "PROCESSED",
			UserID: int64(5),
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := mocks.NewMockOrderRepository(ctrl)
	r.EXPECT().SelectAllByUser(context.Background(), int64(5)).Times(1).Return(successReply, nil)
	r.EXPECT().SelectAllByUser(context.Background(), int64(9)).Times(1).Return(nil, errors.New("incorrect userID provided"))

	tests := []struct {
		name    string
		ctx     context.Context
		UserID  int64
		want    []model.Order
		wantErr bool
	}{
		{
			name:    "success",
			ctx:     context.Background(),
			UserID:  int64(5),
			want:    successReply,
			wantErr: false,
		},
		{
			name:    "success",
			ctx:     context.Background(),
			UserID:  int64(9),
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.SelectAllByUser(tt.ctx, tt.UserID)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderRepository.SelectAllByUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OrderRepository.SelectAllByUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
