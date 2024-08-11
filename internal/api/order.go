package api

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"net/http"

	"github.com/igortoigildin/go-rewards-app/internal/logger"
	"go.uber.org/zap"
)

func (app *app) insertOrderHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	user, err := app.contextGetUser(r)
	if err != nil {
		logger.Log.Info("missing user info:", zap.Error(err))
		rw.WriteHeader(http.StatusInternalServerError)
	}

	number, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Log.Info("error while reading from reqest body", zap.Error(err))
		rw.WriteHeader(http.StatusBadRequest)
	}

	valid, err := app.services.OrderService.ValidateOrder(string(number))
	if err != nil || !valid {
		logger.Log.Info("error while validating order:", zap.Error(err))
		rw.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	id, err := app.services.OrderService.InsertOrder(ctx, string(number), user.ID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			logger.Log.Info("new order accepted successfully")
			rw.WriteHeader(http.StatusAccepted)
			return
		default:
			logger.Log.Info("error while inserting order", zap.Error(err))
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	switch {
	case id == user.ID:
		logger.Log.Info("this order already added by this user")
		rw.WriteHeader(http.StatusOK)
		return
	case id == -1:
		logger.Log.Info("order added successfully")
		rw.WriteHeader(http.StatusOK)
		return
	default:
		logger.Log.Info("this order already added by another user")
		rw.WriteHeader(http.StatusConflict)
		return
	}
}
