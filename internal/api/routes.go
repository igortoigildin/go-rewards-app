package api

import (
	"net/http"

	"github.com/igortoigildin/go-rewards-app/config"
	"github.com/igortoigildin/go-rewards-app/internal/service"
)

func router(services *service.Service, cfg *config.Config) *http.ServeMux {
	app := newApp(*services, cfg)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/user/register", app.registerUserHandler)
	mux.HandleFunc("POST /api/user/login", app.createAuthTokenHandler)
	mux.HandleFunc("POST /api/user/orders", app.auth(app.loadOrderHandler))

	return mux
}
