package api

import (
	"context"
	"net/http"

	"github.com/igortoigildin/go-rewards-app/internal/logger"
	"go.uber.org/zap"
)

func balanceHandler(userService UserService, withdrawalService WithdrawalService) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		user, err := contextGetUser(r)
		if err != nil {
			logger.Log.Info("missing user info:", zap.Error(err))
			rw.WriteHeader(http.StatusInternalServerError)
		}

		currentBalance, err := userService.Balance(ctx, user.UserID)
		if err != nil {
			logger.Log.Info("error while obtaining current balance:", zap.Error(err))
			rw.WriteHeader(http.StatusInternalServerError)
		}

		withdraws, err := withdrawalService.WithdrawalsForUser(ctx, user.UserID)
		if err != nil {
			logger.Log.Info("error while obtaining withdrawn balance:", zap.Error(err))
			rw.WriteHeader(http.StatusInternalServerError)
		}

		var withdrawnBalance float64
		for _, withdraw := range withdraws {
			withdrawnBalance += withdraw.Sum
		}

		data := struct {
			Current   float64 `json:"current"`
			Withdrawn float64 `json:"withdrawn"`
		}{
			Current:   float64(currentBalance),
			Withdrawn: float64(withdrawnBalance),
		}

		err = writeJSON(rw, http.StatusOK, data, nil)
		if err != nil {
			logger.Log.Info("error while encoding response:", zap.Error(err))
			rw.WriteHeader(http.StatusInternalServerError)
		}
	})
}
