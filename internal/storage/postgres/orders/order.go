package order

import (
	"context"
	"database/sql"

	orderEntity "github.com/igortoigildin/go-rewards-app/internal/entities/order"
)

type OrderRepository struct {
	DB *sql.DB
}

// NewOrderRepository returns a new instance of the repository.
func NewOrderRepository(DB *sql.DB) *OrderRepository {
	return &OrderRepository{
		DB: DB,
	}
}

// SaveOrder method saves new order in Orders table or
// returns the user_id who already inserted provided order
func (rep *OrderRepository) InsertOrder(ctx context.Context, order *orderEntity.Order) (int64, error) {
	var userID int64
	query := `
	WITH new_orders AS (INSERT INTO orders (number, status, user_id)
	VALUES ($1, $2, $3) ON CONFLICT (number) DO NOTHING RETURNING user_id)
	SELECT COALESCE ((NULL), (SELECT user_id FROM orders WHERE number = $1));
	`
	args := []any{
		order.Number,
		order.Status,
		order.UserID,
	}

	err := rep.DB.QueryRowContext(ctx, query, args...).Scan(&userID)
	if err != nil {
		return 0, err

	}
	return userID, nil
}
