package repository

import "github.com/fannyhasbi/lab-tools-lending/types"

type ToolQuery interface {
	FindByID(id int64) QueryResult
	GetAvailableTools() QueryResult
}

type ToolRepository interface {
	Save(tool *types.Tool) (int64, error)
	IncreaseStock(toolID int64) error
	DecreaseStock(toolID int64) error
}
