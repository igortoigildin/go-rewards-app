package api

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	model "github.com/igortoigildin/go-rewards-app/internal/entities/withdrawal"
	"github.com/igortoigildin/go-rewards-app/internal/logger"
	"github.com/igortoigildin/go-rewards-app/internal/service"
	"go.uber.org/zap"
)

//go:generate mockgen -package mocks -destination=../../mocks/withdrawalService.go github.com/igortoigildin/go-rewards-app/internal/api WithdrawalService
type WithdrawalService interface {
	Withdraw(ctx context.Context, order string, sum float64, userID int64) error
	WithdrawalsForUser(ctx context.Context, userID int64) ([]model.Withdrawal, error)
}

func withdrawHandler(withdrawalService WithdrawalService) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		user, err := contextGetUser(r)
		if err != nil {
			logger.Log.Info("missing user info:", zap.Error(err))
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		order := struct {
			Order string  `json:"order"`
			Sum   float64 `json:"sum"`
		}{}

		err = readJSON(r, &order)
		if err != nil {
			logger.Log.Info("error while decoding json:", zap.Error(err))
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		valid, err := ValidateOrder(order.Order)
		if err != nil || !valid {
			logger.Log.Info("error while validating order:", zap.Error(err))
			rw.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		err = withdrawalService.Withdraw(ctx, order.Order, order.Sum, user.UserID)
		if err != nil {
			switch {
			case errors.Is(err, service.ErrNotEnoughFunds):
				logger.Log.Info("not enough funds:", zap.Error(err))
				rw.WriteHeader(http.StatusPaymentRequired)
			default:
				logger.Log.Info("error while making withdrawal:", zap.Error(err))
			}
			return
		}
	}
}

func withdrawalsHandler(withdrawalService WithdrawalService) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		user, err := contextGetUser(r)
		if err != nil {
			logger.Log.Info("missing user info:", zap.Error(err))
			rw.WriteHeader(http.StatusInternalServerError)
		}

		trans, err := withdrawalService.WithdrawalsForUser(ctx, user.UserID)
		if err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				rw.WriteHeader(http.StatusNoContent)
				return
			default:
				logger.Log.Info("error requesting withdrawals", zap.Error(err))
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		err = writeJSON(rw, http.StatusOK, trans, nil)
		if err != nil {
			logger.Log.Info("error while encoding response", zap.Error(err))
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
