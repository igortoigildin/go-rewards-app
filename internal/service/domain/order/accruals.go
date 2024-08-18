package order

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/igortoigildin/go-rewards-app/config"
	orderEntity "github.com/igortoigildin/go-rewards-app/internal/entities/order"
	"github.com/igortoigildin/go-rewards-app/internal/logger"
	"go.uber.org/zap"
)

var statusInvalid string = "INVALID"
var statusProcessed string = "PROCESSED"

func (o *OrderService) UpdateAccruals(cfg *config.Config, order *orderEntity.Order) {

	for order.Status != statusInvalid && order.Status != statusProcessed {
		url := cfg.FlagAccSysAddr + fmt.Sprintf("/api/orders/%v", order.Number)
		resp, err := http.Get(url)
		if err != nil {
			logger.Log.Info("error while reaching accrual system", zap.Error(err))
		}

		newOrder := struct {
			Order   string  `json:"number"`
			Status  string  `json:"status"`
			Accrual float64 `json:"accrual,omitempty"`
		}{}
		err = json.NewDecoder(resp.Body).Decode(&newOrder)
		if err != nil {
			logger.Log.Info("error while decoding accrual response", zap.Error(err))
		}

		order.Status = newOrder.Status
		order.Accrual = &newOrder.Accrual
		resp.Body.Close()

		err = o.OrderRepository.UpdateOrderAndBalance(context.Background(), order)
		if err != nil {
			logger.Log.Info("error while updating order with new status", zap.Error(err))
		}
	}
}
