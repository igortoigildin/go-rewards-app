package api

import (
	"context"
	"errors"
	"net/http"
	"time"

	userEntity "github.com/igortoigildin/go-rewards-app/internal/entities/user"
	"github.com/igortoigildin/go-rewards-app/internal/logger"
	"go.uber.org/zap"
)

var ErrDuplicateLogin = errors.New("duplicate login")

type UserService interface {
	Find(ctx context.Context, login string) (*userEntity.User, error)
	Create(ctx context.Context, user *userEntity.User) error
}

type TokenService interface {
	NewToken(ctx context.Context, userID int64, ttl time.Duration) (*userEntity.Token, error)
	FindUserByToken(tokenHash []byte) (*userEntity.User, error)
}

func registerUserHandler(userService UserService, tokenService TokenService) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		var input struct {
			Login    string `json:"login"`
			Password string `json:"password"`
		}
		err := readJSON(r, &input)
		if err != nil {
			logger.Log.Info("cannot decode request JSON body", zap.Error(err))
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		user := &userEntity.User{
			Login: input.Login,
		}

		err = user.Password.Set(input.Password)
		if err != nil {
			logger.Log.Info("error while setting password", zap.Error(err))
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = userService.Create(ctx, user)
		if err != nil {
			switch {
			case errors.Is(err, ErrDuplicateLogin):
				logger.Log.Info("a user with this login already exists", zap.Error(err))
				rw.WriteHeader(http.StatusConflict)
			default:
				logger.Log.Info("error while saving user", zap.Error(err))
				rw.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		// Create new auth token for new registered user
		token, err := tokenService.NewToken(ctx, user.ID, 24*time.Hour)
		if err != nil {
			logger.Log.Info("error while ctreating new token", zap.Error(err))
			rw.WriteHeader(http.StatusInternalServerError)
		}

		// Initialize a new cookie containing the new token create
		cookie := http.Cookie{
			Name:     "token",
			Value:    token.Plaintext,
			Expires:  token.Expiry,
			HttpOnly: true,
		}
		http.SetCookie(rw, &cookie)
		rw.WriteHeader(http.StatusOK)
	})
}
