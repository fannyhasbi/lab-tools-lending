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

func (br *BorrowRepositoryPostgres) UpdateStatus(id int64, status types.BorrowStatus) error {
	_, err := br.DB.Exec(`UPDATE borrows SET status = $1 WHERE id = $2`, status, id)
	return err
}

func (br *BorrowRepositoryPostgres) UpdateConfirm(id int64, confirmedAt time.Time, confirmedBy string) error {
	_, err := br.DB.Exec(`UPDATE borrows SET confirmed_at = $1, confirmed_by = $2 WHERE id = $3`, confirmedAt, confirmedBy, id)
	return err
}
