package user

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	model "github.com/igortoigildin/go-rewards-app/internal/entities/user"
	"github.com/igortoigildin/go-rewards-app/mocks"
)

func TestUserRepository_Create(t *testing.T) {
	newUser := model.User{
		Login: "dummy_login",
		Password: model.Password{
			Hash: []byte("dsfsd"),
		},
	}
	blank := model.User{}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := mocks.NewMockUserRepository(ctrl)
	r.EXPECT().Create(context.Background(), &newUser).Times(1).Return(nil)
	r.EXPECT().Create(context.Background(), &blank).Times(1).Return(errors.New("user data not provided"))

	tests := []struct {
		name    string
		ctx     context.Context
		user    *model.User
		wantErr bool
	}{
		{
			name:    "success",
			ctx:     context.Background(),
			user:    &newUser,
			wantErr: false,
		},
		{
			name:    "fail",
			ctx:     context.Background(),
			user:    &blank,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := r.Create(tt.ctx, tt.user); (err != nil) != tt.wantErr {
				t.Errorf("UserRepository.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
