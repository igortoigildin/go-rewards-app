package api

import (
	"crypto/sha256"
	"errors"
	"net/http"

	"github.com/igortoigildin/go-rewards-app/internal/logger"
	"go.uber.org/zap"
)

func auth(tokenService TokenService, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			switch {
			case errors.Is(err, http.ErrNoCookie):
				logger.Log.Info("cookie not found")
				rw.WriteHeader(http.StatusUnauthorized)
			default:
				logger.Log.Info("cookies cannot be read")
				rw.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
		plaintext := cookie.Value
		hash := sha256.Sum256([]byte(plaintext))
		user, err := tokenService.FindUserByToken(hash[:]) // ctx add
		if err != nil {
			switch {
			case errors.Is(err, ErrRecordNotFound):
				logger.Log.Info("user with such token not found", zap.Error(err))
				rw.WriteHeader(http.StatusUnauthorized)
			default:
				logger.Log.Info("error", zap.Error(err))
				rw.WriteHeader(http.StatusInternalServerError)
			}
		}
		r = contextSetUser(r, user)
		next.ServeHTTP(rw, r)
	})
}
