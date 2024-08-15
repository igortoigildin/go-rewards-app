package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/igortoigildin/go-rewards-app/internal/logger"
	"github.com/igortoigildin/go-rewards-app/internal/service"
	"go.uber.org/zap"
)

// func (app *app) withdrawHandler(rw http.ResponseWriter, r *http.Request) {
// 	ctx, cancel := context.WithCancel(r.Context())
// 	defer cancel()

// 	user, err := app.contextGetUser(r)
// 	if err != nil {
// 		logger.Log.Info("missing user info:", zap.Error(err))
// 		rw.WriteHeader(http.StatusInternalServerError)
// 	}

// 	order := struct {
// 		order string
// 		sum   int
// 	}{}

// 	err = app.readJSON(r, order)
// 	if err != nil {
// 		logger.Log.Info("missing user info:", zap.Error(err))
// 		rw.WriteHeader(http.StatusInternalServerError)
// 	}

// 	err = app.services.WithdrawalService.Withdraw(ctx, order.order, order.sum, user.ID)
// 	if err != nil {
// 		switch {
// 		case errors.Is(err, service.ErrNotEnoughFunds):
// 			logger.Log.Info("not enough funds:", zap.Error(err))
// 			rw.WriteHeader(http.StatusPaymentRequired)
// 		default:
// 			logger.Log.Info("error while making withdrawal:", zap.Error(err))
// 		}
// 	}
// }

type WithdrawalService interface {
	Withdraw(ctx context.Context, order string, sum int, userID int64) error
}

func withdrawalHandler(withdrawalService WithdrawalService) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		user, err := contextGetUser(r)
		if err != nil {
			logger.Log.Info("missing user info:", zap.Error(err))
			rw.WriteHeader(http.StatusInternalServerError)
		}

		order := struct {
			order string
			sum   int
		}{}

		err = readJSON(r, order)
		if err != nil {
			logger.Log.Info("missing user info:", zap.Error(err))
			rw.WriteHeader(http.StatusInternalServerError)
		}

		valid, err := ValidateOrder(order.order)
		if err != nil || !valid {
			logger.Log.Info("error while validating order:", zap.Error(err))
			rw.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		err = withdrawalService.Withdraw(ctx, order.order, order.sum, user.UserID)
		if err != nil {
			switch {
			case errors.Is(err, service.ErrNotEnoughFunds):
				logger.Log.Info("not enough funds:", zap.Error(err))
				rw.WriteHeader(http.StatusPaymentRequired)
			default:
				logger.Log.Info("error while making withdrawal:", zap.Error(err))
			}
		}

	}
}
