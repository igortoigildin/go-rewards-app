package api

import (
	"context"

	"github.com/igortoigildin/go-rewards-app/config"
	"github.com/igortoigildin/go-rewards-app/internal/service"
)

func RunAccrualUpdates(ctx context.Context, cfg *config.Config, services *service.Service) {
	services.OrderService.UpdateAccruals(ctx, cfg)
}
