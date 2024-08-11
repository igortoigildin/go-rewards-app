package api

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/igortoigildin/go-rewards-app/config"
	"github.com/igortoigildin/go-rewards-app/internal/logger"
	"github.com/igortoigildin/go-rewards-app/internal/service"
	"go.uber.org/zap"
)

func router(services *service.Service, cfg *config.Config) *http.ServeMux {
	app := newApp(*services, cfg)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/user/register", app.registerUserHandler)
	mux.HandleFunc("POST /api/user/login", app.createAuthTokenHandler)
	mux.HandleFunc("POST /api/user/orders", app.auth(app.insertOrderHandler))
	mux.HandleFunc("GET /api/user/orders", app.auth(app.allOrdersHandler))

	return mux
}

func (app *app) allOrdersHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	user, err := app.contextGetUser(r)
	if err != nil {
		logger.Log.Info("missing user info:", zap.Error(err))
		rw.WriteHeader(http.StatusInternalServerError)
	}

	orders, err := app.services.OrderService.SelectAllByUser(ctx, user.ID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			rw.WriteHeader(http.StatusNoContent)
			return
		default:
			logger.Log.Info("error requesting orders", zap.Error(err))
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	err = app.writeJSON(rw, http.StatusOK, orders, nil)
	if err != nil {
		logger.Log.Info("error while encoding response", zap.Error(err))
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}
