package repository

import "github.com/fannyhasbi/lab-tools-lending/types"

type ToolQuery interface {
	GetTool() QueryResult
}

type ToolRepository interface {
	Save(tool *types.Tool) error
	Update(tool *types.Tool) error
}
