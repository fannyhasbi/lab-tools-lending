package postgres

import (
	"database/sql"
	"time"

	"github.com/fannyhasbi/lab-tools-lending/repository"
	"github.com/fannyhasbi/lab-tools-lending/types"
)

type BorrowRepositoryPostgres struct {
	DB *sql.DB
}

func NewBorrowRepositoryPostgres(DB *sql.DB) repository.BorrowRepository {
	return &BorrowRepositoryPostgres{
		DB: DB,
	}
}

func (br *BorrowRepositoryPostgres) Save(borrow *types.Borrow) (int64, error) {
	row := br.DB.QueryRow(`INSERT INTO borrows (amount, status, user_id, tool_id, reason, duration)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`, borrow.Amount, borrow.Status, borrow.UserID, borrow.ToolID, borrow.Reason, borrow.Duration)

	var id int64
	err := row.Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}

func (br *BorrowRepositoryPostgres) Update(borrow *types.Borrow) (types.Borrow, error) {
	row := br.DB.QueryRow(`UPDATE borrows SET
		amount = $1,
		duration = $2,
		status = $3,
		user_id = $4,
		tool_id = $5
		WHERE id = $6
		RETURNING id, amount, duration, status, user_id, tool_id, created_at, confirmed_at
	`, borrow.Amount, borrow.Duration, borrow.Status, borrow.UserID, borrow.ToolID, borrow.ID)

	b := types.Borrow{}
	err := row.Scan(
		&b.ID,
		&b.Amount,
		&b.Duration,
		&b.Status,
		&b.UserID,
		&b.ToolID,
		&b.CreatedAt,
		&b.ConfirmedAt,
	)
	if err != nil {
		return types.Borrow{}, err
	}

	return b, nil
}

func (br *BorrowRepositoryPostgres) UpdateConfirmedAt(id int64, confirmedAt time.Time) error {
	_, err := br.DB.Exec(`UPDATE borrows SET confirmed_at = $1 WHERE id = $2`, confirmedAt, id)
	return err
}
