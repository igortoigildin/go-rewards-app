package order

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/igortoigildin/go-rewards-app/config"
	"github.com/igortoigildin/go-rewards-app/internal/entities/order"
	"github.com/igortoigildin/go-rewards-app/internal/logger"
	"go.uber.org/zap"
)

const statusInvalid string = "INVALID"
const statusProcessed string = "PROCESSED"

// Selects orders with status "NEW" or "PROCESSING" and sends it to accruals API.
func (o *OrderService) SendOrdersToAccrualAPI(ctx context.Context, cfg *config.Config) {
	for {
		ctx, cancel := context.WithTimeout(ctx, cfg.ContextTimout)
		defer cancel()

		orders, err := o.OrderRepository.SelectForAccrualCalc(ctx)
		if err != nil {
			// check if no new orders with status "INVALID" or "PROCESSING" available
			if errors.Is(err, sql.ErrNoRows) {
				time.Sleep(1 * time.Second)
				continue
			}
			logger.Log.Error("error while obtaining orders for accrual calcs", zap.Error(err))
			return
		}
		var numJobs = len(orders)
		var wg sync.WaitGroup
		newOrdersChan := make(chan order.Order, numJobs) // chan for jobs
		results := make(chan order.Order, numJobs)       // chan for results

		// Worker pool
		for w := 1; w <= cfg.FlagRateLimit; w++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				processOrders(newOrdersChan, results, cfg)
			}()
		}

		// Send selected new orders to jobs (newOrdersChan) chan
		for _, v := range orders {
			newOrdersChan <- v
		}
		close(newOrdersChan)

		// Read from results chan
		for a := 1; a <= numJobs; a++ {
			order := <-results
			// If recievied order status not expected - resend order to jobs chan
			if order.Status != statusInvalid && order.Status != statusProcessed {
				newOrdersChan <- order
				continue
			}
			err = o.OrderRepository.UpdateOrderAndBalance(ctx, &order)
			if err != nil {
				logger.Log.Info("error while updating order with new status", zap.Error(err))
				return
			}
		}

		wg.Wait()
	}
}

// Worker for sending orders from newOrdersChan to Accrual API and sending results to results chan.
func processOrders(newOrdersChan chan order.Order, results chan<- order.Order, cfg *config.Config) {
	for i := range newOrdersChan {
		logger.Log.Info("sending new order number to accrual system for processing: %s", zap.String("number", i.Number))

		url := cfg.FlagAccSysAddr + fmt.Sprintf("/api/orders/%v", i.Number)
		resp, err := http.Get(url)
		if err != nil {
			logger.Log.Info("error while reaching accrual system", zap.Error(err))
			return
		}
		switch resp.StatusCode {
		case http.StatusNoContent:
			logger.Log.Info("order number not registered in accrual system")
			continue
		case http.StatusTooManyRequests:
			wait(resp)
		}

		updOrder := order.Order{
			Number: i.Number,
			UserID: i.UserID,
		}
		err = json.NewDecoder(resp.Body).Decode(&updOrder)
		if err != nil {
			newOrdersChan <- i // resend order request in case of unexpected response
			resp.Body.Close()
			continue
		}
		resp.Body.Close()
		results <- updOrder
	}
}

// Get "Retry-After" header from response and pause accoringly before next attempt.
// If not such header received, wait 30 sec by default.
func wait(resp *http.Response) {
	logger.Log.Info("accrual system response: too many requests")
	headers := resp.Header["Retry-After"]
	waitTime, err := strconv.Atoi(headers[0])
	if err != nil {
		logger.Log.Error("error while converting time pause", zap.Error(err))
	}
	if waitTime == 0 {
		waitTime = 30
	}
	time.Sleep(time.Duration(waitTime) * time.Second)
}
