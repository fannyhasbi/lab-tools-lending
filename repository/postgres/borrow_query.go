package postgres

import (
	"database/sql"
	"fmt"

	"github.com/fannyhasbi/lab-tools-lending/repository"
	"github.com/fannyhasbi/lab-tools-lending/types"
	"github.com/lib/pq"
)

type BorrowQueryPostgres struct {
	DB *sql.DB
}

func NewBorrowQueryPostgres(DB *sql.DB) repository.BorrowQuery {
	return &BorrowQueryPostgres{
		DB: DB,
	}
}

func (bq BorrowQueryPostgres) FindByID(id int64) repository.QueryResult {
	row := bq.DB.QueryRow(`
	SELECT b.id, b.amount, b.duration, b.status, b.user_id, b.tool_id, b.created_at, b.confirmed_at, b.reason, t.name AS tool_name, t.stock AS tool_stock, u.name AS user_name, u.nim
	FROM borrows b
	INNER JOIN tools t
		ON t.id = b.tool_id
	INNER JOIN users u
		ON u.id = b.user_id
	WHERE b.id = $1
	`, id)

	borrow := types.Borrow{}
	result := repository.QueryResult{}

	err := row.Scan(
		&borrow.ID,
		&borrow.Amount,
		&borrow.Duration,
		&borrow.Status,
		&borrow.UserID,
		&borrow.ToolID,
		&borrow.CreatedAt,
		&borrow.ConfirmedAt,
		&borrow.Reason,
		&borrow.Tool.Name,
		&borrow.Tool.Stock,
		&borrow.User.Name,
		&borrow.User.NIM,
	)

	if err != nil {
		result.Error = err
		return result
	}

	result.Result = borrow
	return result
}

func (bq BorrowQueryPostgres) FindByUserIDAndStatus(id int64, status types.BorrowStatus) repository.QueryResult {
	row := bq.DB.QueryRow(`
		SELECT b.id, b.amount, b.duration, b.status, b.user_id, b.tool_id, b.created_at, b.confirmed_at, t.name AS tool_name, u.name AS user_name
		FROM borrows b
		INNER JOIN tools t
			ON t.id = b.tool_id
		INNER JOIN users u
			ON u.id = b.user_id
		WHERE b.user_id = $1
			AND b.status = $2
		ORDER BY b.id DESC
	`, id, status)

	borrow := types.Borrow{}
	result := repository.QueryResult{}

	err := row.Scan(
		&borrow.ID,
		&borrow.Amount,
		&borrow.Duration,
		&borrow.Status,
		&borrow.UserID,
		&borrow.ToolID,
		&borrow.CreatedAt,
		&borrow.ConfirmedAt,
		&borrow.Tool.Name,
		&borrow.User.Name,
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
		SELECT b.id, b.amount, b.duration, b.status, b.user_id, b.tool_id, b.created_at, b.confirmed_at, t.name AS tool_name
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
				&temp.Duration,
				&temp.Status,
				&temp.UserID,
				&temp.ToolID,
				&temp.CreatedAt,
				&temp.ConfirmedAt,
				&temp.Tool.Name,
			)

			borrows = append(borrows, temp)
		}
		result.Result = borrows
	}
	return result
}

func (bq BorrowQueryPostgres) GetByStatus(status types.BorrowStatus) repository.QueryResult {
	rows, err := bq.DB.Query(`
		SELECT b.id, b.amount, b.duration, b.status, b.user_id, b.tool_id, b.created_at, b.confirmed_at, t.name AS tool_name, u.name AS user_name
		FROM borrows b
		INNER JOIN tools t
			ON t.id = b.tool_id
		INNER JOIN users u
			ON u.id = b.user_id
		WHERE b.status = $1
		ORDER BY b.id ASC
	`, status)

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
				&temp.Duration,
				&temp.Status,
				&temp.UserID,
				&temp.ToolID,
				&temp.CreatedAt,
				&temp.ConfirmedAt,
				&temp.Tool.Name,
				&temp.User.Name,
			)

			borrows = append(borrows, temp)
		}
		result.Result = borrows
	}
	return result
}

func (bq BorrowQueryPostgres) GetByUserIDAndMultipleStatus(id int64, statuses []types.BorrowStatus) repository.QueryResult {
	rows, err := bq.DB.Query(`
		SELECT b.id, b.amount, b.duration, b.status, b.user_id, b.tool_id, b.created_at, b.confirmed_at, t.name AS tool_name, u.name AS user_name
		FROM borrows b
		INNER JOIN tools t
			ON t.id = b.tool_id
		INNER JOIN users u
			ON u.id = b.user_id
		WHERE b.user_id = $1 AND b.status = ANY($2)
		ORDER BY b.id ASC
	`, id, pq.Array(statuses))

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
				&temp.Duration,
				&temp.Status,
				&temp.UserID,
				&temp.ToolID,
				&temp.CreatedAt,
				&temp.ConfirmedAt,
				&temp.Tool.Name,
				&temp.User.Name,
			)

			borrows = append(borrows, temp)
		}
		result.Result = borrows
	}
	return result
}

func (bq BorrowQueryPostgres) GetReport() repository.QueryResult {
	rows, err := bq.DB.Query(fmt.Sprintf(`
		SELECT b.id, b.amount, b.duration, b.status, b.user_id, b.tool_id, b.created_at, b.confirmed_at, b.confirmed_by, t.name AS tool_name, u.name AS user_name
		FROM borrows b
		INNER JOIN tools t
			ON t.id = b.tool_id
		INNER JOIN users u
			ON u.id = b.user_id
		WHERE b.status IN ('%s', '%s')
		ORDER BY b.id ASC
	`, types.GetBorrowStatus("progress"), types.GetBorrowStatus("returned")))

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
				&temp.Duration,
				&temp.Status,
				&temp.UserID,
				&temp.ToolID,
				&temp.CreatedAt,
				&temp.ConfirmedAt,
				&temp.ConfirmedBy,
				&temp.Tool.Name,
				&temp.User.Name,
			)

			borrows = append(borrows, temp)
		}
		result.Result = borrows
	}
	return result
}
