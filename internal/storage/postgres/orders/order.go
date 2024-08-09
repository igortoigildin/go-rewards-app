package orders

import (
	"context"
	"database/sql"
	"errors"

	"github.com/igortoigildin/go-rewards-app/internal/storage"
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
func (rep *OrderRepository) SaveOrder(ctx context.Context, orderNumber int) (int64, error) {
	var user int64
	firstQuery := `
	WITH res AS (INSERT INTO orders (number) VALUES ($1, $2, $3, $4) 
	ON CONFLICT (number) DO NOTHING RETURNING (user_id))
	SELECT * FROM res`
	secondQuery := `SELECT user_id FROM orders WHERE number = $5`
	

	err := rep.DB.QueryRowContext(ctx, firstQuery, orderNumber).Scan(&user)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return -1, storage.ErrRecordNotFound
		default:
			return user, err
		}
	}
	return user, nil
}