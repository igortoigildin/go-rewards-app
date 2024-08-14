package api

import (
	"net/http"

	"github.com/igortoigildin/go-rewards-app/internal/service"
)

// func Router() *http.ServeMux {

// mux := http.NewServeMux()
// mux.HandleFunc("POST /api/user/register", registerUserHandler()   )
// mux.HandleFunc("POST /api/user/login", app.createAuthTokenHandler)
// mux.HandleFunc("POST /api/user/orders", app.auth(app.insertOrderHandler))
// mux.HandleFunc("GET /api/user/orders", app.auth(app.allOrdersHandler))
// mux.HandleFunc("GET /api/user/balance", app.auth(app.balanceHandler))
// mux.HandleFunc("POST /api/user/balance/withdraw", app.auth(app.withdrawHandler))

// return mux

// TODO: add handler for /api/user/withdrawals

func Router(services *service.Service) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/user/register", registerUserHandler(services.UserService, services.TokenService))
	mux.HandleFunc("POST /api/user/login", createAuthTokenHandler(services.UserService, services.TokenService))
	mux.HandleFunc("POST /api/user/orders", auth(services.TokenService, insertOrderHandler(services.OrderService)))
	mux.HandleFunc("GET /api/user/orders", auth(services.TokenService, allOrdersHandler(services.OrderService)))

	return mux
}
