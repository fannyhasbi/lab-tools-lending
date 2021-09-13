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
		SELECT tr.id, tr.user_id, tr.tool_id, tr.status, tr.created_at, tr.additional_info, t.name AS tool_name, u.name AS user_name, u.nim
		FROM tool_returning tr
		INNER JOIN tools t
			ON t.id = tr.tool_id
		INNER JOIN users u
			ON u.id = tr.user_id
		WHERE tr.id = $1
	`, id)

	ret := types.ToolReturning{}
	result := repository.QueryResult{}

	err := row.Scan(
		&ret.ID,
		&ret.UserID,
		&ret.ToolID,
		&ret.Status,
		&ret.CreatedAt,
		&ret.AdditionalInfo,
		&ret.Tool.Name,
		&ret.User.Name,
		&ret.User.NIM,
	)

	if err != nil {
		result.Error = err
		return result
	}

	result.Result = ret
	return result
}

func (trq ToolReturningQueryPostgres) FindByUserIDAndStatus(id int64, status types.ToolReturningStatus) repository.QueryResult {
	row := trq.DB.QueryRow(`
		SELECT id, user_id, tool_id, status, created_at, additional_info
		FROM tool_returning
		WHERE user_id = $1 AND status = $2
		ORDER BY id DESC
	`, id, status)

	ret := types.ToolReturning{}
	result := repository.QueryResult{}

	err := row.Scan(
		&ret.ID,
		&ret.UserID,
		&ret.ToolID,
		&ret.Status,
		&ret.CreatedAt,
		&ret.AdditionalInfo,
	)

	if err != nil {
		result.Error = err
		return result
	}

	result.Result = ret
	return result
}

func (trq ToolReturningQueryPostgres) GetByStatus(status types.ToolReturningStatus) repository.QueryResult {
	rows, err := trq.DB.Query(`
		SELECT tr.id, tr.user_id, tr.tool_id, tr.status, tr.created_at, tr.additional_info, t.name AS tool_name, u.name AS user_name
		FROM tool_returning tr
		INNER JOIN tools t
			ON t.id = tr.tool_id
		INNER JOIN users u
			ON u.id = tr.user_id
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
				&temp.UserID,
				&temp.ToolID,
				&temp.Status,
				&temp.CreatedAt,
				&temp.AdditionalInfo,
				&temp.Tool.Name,
				&temp.User.Name,
			)

			rets = append(rets, temp)
		}
		result.Result = rets
	}
	return result
}

func (trq ToolReturningQueryPostgres) GetReport(year, month int) repository.QueryResult {
	rows, err := trq.DB.Query(`SELECT tr.id, tr.user_id, tr.tool_id, tr.status, tr.created_at, tr.confirmed_at, tr.confirmed_by, t.name AS tool_name, u.name AS user_name
		FROM tool_returning tr
		INNER JOIN tools t
			ON t.id = tr.tool_id
		INNER JOIN users u
			ON u.id = tr.user_id
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
				&temp.UserID,
				&temp.ToolID,
				&temp.Status,
				&temp.CreatedAt,
				&temp.ConfirmedAt,
				&temp.ConfirmedBy,
				&temp.Tool.Name,
				&temp.User.Name,
			)

			rets = append(rets, temp)
		}
		result.Result = rets
	}
	return result
}
