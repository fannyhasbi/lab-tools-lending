package repository

type ToolQuery interface {
	FindByID(id int64) QueryResult
	GetAvailableTools() QueryResult
}

type ToolRepository interface {
	IncreaseStock(toolID int64) error
	DecreaseStock(toolID int64) error
}
