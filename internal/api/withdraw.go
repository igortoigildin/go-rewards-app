package api

// func (app *app) withdrawHandler(rw http.ResponseWriter, r *http.Request) {
// 	ctx, cancel := context.WithCancel(r.Context())
// 	defer cancel()

// 	user, err := app.contextGetUser(r)
// 	if err != nil {
// 		logger.Log.Info("missing user info:", zap.Error(err))
// 		rw.WriteHeader(http.StatusInternalServerError)
// 	}

// 	order := struct {
// 		order string
// 		sum   int
// 	}{}

// 	err = app.readJSON(r, order)
// 	if err != nil {
// 		logger.Log.Info("missing user info:", zap.Error(err))
// 		rw.WriteHeader(http.StatusInternalServerError)
// 	}

// 	err = app.services.WithdrawalService.Withdraw(ctx, order.order, order.sum, user.ID)
// 	if err != nil {
// 		switch {
// 		case errors.Is(err, service.ErrNotEnoughFunds):
// 			logger.Log.Info("not enough funds:", zap.Error(err))
// 			rw.WriteHeader(http.StatusPaymentRequired)
// 		default:
// 			logger.Log.Info("error while making withdrawal:", zap.Error(err))
// 		}
// 	}
// }
