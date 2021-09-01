package repository

import "github.com/fannyhasbi/lab-tools-lending/types"

type BorrowQuery interface {
	FindByUserIDAndStatus(id int64, status types.BorrowStatus) QueryResult
	FindByUserID(id int64) QueryResult
	GetByStatus(status types.BorrowStatus) QueryResult
}

type BorrowRepository interface {
	Save(borrow *types.Borrow) (types.Borrow, error)
	Update(borrow *types.Borrow) (types.Borrow, error)
}
