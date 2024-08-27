package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/igortoigildin/go-rewards-app/config"
	"github.com/igortoigildin/go-rewards-app/internal/api"
	"github.com/igortoigildin/go-rewards-app/internal/entities/order"
	"github.com/igortoigildin/go-rewards-app/internal/logger"
	"github.com/igortoigildin/go-rewards-app/internal/service"
	"github.com/igortoigildin/go-rewards-app/internal/storage"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

func main() {
	cfg := config.LoadConfig()
	if err := logger.Initialize(cfg.FlagLogLevel); err != nil {
		log.Fatalf("can't initialize logger: %v", err)
	}

	db, err := sql.Open("pgx", cfg.FlagDBURI)
	if err != nil {
		logger.Log.Fatal("error while connecting to DB", zap.Error(err))
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Log.Fatal("error while closing db connection", zap.Error(err))
		}
	}()
	logger.Log.Info("database connection pool established")

	instance, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		logger.Log.Fatal("migration error", zap.Error(err))
	}

	migrator, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", instance)
	if err != nil {
		logger.Log.Fatal("migration error", zap.Error(err))
	}
	err = migrator.Up()
	if err != nil || errors.Is(err, migrate.ErrNoChange) {
		logger.Log.Fatal("migration error", zap.Error(err))
	}
	logger.Log.Info("database connection pool established")

	repository := storage.NewRepository(db)
	services := service.NewService(repository)

	go func(cfg *config.Config) {
		for {
			orders, err := services.OrderService.OrderRepository.SelectForAccrualCalc()
			if err != nil {
				logger.Log.Fatal("error while obtaining orders for accrual calcs", zap.Error(err))
			}
			var numJobs = len(orders)

			jobs := make(chan int64, numJobs)
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

				err = services.OrderService.OrderRepository.UpdateOrderAndBalance(context.Background(), &order)
				if err != nil {
					logger.Log.Info("error while updating order with new status", zap.Error(err))
					return
				}
			}
		}
	}(cfg)

	err = http.ListenAndServe(cfg.FlagRunAddr, api.Router(services, cfg))
	if err != nil {
		logger.Log.Fatal("database migrations applied", zap.Error(err))
	}
}

func worker(jobs <-chan int64, results chan<- order.Order, cfg *config.Config) {
	for i := range jobs {
		url := cfg.FlagAccSysAddr + fmt.Sprintf("/api/orders/%v", i)

		resp, err := http.Get(url)
		if err != nil {
			logger.Log.Info("error while reaching accrual system", zap.Error(err))
			return
		}

		var newOrder order.Order

		err = json.NewDecoder(resp.Body).Decode(&newOrder)
		if err != nil {
			if err == io.EOF {
				logger.Log.Info("EOF response, continue")
				resp.Body.Close()
				continue
			}
			logger.Log.Info("error while decoding accrual response", zap.Error(err))
			return
		}
		resp.Body.Close()

		results <- newOrder
	}
}
