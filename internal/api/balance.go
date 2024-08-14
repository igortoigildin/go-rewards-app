package api

// func (app *app) balanceHandler(rw http.ResponseWriter, r *http.Request) {
// 	ctx, cancel := context.WithCancel(r.Context())
// 	defer cancel()

// 	user, err := app.contextGetUser(r)
// 	if err != nil {
// 		logger.Log.Info("missing user info:", zap.Error(err))
// 		rw.WriteHeader(http.StatusInternalServerError)
// 	}

// 	_, err = app.services.OrderService.RequestBalance(ctx, user.ID)
// 	if err != nil {
// 		logger.Log.Info("error while obtaining user balance:", zap.Error(err))
// 		rw.WriteHeader(http.StatusInternalServerError)
// 	}

// 	// TODO add JSON reply

// }
