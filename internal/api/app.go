package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/igortoigildin/go-rewards-app/config"
	entities "github.com/igortoigildin/go-rewards-app/internal/entities/user"
	"github.com/igortoigildin/go-rewards-app/internal/logger"
	"github.com/igortoigildin/go-rewards-app/internal/service"
	"github.com/igortoigildin/go-rewards-app/internal/storage"
	"go.uber.org/zap"
)

type app struct {
	services service.Service
	cfg      *config.Config
}

func newApp(service service.Service, cfg *config.Config) *app {
	return &app{
		services: service,
		cfg:      cfg,
	}
}

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

	user := &entities.User{
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

	err = app.writeJSON(rw, http.StatusOK, user, nil)
	if err != nil {
		logger.Log.Info("error while encoding response", zap.Error(err))
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}

