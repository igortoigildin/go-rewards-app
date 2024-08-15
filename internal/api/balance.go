package api

import (
	"context"
	"net/http"

	"github.com/igortoigildin/go-rewards-app/internal/logger"
	"go.uber.org/zap"
)

func balanceHandler(orderService OrderService) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		user, err := contextGetUser(r)
		if err != nil {
			logger.Log.Info("missing user info:", zap.Error(err))
			rw.WriteHeader(http.StatusInternalServerError)
		}

		balance, err := orderService.RequestBalance(ctx, user.UserID)
		if err != nil {
			logger.Log.Info("error while obtaining user balance:", zap.Error(err))
			rw.WriteHeader(http.StatusInternalServerError)
		}
		// TODO withdrawn
	})
}
