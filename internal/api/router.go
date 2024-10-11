package api

import (
	"net/http"

	"github.com/igortoigildin/go-rewards-app/config"
	middleware "github.com/igortoigildin/go-rewards-app/internal/lib/middleware"
	"github.com/igortoigildin/go-rewards-app/internal/service"
)

func Router(s *service.Service, cfg *config.Config) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/user/register", middleware.Timeout(cfg.ContextTimout, registerUserHandler(s.UserService, s.TokenService)))
	mux.HandleFunc("POST /api/user/login", middleware.Timeout(cfg.ContextTimout, createAuthTokenHandler(s.UserService, s.TokenService)))
	mux.HandleFunc("POST /api/user/orders", middleware.Timeout(cfg.ContextTimout, auth(s.TokenService, insertOrderHandler(s.OrderService))))
	mux.HandleFunc("GET /api/user/orders", middleware.Timeout(cfg.ContextTimout, auth(s.TokenService, allOrdersHandler(s.OrderService))))
	mux.HandleFunc("GET /api/user/balance", middleware.Timeout(cfg.ContextTimout, auth(s.TokenService, balanceHandler(s.UserService, s.WithdrawalService))))
	mux.HandleFunc("POST /api/user/balance/withdraw", middleware.Timeout(cfg.ContextTimout, auth(s.TokenService, withdrawHandler(s.WithdrawalService))))
	mux.HandleFunc("GET /api/user/withdrawals", middleware.Timeout(cfg.ContextTimout, auth(s.TokenService, withdrawalsHandler(s.WithdrawalService))))

	return mux
}
