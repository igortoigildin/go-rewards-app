package withdrawal

import (
	"context"
	"database/sql"
	"fmt"

	withdrawalEntity "github.com/igortoigildin/go-rewards-app/internal/entities/withdrawal"
)

type WithdrawalRepository struct {
	db *sql.DB
}

// NewUserRepository returns a new instance of the repository.
func NewWithdrawalRepository(db *sql.DB) *WithdrawalRepository {
	return &WithdrawalRepository{
		db: db,
	}
}

func (rep *WithdrawalRepository) Create(ctx context.Context, withdrawal *withdrawalEntity.Withdrawal) error {
	const op = "storage.postgres.withdrawal"

	tx, err := rep.db.Begin()
	if err != nil {
		return err
	}

	query := `
	INSERT INTO withdrawals (order_id, sum, user_id, date) VALUES ($1, $2, $3, now() AT TIME ZONE 'MSK')`
	args := []any{
		withdrawal.Order,
		withdrawal.Sum,
		withdrawal.UserID,
	}
	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("%s: insert withdrawal: unable to rollback: %v", op, rollbackErr)
		}
		return err
	}
	query = `UPDATE users SET balance = balance - $1 WHERE user_id = $2`
	args = []any{
		withdrawal.Sum,
		withdrawal.UserID,
	}
	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("%s: update user: unable to rollback: %v", op, rollbackErr)
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return tx.Commit()
}

func (rep *WithdrawalRepository) SelectAllForUserID(ctx context.Context, userID int64) ([]withdrawalEntity.Withdrawal, error) {
	var withdrawals []withdrawalEntity.Withdrawal
	query := `
	SELECT order_id, sum, date FROM withdrawals WHERE user_id = $1 ORDER BY date;`
	rows, err := rep.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var withdrawal withdrawalEntity.Withdrawal
		err = rows.Scan(&withdrawal.Order, &withdrawal.Sum, &withdrawal.Date)
		if err != nil {
			return nil, err
		}
		withdrawals = append(withdrawals, withdrawal)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return withdrawals, nil
}
