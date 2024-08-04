package storage

import "github.com/igortoigildin/go-rewards-app/internal/models"

type Storage interface {
	CreateUser(user *models.User) error
	GetUserByLogin(login string) (*models.User, error)
}
