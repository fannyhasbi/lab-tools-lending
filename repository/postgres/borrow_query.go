package postgres

import (
	"database/sql"

	"github.com/fannyhasbi/lab-tools-lending/repository"
)

type BorrowQueryPostgres struct {
	DB *sql.DB
}

func NewBorrowQueryPostgres(DB *sql.DB) repository.BorrowQuery {
	return &BorrowQueryPostgres{
		DB: DB,
	}
}

func (bq BorrowQueryPostgres) FindByUserID(id int64) repository.QueryResult {
	return repository.QueryResult{}
}
