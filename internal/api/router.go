package api

import (
	"net/http"

	"github.com/igortoigildin/go-rewards-app/config"
	"github.com/igortoigildin/go-rewards-app/internal/service"
)

func Router(s *service.Service, cfg *config.Config) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/user/register", registerUserHandler(s.UserService, s.TokenService))
	mux.HandleFunc("POST /api/user/login", createAuthTokenHandler(s.UserService, s.TokenService))
	mux.HandleFunc("POST /api/user/orders", auth(s.TokenService, insertOrderHandler(s.OrderService)))
	mux.HandleFunc("GET /api/user/orders", auth(s.TokenService, allOrdersHandler(s.OrderService)))
	mux.HandleFunc("GET /api/user/balance", auth(s.TokenService, balanceHandler(s.UserService, s.WithdrawalService)))
	mux.HandleFunc("POST /api/user/balance/withdraw", auth(s.TokenService, withdrawHandler(s.WithdrawalService)))
	mux.HandleFunc("GET /api/user/withdrawals", auth(s.TokenService, withdrawalsHandler(s.WithdrawalService)))

	return mux
}
