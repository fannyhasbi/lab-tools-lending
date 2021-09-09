package repository

import (
	"time"

	"github.com/fannyhasbi/lab-tools-lending/types"
)

type BorrowQuery interface {
	FindByID(id int64) QueryResult
	FindByUserIDAndStatus(id int64, status types.BorrowStatus) QueryResult
	FindByUserID(id int64) QueryResult
	GetByStatus(status types.BorrowStatus) QueryResult
	GetByUserIDAndMultipleStatus(id int64, statuses []types.BorrowStatus) QueryResult
	GetReport() QueryResult
}

type BorrowRepository interface {
	Save(borrow *types.Borrow) (int64, error)
	UpdateStatus(id int64, status types.BorrowStatus) error
	UpdateConfirm(id int64, confirmedAt time.Time, confirmedBy string) error
}
