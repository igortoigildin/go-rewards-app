package api

import (
	"crypto/sha256"
	"errors"
	"net/http"

	"github.com/igortoigildin/go-rewards-app/internal/logger"
	"github.com/igortoigildin/go-rewards-app/internal/storage"
	"go.uber.org/zap"
)

func (app *app) auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			switch {
			case errors.Is(err, http.ErrNoCookie):
				logger.Log.Info("cookie not found")
				w.WriteHeader(http.StatusBadRequest)
			default:
				logger.Log.Info("cookies cannot be read")
				w.WriteHeader(http.StatusInternalServerError)
			}
			return 
		}
		plaintext := cookie.Value 
		hash := sha256.Sum256([]byte(plaintext))
		_, err = app.services.TokenService.FindUserByToken(hash[:])
		if err != nil {
			switch {
			case errors.Is(err, storage.ErrRecordNotFound):
				logger.Log.Info("user with such token not found", zap.Error(err))
				w.WriteHeader(http.StatusUnauthorized)
			default:
				logger.Log.Info("error", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
		next.ServeHTTP(w, r)
	})
}