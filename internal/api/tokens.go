package api

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/igortoigildin/go-rewards-app/internal/logger"
	"github.com/igortoigildin/go-rewards-app/internal/storage"
	"go.uber.org/zap"
)

func (app *app) createAuthTokenHandler(rw http.ResponseWriter, r *http.Request) {
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

	user, err := app.services.UserService.Find(ctx, input.Login)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrRecordNotFound):
			logger.Log.Info("user not found", zap.Error(err))
			rw.WriteHeader(http.StatusUnauthorized)
			return
		default:
			logger.Log.Info("internal error", zap.Error(err))
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	match, err := user.Password.Matches(input.Password)
	if err != nil {
		logger.Log.Info("internal error", zap.Error(err))
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !match {
		logger.Log.Info("incorrect password", zap.Error(err))
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	token, err := app.services.TokenService.NewToken(ctx, user.ID, 24*time.Hour)
	if err != nil {
		logger.Log.Info("error while ctreating new token", zap.Error(err))
		rw.WriteHeader(http.StatusInternalServerError)
	}
	err = app.writeJSON(rw, http.StatusOK, token, nil)
	if err != nil {
		logger.Log.Info("error while encoding response", zap.Error(err))
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}
