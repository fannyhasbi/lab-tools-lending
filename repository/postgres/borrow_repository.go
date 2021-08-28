package postgres

import (
	"database/sql"

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

func (br *BorrowRepositoryPostgres) Save(borrow *types.Borrow) (types.Borrow, error) {
	row := br.DB.QueryRow(`INSERT INTO borrows (amount, status, user_id, tool_id) VALUES ($1, $2, $3, $4)
		RETURNING id, amount, return_date, status, user_id, tool_id, created_at`, borrow.Amount, borrow.Status, borrow.UserID, borrow.ToolID)

	b := types.Borrow{}
	err := row.Scan(
		&b.ID,
		&b.Amount,
		&b.ReturnDate,
		&b.Status,
		&b.UserID,
		&b.ToolID,
		&b.CreatedAt,
	)
	if err != nil {
		return types.Borrow{}, err
	}

	return b, nil
}

func (br *BorrowRepositoryPostgres) Update(borrow *types.Borrow) (types.Borrow, error) {
	row := br.DB.QueryRow(`UPDATE borrows SET
		amount = $1,
		return_date = $2,
		status = $3,
		user_id = $4,
		tool_id = $5
		WHERE id = $6
		RETURNING id, amount, return_date, status, user_id, tool_id, created_at
	`, borrow.Amount, borrow.ReturnDate, borrow.Status, borrow.UserID, borrow.ToolID, borrow.ID)

	b := types.Borrow{}
	err := row.Scan(
		&b.ID,
		&b.Amount,
		&b.ReturnDate,
		&b.Status,
		&b.UserID,
		&b.ToolID,
		&b.CreatedAt,
	)
	if err != nil {
		return types.Borrow{}, err
	}

	return b, nil
}
