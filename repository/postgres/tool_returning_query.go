package postgres

import (
	"database/sql"

	"github.com/fannyhasbi/lab-tools-lending/repository"
	"github.com/fannyhasbi/lab-tools-lending/types"
)

type ToolReturningQueryPostgres struct {
	DB *sql.DB
}

func NewToolReturningQueryPostgres(DB *sql.DB) repository.ToolReturningQuery {
	return &ToolReturningQueryPostgres{
		DB: DB,
	}
}

func (trq ToolReturningQueryPostgres) FindByID(id int64) repository.QueryResult {
	row := trq.DB.QueryRow(`
		SELECT tr.id, tr.borrow_id, tr.status, tr.created_at, tr.additional_info, b.amount, b.duration, b.tool_id, b.confirmed_at AS borrow_confirmed_at, t.name AS tool_name, b.user_id, u.name AS user_name, u.nim, u.address
		FROM tool_returning tr
		INNER JOIN borrows b
			ON b.id = tr.borrow_id
		INNER JOIN tools t
			ON t.id = b.tool_id
		INNER JOIN users u
			ON u.id = b.user_id
		WHERE tr.id = $1
	`, id)

	ret := types.ToolReturning{}
	result := repository.QueryResult{}

	err := row.Scan(
		&ret.ID,
		&ret.BorrowID,
		&ret.Status,
		&ret.CreatedAt,
		&ret.AdditionalInfo,
		&ret.Borrow.Amount,
		&ret.Borrow.Duration,
		&ret.Borrow.ToolID,
		&ret.Borrow.ConfirmedAt,
		&ret.Borrow.Tool.Name,
		&ret.Borrow.UserID,
		&ret.Borrow.User.Name,
		&ret.Borrow.User.NIM,
		&ret.Borrow.User.Address,
	)

	if err != nil {
		result.Error = err
		return result
	}

	result.Result = ret
	return result
}

func (trq ToolReturningQueryPostgres) GetByUserIDAndStatus(id int64, status types.ToolReturningStatus) repository.QueryResult {
	rows, err := trq.DB.Query(`
		SELECT tr.id, tr.borrow_id, tr.status, tr.created_at, tr.additional_info, b.user_id
		FROM tool_returning tr
		INNER JOIN borrows b
			ON b.id = tr.borrow_id
		WHERE b.user_id = $1 AND tr.status = $2
		ORDER BY tr.id ASC
	`, id, status)

	rets := []types.ToolReturning{}
	result := repository.QueryResult{}

	if err != nil {
		result.Error = err
	} else {
		for rows.Next() {
			temp := types.ToolReturning{}
			rows.Scan(
				&temp.ID,
				&temp.BorrowID,
				&temp.Status,
				&temp.CreatedAt,
				&temp.AdditionalInfo,
				&temp.Borrow.UserID,
			)

			rets = append(rets, temp)
		}
		result.Result = rets
	}

	result.Result = rets
	return result
}

func (trq ToolReturningQueryPostgres) GetByStatus(status types.ToolReturningStatus) repository.QueryResult {
	rows, err := trq.DB.Query(`
		SELECT tr.id, tr.borrow_id, tr.status, tr.created_at, tr.additional_info, t.name AS tool_name, u.name AS user_name
		FROM tool_returning tr
		INNER JOIN borrows b
			ON b.id = tr.borrow_id
		INNER JOIN tools t
			ON t.id = b.tool_id
		INNER JOIN users u
			ON u.id = b.user_id
		WHERE tr.status = $1
		ORDER BY tr.id ASC
	`, status)

	rets := []types.ToolReturning{}
	result := repository.QueryResult{}

	if err != nil {
		result.Error = err
	} else {
		for rows.Next() {
			temp := types.ToolReturning{}
			rows.Scan(
				&temp.ID,
				&temp.BorrowID,
				&temp.Status,
				&temp.CreatedAt,
				&temp.AdditionalInfo,
				&temp.Borrow.Tool.Name,
				&temp.Borrow.User.Name,
			)

			rets = append(rets, temp)
		}
		result.Result = rets
	}
	return result
}

func (trq ToolReturningQueryPostgres) GetReport(year, month int) repository.QueryResult {
	rows, err := trq.DB.Query(`SELECT tr.id, tr.borrow_id, tr.status, tr.created_at, tr.confirmed_at, tr.confirmed_by, b.amount, t.name AS tool_name, u.name AS user_name
		FROM tool_returning tr
		INNER JOIN borrows b
			ON b.id = tr.borrow_id
		INNER JOIN tools t
			ON t.id = b.tool_id
		INNER JOIN users u
			ON u.id = b.user_id
		WHERE tr.status = $1
			AND DATE_PART('year', tr.confirmed_at) = $2
			AND DATE_PART('month', tr.confirmed_at) = $3
		ORDER BY tr.id ASC
	`, types.GetToolReturningStatus("complete"), year, month)

	rets := []types.ToolReturning{}
	result := repository.QueryResult{}

	if err != nil {
		result.Error = err
	} else {
		for rows.Next() {
			temp := types.ToolReturning{}
			rows.Scan(
				&temp.ID,
				&temp.BorrowID,
				&temp.Status,
				&temp.CreatedAt,
				&temp.ConfirmedAt,
				&temp.ConfirmedBy,
				&temp.Borrow.Amount,
				&temp.Borrow.Tool.Name,
				&temp.Borrow.User.Name,
			)

			rets = append(rets, temp)
		}
		result.Result = rets
	}
	return result
}
