package api

import (
	"net/http"
)

func router(app *app) *http.ServeMux {

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/user/register", app.registerUserHandler)
	mux.HandleFunc("POST /api/user/login", app.createAuthTokenHandler)
	mux.HandleFunc("POST /api/user/orders", app.auth(app.insertOrderHandler))
	mux.HandleFunc("GET /api/user/orders", app.auth(app.allOrdersHandler))

	return mux
}
