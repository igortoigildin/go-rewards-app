package timeout

import (
	"context"
	"net/http"
	"time"

	"github.com/igortoigildin/go-rewards-app/internal/logger"
)

// Timeout middleware.
func Timeout(timeout time.Duration, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		r = r.WithContext(ctx)

		processDone := make(chan bool)
		go func() {
			next.ServeHTTP(rw, r)
			processDone <- true
		}()

		select {
		case <-ctx.Done():
			logger.Log.Info("HTTP Request timed out")
			rw.WriteHeader(http.StatusRequestTimeout)
			rw.Write([]byte("timed out"))
		case <-processDone:
		}
	})
}
