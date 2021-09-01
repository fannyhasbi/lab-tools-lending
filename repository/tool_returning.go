package repository

import "github.com/fannyhasbi/lab-tools-lending/types"

type ToolReturningQuery interface {
	FindByID(id int64) QueryResult
	FindByUserIDAndStatus(id int64, status types.ToolReturningStatus) QueryResult
	GetByStatus(status types.ToolReturningStatus) QueryResult
}

type ToolReturningRepository interface {
	Save(toolReturning *types.ToolReturning) (types.ToolReturning, error)
	UpdateStatus(id int64, status types.ToolReturningStatus) error
}
