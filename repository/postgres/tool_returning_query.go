package postgres

import (
	"database/sql"

	"github.com/fannyhasbi/lab-tools-lending/repository"
)

type ToolReturningQueryPostgres struct {
	DB *sql.DB
}

func NewToolReturningQueryPostgres(DB *sql.DB) repository.ToolReturningQuery {
	return &ToolReturningQueryPostgres{
		DB: DB,
	}
}

func (trq ToolReturningQueryPostgres) FindByUserID(id int64) repository.QueryResult {
	return repository.QueryResult{}
}
