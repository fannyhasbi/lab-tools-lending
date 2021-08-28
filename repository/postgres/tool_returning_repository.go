package postgres

import (
	"database/sql"

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
	return types.ToolReturning{}, nil
}
