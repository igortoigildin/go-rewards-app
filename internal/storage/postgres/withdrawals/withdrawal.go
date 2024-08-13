package withdrawal

import (
	"context"
	"database/sql"

	withdrawalEntity "github.com/igortoigildin/go-rewards-app/internal/entities/withdrawal"
)

type WithdrawalRepository struct {
	DB *sql.DB
}

// NewUserRepository returns a new instance of the repository.
func NewWithdrawalRepository(DB *sql.DB) *WithdrawalRepository {
	return &WithdrawalRepository{
		DB: DB,
	}
}

func (rep *WithdrawalRepository) Create(ctx context.Context, withdrawal *withdrawalEntity.Withdrawal) error {
	query := `
	INSERT INTO withdrawals (order, sum, date, user_id) VALUES ($1, $2, $3, $4)`
	args := []any{
		withdrawal.Order,
		withdrawal.Sum,
		withdrawal.Date,
		withdrawal.UserID,
	}

	_, err := rep.DB.Exec(query, args)
	if err != nil {
		return err
	}
	return nil
}

func (rep *WithdrawalRepository) SelectAllForUserID(ctx context.Context, userID int64) ([]withdrawalEntity.Withdrawal, error) {
	var withdrawals []withdrawalEntity.Withdrawal

	query := `
	SELECT order, sum, date FROM withdrawals WHERE user_id = $1 ORDER BY date;`

	rows, err := rep.DB.QueryContext(ctx, query, userID)
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
