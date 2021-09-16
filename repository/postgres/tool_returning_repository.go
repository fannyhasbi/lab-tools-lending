package postgres

import (
	"database/sql"
	"time"

	"github.com/fannyhasbi/lab-tools-lending/repository"
	"github.com/fannyhasbi/lab-tools-lending/types"
)

type ToolReturningRepositoryPostgres struct {
	DB *sql.DB
}

func NewToolReturningRepositoryPostgres(DB *sql.DB) repository.ToolReturningRepository {
	return &ToolReturningRepositoryPostgres{
		DB: DB,
	}
}

func (trr *ToolReturningRepositoryPostgres) Save(toolReturning *types.ToolReturning) (types.ToolReturning, error) {
	stmt, err := trr.DB.Prepare(`INSERT INTO tool_returning (borrow_id, status, additional_info) VALUES ($1, $2, $3)
	RETURNING id, borrow_id, status, created_at, additional_info`)
	if err != nil {
		return types.ToolReturning{}, err
	}

	row := stmt.QueryRow(toolReturning.BorrowID, toolReturning.Status, toolReturning.AdditionalInfo)

	ret := types.ToolReturning{}
	err = row.Scan(
		&ret.ID,
		&ret.BorrowID,
		&ret.Status,
		&ret.CreatedAt,
		&ret.AdditionalInfo,
	)
	if err != nil {
		return types.ToolReturning{}, err
	}

	return ret, nil
}

func (trr *ToolReturningRepositoryPostgres) UpdateStatus(id int64, status types.ToolReturningStatus) error {
	_, err := trr.DB.Exec(`UPDATE tool_returning SET status = $1 WHERE id = $2`, status, id)
	return err
}

func (trr *ToolReturningRepositoryPostgres) UpdateConfirm(id int64, datetime time.Time, confirmedBy string) error {
	_, err := trr.DB.Exec(`UPDATE tool_returning SET confirmed_at = $1, confirmed_by = $2 WHERE id = $3`, datetime, confirmedBy, id)
	return err
}
