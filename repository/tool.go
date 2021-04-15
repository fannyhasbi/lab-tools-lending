package repository

import "github.com/fannyhasbi/lab-tools-lending/types"

type ToolQuery interface {
	FindByID(id int64) QueryResult
	GetAvailableTools() QueryResult
}

type ToolRepository interface {
	Save(tool *types.Tool) error
	Update(tool *types.Tool) error
}
