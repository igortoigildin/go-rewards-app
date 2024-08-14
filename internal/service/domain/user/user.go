package user

import (
	"context"

	userEntity "github.com/igortoigildin/go-rewards-app/internal/entities/user"
	"github.com/igortoigildin/go-rewards-app/internal/logger"
	"go.uber.org/zap"
)

type UserService struct {
	UserRepository UserRepository
}

func NewUserService(UserRepository UserRepository) *UserService {
	return &UserService{
		UserRepository: UserRepository,
	}
}

func (u *UserService) Create(ctx context.Context, user *userEntity.User) error {
	if err := u.UserRepository.Create(ctx, user); err != nil {
		logger.Log.Info("failed to create user", zap.Error(err))
		return err
	}
	return nil
}

func (u *UserService) Find(ctx context.Context, login string) (*userEntity.User, error) {
	user, err := u.UserRepository.Find(ctx, login)
	if err != nil {
		logger.Log.Info("failed to find user", zap.Error(err))
		return nil, err
	}
	return user, nil
}
