package api

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	processJSON "github.com/igortoigildin/go-rewards-app/internal/lib/processJSON"
	"github.com/igortoigildin/go-rewards-app/internal/logger"
	"go.uber.org/zap"
)

func createAuthTokenHandler(userService UserService, tokenService TokenService) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var input struct {
			Login    string `json:"login"`
			Password string `json:"password"`
		}

		err := processJSON.ReadJSON(r, &input)
		if err != nil {
			logger.Log.Info("cannot decode request JSON body", zap.Error(err))
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		user, err := userService.Find(ctx, input.Login)
		if err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				logger.Log.Info("user not found", zap.Error(err))
				rw.WriteHeader(http.StatusUnauthorized)
				return
			default:
				logger.Log.Info("internal error", zap.Error(err))
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		match, err := user.Password.Matches(input.Password)
		if err != nil {
			logger.Log.Info("internal error", zap.Error(err))
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !match {
			logger.Log.Info("incorrect password", zap.Error(err))
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}

		token, err := tokenService.NewToken(ctx, user.UserID, 24*time.Hour)
		if err != nil {
			logger.Log.Info("error while ctreating new token", zap.Error(err))
			rw.WriteHeader(http.StatusInternalServerError)
		}

		// Initialize a new cookie containing the new token create
		cookie := http.Cookie{
			Name:     "token",
			Value:    token.Plaintext,
			Expires:  token.Expiry,
			HttpOnly: true,
		}
		http.SetCookie(rw, &cookie)

		err = processJSON.WriteJSON(rw, http.StatusOK, token, nil)
		if err != nil {
			logger.Log.Info("error while encoding response", zap.Error(err))
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}
