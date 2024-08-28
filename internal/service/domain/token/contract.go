package token

import (
	"context"

	userEntity "github.com/igortoigildin/go-rewards-app/internal/entities/user"
)

//go:generate mockgen -package mocks -destination=../../../../mocks/tokenRepository.go github.com/igortoigildin/go-rewards-app/internal/service/domain/token TokenRepository
type TokenRepository interface {
	Insert(ctx context.Context, token *userEntity.Token) error
	FindUserByToken(ctx context.Context, tokenHash []byte) (*userEntity.User, error)
}
