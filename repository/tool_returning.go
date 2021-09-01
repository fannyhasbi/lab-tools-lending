package repository

import "github.com/fannyhasbi/lab-tools-lending/types"

type ToolReturningQuery interface {
	FindByUserIDAndStatus(id int64, status types.ToolReturningStatus) QueryResult
}

type ToolReturningRepository interface {
	Save(toolReturning *types.ToolReturning) (types.ToolReturning, error)
	UpdateStatus(id int64, status types.ToolReturningStatus) error
}
