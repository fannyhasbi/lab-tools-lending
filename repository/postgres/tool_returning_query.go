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
		&ret.ReturnedAt,
		&ret.AdditionalInfo,
	)

	if err != nil {
		result.Error = err
		return result
	}

	result.Result = ret
	return result
}
