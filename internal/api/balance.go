package api

import (
	"net/http"

	ctxPac "github.com/igortoigildin/go-rewards-app/internal/lib/context"
	processJSON "github.com/igortoigildin/go-rewards-app/internal/lib/processJSON"
	"github.com/igortoigildin/go-rewards-app/internal/logger"
	"go.uber.org/zap"
)

func balanceHandler(userService UserService, withdrawalService WithdrawalService) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		user, err := ctxPac.ContextGetUser(r)
		if err != nil {
			logger.Log.Info("missing user info:", zap.Error(err))
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		currentBalance, err := userService.Balance(ctx, user.UserID)
		if err != nil {
			logger.Log.Info("error while obtaining current balance:", zap.Error(err))
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		withdraws, err := withdrawalService.WithdrawalsForUser(ctx, user.UserID)
		if err != nil {
			logger.Log.Info("error while obtaining withdrawn balance:", zap.Error(err))
			rw.WriteHeader(http.StatusInternalServerError)
			return
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

		err = processJSON.WriteJSON(rw, http.StatusOK, data, nil)
		if err != nil {
			logger.Log.Info("error while encoding response:", zap.Error(err))
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}
