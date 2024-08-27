package order

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/igortoigildin/go-rewards-app/config"
	"github.com/igortoigildin/go-rewards-app/internal/entities/order"
	"github.com/igortoigildin/go-rewards-app/internal/logger"
	"go.uber.org/zap"
)

const statusInvalid string = "INVALID"
const statusProcessed string = "PROCESSED"

func (o *OrderService) SendOrdersToAccrualAPI(cfg *config.Config) {
	for {
		orders, err := o.OrderRepository.SelectForAccrualCalc()
		if err != nil {
			logger.Log.Fatal("error while obtaining orders for accrual calcs", zap.Error(err))
		}
		var numJobs = len(orders)

		if numJobs == 0 {
			continue
		}

		jobs := make(chan order.Order, numJobs)
		results := make(chan order.Order, numJobs)

		for w := 1; w <= 3; w++ {
			go worker(jobs, results, cfg)
		}

		for _, v := range orders {
			jobs <- v
		}
		close(jobs)

		for a := 1; a <= numJobs; a++ {
			order := <-results

			if order.Status != statusInvalid && order.Status != statusProcessed {
				jobs <- order
				continue
			}

			err = o.OrderRepository.UpdateOrderAndBalance(context.Background(), &order)
			if err != nil {
				logger.Log.Info("error while updating order with new status", zap.Error(err))
				return
			}
		}
	}
}

func worker(jobs <-chan order.Order, results chan<- order.Order, cfg *config.Config) {
	for i := range jobs {
		logger.Log.Info("sending new order number to accrual system for processing: %s", zap.String("number", i.Number))

		url := cfg.FlagAccSysAddr + fmt.Sprintf("/api/orders/%v", i.Number)
		resp, err := http.Get(url)
		if err != nil {
			logger.Log.Info("error while reaching accrual system", zap.Error(err))
			return
		}

		switch resp.StatusCode {
		case 204:
			logger.Log.Info("order number not registered in accrual system")
			continue
		case 429:
			logger.Log.Info("accrual system response: too many requests")
			headers := resp.Header["Retry-After"]
			pause, err := strconv.Atoi(headers[0])
			if err != nil {
				logger.Log.Error("error while converting time pause", zap.Error(err))
				continue
			}
			time.Sleep(time.Duration(pause) * time.Second)
		}


		var updOrder order.Order
		updOrder.Number =  i.Number
		updOrder.UserID = i.UserID
		err = json.NewDecoder(resp.Body).Decode(&updOrder)
		if err != nil {
			switch {
			case err == io.EOF:
				resp.Body.Close()
				continue
			default:
				logger.Log.Info("error while decoding accrual response", zap.Error(err))
				return
			}	
		}
		resp.Body.Close()

		results <- updOrder
	}
}