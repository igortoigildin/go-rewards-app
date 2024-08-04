package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/igortoigildin/go-rewards-app/config"
	"github.com/igortoigildin/go-rewards-app/internal/logger"
	"github.com/igortoigildin/go-rewards-app/internal/models"
	"github.com/igortoigildin/go-rewards-app/internal/storage"
	"go.uber.org/zap"
)

func router(ctx context.Context, cfg *config.Config) *http.ServeMux {
	repo := storage.InitPostgresRepo(ctx, cfg)
	app := newApp(repo, cfg)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/user/register", app.registerUserHandler)

	return mux
}

func (app *app) registerUserHandler(rw http.ResponseWriter, r *http.Request) {
	_, cancel := context.WithCancel(r.Context())
	defer cancel()

	var input struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&input); err != nil {
		logger.Log.Debug("cannot decode request JSON body", zap.Error(err))
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	user := &models.User{
		Login: input.Login,
	}

	err := user.Password.Set(input.Password)
	if err != nil {
		logger.Log.Info("error while setting password", zap.Error(err))
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = app.storage.CreateUser(user)
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

	enc := json.NewEncoder(rw)
	if err := enc.Encode(user); err != nil {
		logger.Log.Info("error while encoding response", zap.Error(err))
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	logger.Log.Info("sending HTTP 200 response")
}
