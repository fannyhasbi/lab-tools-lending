package repository

import "github.com/fannyhasbi/lab-tools-lending/types"

type ToolReturningQuery interface {
	FindByUserID(id int64) QueryResult
}

type ToolReturningRepository interface {
	Save(toolReturning *types.ToolReturning) (types.ToolReturning, error)
}
