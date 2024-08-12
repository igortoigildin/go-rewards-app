package api

import (
	"context"
	"errors"
	"net/http"
	"time"

	userEntity "github.com/igortoigildin/go-rewards-app/internal/entities/user"
	"github.com/igortoigildin/go-rewards-app/internal/logger"
	"github.com/igortoigildin/go-rewards-app/internal/storage"
	"go.uber.org/zap"
)

func (app *app) registerUserHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	var input struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	err := app.readJSON(r, &input)
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

	err = app.services.UserService.Create(ctx, user)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrDuplicateLogin):
			logger.Log.Info("a user with this login already exists", zap.Error(err))
			rw.WriteHeader(http.StatusConflict)
		default:
			logger.Log.Info("error while saving user", zap.Error(err))
			rw.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// Create new auth token for new registered user
	token, err := app.services.TokenService.NewToken(ctx, user.ID, 24*time.Hour)
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
}
