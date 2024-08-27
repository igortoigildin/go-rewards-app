package order

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/igortoigildin/go-rewards-app/internal/entities/order"
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

func (rep *OrderRepository) UpdateOrderAndBalance(ctx context.Context, order *orderEntity.Order) error {
	const op = "postgres.orders.order.UpdateOrderAndBalance"

	tx, err := rep.db.Begin()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
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
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("%s: update order: unable to rollback: %v", op, rollbackErr)
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	query = `UPDATE users SET balance = balance + $1 WHERE user_id = $2`
	args = []any{
		order.Accrual,
		order.UserID,
	}
	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("%s: update order: unable to rollback: %v", op, rollbackErr)
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return tx.Commit()
}

func (rep *OrderRepository) InsertOrder(ctx context.Context, order *orderEntity.Order) (int64, error) {
	var userID int64
	err := rep.db.QueryRowContext(ctx, `SELECT user_id FROM orders WHERE number = $1;`, order.Number).Scan(&userID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			// no such order exists
		default:
			return 0, err
		}
	}
	if userID != 0 { // order already exists, return
		return userID, nil
	}
	query := "INSERT INTO orders (number, status, user_id, uploaded_at)" +
		"VALUES ($1, $2, $3, now() AT TIME ZONE 'MSK')"
	args := []any{
		order.Number,
		order.Status,
		order.UserID,
	}
	_, err = rep.db.ExecContext(ctx, query, args...) // insert order accordingly
	if err != nil {
		return 0, err
	}
	return 0, nil
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
		err = rows.Scan(&order.Number, &order.Accrual, &order.Status, &order.UploadedAt)
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
func (rep *OrderRepository) SelectForAccrualCalc() ([]order.Order, error) {
	var orders []order.Order
	query := `
	SELECT * FROM orders WHERE status = $1 or status = $2`
	args := []any{statusNew, statusProcessing}
	rows, err := rep.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var order order.Order
		err = rows.Scan(&order.Number, &order.Status, &order.Accrual, &order.UploadedAt, &order.UserID)
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
