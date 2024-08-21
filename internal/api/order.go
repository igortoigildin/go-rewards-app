package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/igortoigildin/go-rewards-app/config"
	orderEntity "github.com/igortoigildin/go-rewards-app/internal/entities/order"
	"github.com/igortoigildin/go-rewards-app/internal/logger"
	"go.uber.org/zap"
)

//go:generate mockgen -package mocks -destination=../../mocks/orderService.go github.com/igortoigildin/go-rewards-app/internal/api OrderService
type OrderService interface {
	InsertOrder(ctx context.Context, number string, userID int64) (int64, error)
	SelectAllByUser(ctx context.Context, userID int64) ([]orderEntity.Order, error)
	UpdateAccruals(cfg *config.Config, order *orderEntity.Order)
}

func insertOrderHandler(orderService OrderService, cfg *config.Config) http.HandlerFunc {
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
			rw.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		valid, err := ValidateOrder(string(number))
		if err != nil || !valid {
			logger.Log.Info("error while validating order:", zap.Error(err))
			rw.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		id, err := orderService.InsertOrder(ctx, string(number), user.UserID)
		if err != nil {
			logger.Log.Info("error while inserting order", zap.Error(err))
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		var accraualValue float64
		order := orderEntity.Order{
			Number:  string(number),
			Status:  "NEW",
			Accrual: &accraualValue,
			UserID:  user.UserID,
		}

		switch {
		case id == user.UserID:
			logger.Log.Info("this order already added by this user")
			rw.WriteHeader(http.StatusOK)
			return
		case id == 0:
			logger.Log.Info("order added successfully")
			rw.WriteHeader(http.StatusAccepted)
			go orderService.UpdateAccruals(cfg, &order) // Send reqest to accrual api
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

		orders, err := orderService.SelectAllByUser(ctx, user.UserID)
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

		js, err := json.Marshal(orders)
		if err != nil {
			logger.Log.Info("error while marshalling", zap.Error(err))
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write(js)
	})
}
