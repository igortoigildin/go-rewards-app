package order

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/igortoigildin/go-rewards-app/config"
	orderEntity "github.com/igortoigildin/go-rewards-app/internal/entities/order"
	"github.com/igortoigildin/go-rewards-app/internal/logger"
	"github.com/igortoigildin/go-rewards-app/internal/storage"
	"go.uber.org/zap"
)

const (
	statusNew        = "NEW"
	statusRegistered = "REGISTERED"
)

type OrderService struct {
	OrderRepository storage.OrderRepository
}

func (o *OrderService) UpdateAccruals(cfg *config.Config) {
	for {
		var wg sync.WaitGroup
		orders, _ := o.OrderRepository.SelectForCalc()
		numJobs := len(orders)
		jobs := make(chan int64, numJobs)
		results := make(chan orderEntity.Order, numJobs)

		for w := 1; w <= 3; w++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				worker(jobs, results, cfg)
			}()
		}

		for _, order := range orders {
			wg.Add(1)
			go func() {
				defer wg.Done()
				jobs <- order
			}()
		}
		close(jobs)

		for a := 1; a <= numJobs; a++ {
			response := <-results
			updateOrder(response) // TODO: add func for updating order with new status and accrual data in DB
		}
	}

}

func worker(jobs chan int64, results chan<- orderEntity.Order, cfg *config.Config) {
	for j := range jobs {
		url := cfg.FlagAccSysAddr + fmt.Sprintf("/api/orders/%v", j)

		resp, err := http.Get(url)
		if err != nil {
			logger.Log.Info("error while reaching accrual system", zap.Error(err))
		}

		var order orderEntity.Order
		switch resp.StatusCode {
		case http.StatusOK:
			err := json.NewDecoder(resp.Body).Decode(&order)
			if err != nil {
				logger.Log.Info("error while decoding accrual response", zap.Error(err))
			}

			if order.Status == statusRegistered { // if order status "REGISTERED" - send this number again
				jobs <- j
			} else {
				results <- order
			}
		case http.StatusNoContent:
			logger.Log.Info("order not registered")
		case http.StatusTooManyRequests:
			jobs <- j
			time.Sleep(cfg.PauseDuration * 2)
		default:
			jobs <- j
		}
		time.Sleep(cfg.PauseDuration)
	}
}

// func updateOrders() ([]int64, error){
// 	var orders []int64
// 	orders, err := o.OrderRepository.SelectAllNew()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return orders, nil
// }

// Returns -1 in case of success or returns user id who already added this order.
func (o *OrderService) InsertOrder(ctx context.Context, number string, userID int64) (int64, error) {
	t, err := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	if err != nil {
		logger.Log.Info("time parsing error", zap.Error(err))
		return 0, err
	}

	order := orderEntity.Order{
		Number:      number,
		Status:      newOrder,
		Uploaded_at: t,
		UserID:      userID,
	}
	id, err := o.OrderRepository.InsertOrder(ctx, &order)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// func (o *OrderService) SelectAllNew() ([]int64, error) {
// 	var orders []int64
// 	orders, err := o.OrderRepository.SelectAllNew()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return orders, nil
// }

func (o *OrderService) SelectAllByUser(ctx context.Context, userID int64) ([]orderEntity.Order, error) {
	orders, err := o.OrderRepository.SelectAllByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (o *OrderService) ValidateOrder(number string) (bool, error) {
	res, err := strconv.Atoi(number)
	if err != nil {
		logger.Log.Info("error while converting number", zap.Error(err))
		return false, err
	}
	return Valid(res), nil
}

// Valid check number is valid or not based on Luhn algorithm
func Valid(number int) bool {
	return (number%10+checksum(number/10))%10 == 0
}

func checksum(number int) int {
	var luhn int

	for i := 0; number > 0; i++ {
		cur := number % 10

		if i%2 == 0 { // even
			cur = cur * 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}

		luhn += cur
		number = number / 10
	}
	return luhn % 10
}

// NewOrderService returns a new instance of order service.
func NewOrderService(OrderRepository storage.OrderRepository) *OrderService {
	return &OrderService{
		OrderRepository: OrderRepository,
	}
}
