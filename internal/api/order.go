package api

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"net/http"

	"github.com/igortoigildin/go-rewards-app/config"
	orderEntity "github.com/igortoigildin/go-rewards-app/internal/entities/order"
	"github.com/igortoigildin/go-rewards-app/internal/logger"
	"go.uber.org/zap"
)

type OrderService interface {
	InsertOrder(ctx context.Context, number string, userID int64) (int64, error)
	SelectAllByUser(ctx context.Context, userID int64) ([]orderEntity.Order, error)
	ValidateOrder(number string) (bool, error)
	UpdateAccruals(cfg *config.Config)
	RequestBalance(ctx context.Context, userID int64) (int, error)
}

func insertOrderHandler(orderService OrderService) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		user, err := contextGetUser(r)
		if err != nil {
			logger.Log.Info("missing user info:", zap.Error(err))
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		number, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Log.Info("error while reading from reqest body", zap.Error(err))
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		if len(number) == 0 {
			logger.Log.Info("order not provided", zap.Error(err))
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		valid, err := orderService.ValidateOrder(string(number))
		if err != nil || !valid {
			logger.Log.Info("error while validating order:", zap.Error(err))
			rw.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		id, err := orderService.InsertOrder(ctx, string(number), user.ID)
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
	})
}

func allOrdersHandler(orderService OrderService) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		user, err := contextGetUser(r)
		if err != nil {
			logger.Log.Info("missing user info:", zap.Error(err))
			rw.WriteHeader(http.StatusInternalServerError)
		}

		orders, err := orderService.SelectAllByUser(ctx, user.ID)
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

		err = writeJSON(rw, http.StatusOK, orders, nil)
		if err != nil {
			logger.Log.Info("error while encoding response", zap.Error(err))
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}
