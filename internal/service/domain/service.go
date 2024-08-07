package domain

import (
	"github.com/igortoigildin/go-rewards-app/internal/service"
	"github.com/igortoigildin/go-rewards-app/internal/service/domain/token"
	"github.com/igortoigildin/go-rewards-app/internal/service/domain/user"
	"github.com/igortoigildin/go-rewards-app/internal/storage"
)

// NewService implementation for storage of all services.
func NewService(repositories *storage.Repository) *service.Service {
	return &service.Service{
		UserService:  user.NewUserService(repositories.User),
		TokenService: token.NewTokenService(repositories.Token),
	}
}
