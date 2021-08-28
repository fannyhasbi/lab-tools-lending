package postgres

import (
	"database/sql"

	"github.com/fannyhasbi/lab-tools-lending/repository"
	"github.com/fannyhasbi/lab-tools-lending/types"
)

type BorrowQueryPostgres struct {
	DB *sql.DB
}

func NewBorrowQueryPostgres(DB *sql.DB) repository.BorrowQuery {
	return &BorrowQueryPostgres{
		DB: DB,
	}
}

func (bq BorrowQueryPostgres) FindInitialByUserID(id int64) repository.QueryResult {
	row := bq.DB.QueryRow(`
		SELECT id, amount, return_date, status, user_id, tool_id, created_at
		FROM borrows
		WHERE user_id = $1
			AND status = $2
		ORDER BY id DESC
	`, id, types.GetBorrowStatus("init"))

	borrow := types.Borrow{}
	result := repository.QueryResult{}

	err := row.Scan(
		&borrow.ID,
		&borrow.Amount,
		&borrow.ReturnDate,
		&borrow.Status,
		&borrow.UserID,
		&borrow.ToolID,
		&borrow.CreatedAt,
	)

	if err != nil {
		result.Error = err
		return result
	}

	result.Result = borrow
	return result
}

func (bq BorrowQueryPostgres) FindByUserID(id int64) repository.QueryResult {
	rows, err := bq.DB.Query(`
		SELECT b.id, b.amount, b.return_date, b.status, b.user_id, b.tool_id, b.created_at, t.name AS tool_name
		FROM borrows b
		INNER JOIN tools t
			ON t.id = b.tool_id
		WHERE b.user_id = $1
		ORDER BY b.id DESC
	`, id)

	borrows := []types.Borrow{}
	result := repository.QueryResult{}

	if err != nil {
		result.Error = err
	} else {
		for rows.Next() {
			temp := types.Borrow{}
			rows.Scan(
				&temp.ID,
				&temp.Amount,
				&temp.ReturnDate,
				&temp.Status,
				&temp.UserID,
				&temp.ToolID,
				&temp.CreatedAt,
				&temp.Tool.Name,
			)

			borrows = append(borrows, temp)
		}
		result.Result = borrows
	}
	return result
}
