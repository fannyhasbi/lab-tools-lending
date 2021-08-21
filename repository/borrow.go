package repository

import "github.com/fannyhasbi/lab-tools-lending/types"

type BorrowQuery interface {
	FindInitialByUserID(id int64) QueryResult
	FindByUserID(id int64) QueryResult
}

type BorrowRepository interface {
	Save(borrow *types.Borrow) (types.Borrow, error)
	Update(borrow *types.Borrow) (types.Borrow, error)
}
