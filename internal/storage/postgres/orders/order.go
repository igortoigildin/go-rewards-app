package order

import (
	"context"
	"database/sql"

	orderEntity "github.com/igortoigildin/go-rewards-app/internal/entities/order"
)

const (
	statusNew        = "NEW"
	statusProcessing = "PROCESSING"
)

type OrderRepository struct {
	db *sql.DB
}

// NewOrderRepository returns a new instance of the repository.
func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (rep *OrderRepository) Update(ctx context.Context, order *orderEntity.Order) error {
	tx, err := rep.db.Begin()
	if err != nil {
		return err
	}
	query := `
	UPDATE orders SET status = $1, accrual = $2 WHERE number = $3`
	args := []any{
		order.Status,
		order.Accrual,
		order.Number,
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		tx.Rollback()
		return err
	}

	query = `UPDATE users SET balance = balance + $1 WHERE user_id = $2`
	args = []any{
		order.Accrual,
		order.UserID,
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

// SaveOrder saves new order in db or
// returns the user id who already saved this order.
// Returns -1 if added successfully.
func (rep *OrderRepository) InsertOrder(ctx context.Context, order *orderEntity.Order) (int64, error) {
	var userID int64
	query := `
	WITH new_orders AS (INSERT INTO orders (number, status, user_id, uploaded_at)
	VALUES ($1, $2, $3, now() AT TIME ZONE 'MSK') ON CONFLICT (number) DO NOTHING RETURNING user_id)
	SELECT COALESCE ((-1), (SELECT user_id FROM orders WHERE number = $1));
	`
	args := []any{
		order.Number,
		order.Status,
		order.UserID,
	}

	err := rep.db.QueryRowContext(ctx, query, args...).Scan(&userID)
	if err != nil {
		return 0, err

	}
	return userID, nil
}

func (rep *OrderRepository) SelectAllByUser(ctx context.Context, user int64) ([]orderEntity.Order, error) {
	var orders []orderEntity.Order
	query := `
	SELECT number, accrual, status, uploaded_at FROM orders WHERE user_id = $1 ORDER BY uploaded_at;`

	rows, err := rep.db.QueryContext(ctx, query, user)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var order orderEntity.Order
		err = rows.Scan(&order.Number, &order.Accrual, &order.Status, &order.Uploaded_at)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// Select numbers of all new orders.
func (rep *OrderRepository) SelectForAccrualCalc() ([]int64, error) {
	var orders []int64
	query := `
	SELECT number FROM orders WHERE status = $1 or status = $2`
	args := []any{statusNew, statusProcessing}
	rows, err := rep.db.Query(query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var order int64
		err = rows.Scan(&order)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return orders, nil
}
