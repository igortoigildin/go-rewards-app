package order

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/igortoigildin/go-rewards-app/config"
	orderEntity "github.com/igortoigildin/go-rewards-app/internal/entities/order"
	"github.com/igortoigildin/go-rewards-app/internal/logger"
	"go.uber.org/zap"
)

var statusInvalid string = "INVALID"
var statusProcessed string = "PROCESSED"

// func (o *OrderService) UpdateAccruals(ctx context.Context, cfg *config.Config) {

// 		var wg sync.WaitGroup
// 		orders, err := o.OrderRepository.SelectForAccrualCalc()
// 		if err != nil {
// 			switch {
// 			case errors.Is(err, sql.ErrNoRows):
// 				logger.Log.Info("no new orders for accrual calculation found")
// 			default:
// 				logger.Log.Info("error while selecting orders for accrual recalulation", zap.Error(err))
// 			}
// 		}

// 		fmt.Println("ORDERS", orders)

// 		jobs := make(chan int64, 10) // chan with order numbers for accrual calculation
// 		results := make(chan orderEntity.Order, 10)

// 		for w := 1; w <= cfg.FlagRateLimit; w++ {
// 			wg.Add(1)
// 			go func() {
// 				defer wg.Done()
// 				worker(jobs, results, cfg)
// 			}()
// 		}

// 		for _, order := range orders {
// 			wg.Add(1)
// 			go func() {
// 				defer wg.Done()
// 				jobs <- order
// 			}()
// 		}

// 		for a := 1; a <= len(orders); a++ {
// 			order := <-results
// 			err := o.OrderRepository.Update(ctx, &order) // updating order in DB with calculated accrual accordingly
// 			if err != nil {
// 				logger.Log.Info("error while updating order", zap.Error(err))
// 			}
// 		}

// 		wg.Wait()

// 		close(jobs)

// }

// // Sends recived orders to accrual system and work with responses.
// func worker(jobs chan int64, results chan<- orderEntity.Order, cfg *config.Config) {
// 	for j := range jobs {
// 		url := cfg.FlagAccSysAddr + fmt.Sprintf("/api/orders/%v", j)

// 		fmt.Println("worker startred")

// 		resp, err := http.Get(url)
// 		if err != nil {
// 			logger.Log.Info("error while reaching accrual system", zap.Error(err))
// 		}
// 		resp.Body.Close()

// 		switch resp.StatusCode {
// 		case http.StatusOK:
// 			processOrderStatusOK(resp, jobs, results, j)
// 		case http.StatusNoContent:
// 			logger.Log.Info("order not registered")
// 			//////
// 		case http.StatusTooManyRequests:
// 			jobs <- j
// 			time.Sleep(cfg.PauseDuration * 2)
// 		default:
// 			jobs <- j
// 		}
// 		time.Sleep(cfg.PauseDuration)
// 	}
// }

// func processOrderStatusOK(resp *http.Response, jobs chan int64, results chan<- orderEntity.Order, j int64) {
// 	var order orderEntity.Order
// 	err := json.NewDecoder(resp.Body).Decode(&order)
// 	if err != nil {
// 		logger.Log.Info("error while decoding accrual response", zap.Error(err))
// 	}

// 	switch {
// 	case order.Status == statusRegistered || order.Status == statusProcessing:
// 		jobs <- j // in case "REGISTERED" - send this number again
// 	default:
// 		results <- order
// 	}
// }

func (o *OrderService) UpdateAccruals(ctx context.Context, cfg *config.Config) {

	orders, err := o.OrderRepository.SelectForAccrualCalc()
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			logger.Log.Info("no new orders for accrual calculation found")
			return
		default:
			logger.Log.Info("error while selecting orders for accrual recalulation", zap.Error(err))
			return
		}
	}

	for _, val := range orders {
		var status string
		var order orderEntity.Order
		for status != statusInvalid || status != statusProcessed {

			url := cfg.FlagAccSysAddr + fmt.Sprintf("/api/orders/%v", val)

			fmt.Println("worker startred")

			resp, err := http.Get(url)
			if err != nil {
				logger.Log.Info("error while reaching accrual system", zap.Error(err))
			}
			resp.Body.Close()

			err = json.NewDecoder(resp.Body).Decode(&order)
			if err != nil {
				logger.Log.Info("error while decoding accrual response", zap.Error(err))
			}

			status = order.Status
		}

		fmt.Println("STATUS FINAL - ", status)
		

		err := o.OrderRepository.Update(ctx, &order) // updating order in DB with calculated accrual accordingly
		if err != nil {
			logger.Log.Info("error while updating order", zap.Error(err))
		}

		fmt.Println("ORDER FINAL - ", order)
	}
}
