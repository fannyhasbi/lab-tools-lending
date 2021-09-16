package repository

import "github.com/fannyhasbi/lab-tools-lending/types"

type ToolQuery interface {
	FindByID(id int64) QueryResult
	Get() QueryResult
	GetAvailableTools() QueryResult
	GetPhotos(toolID int64) QueryResult
}

type ToolRepository interface {
	Save(tool *types.Tool) (int64, error)
	Update(tool *types.Tool) error
	SavePhotos(toolID int64, photos []types.TelePhotoSize) error
	DeletePhotos(toolID int64) error
	IncreaseStock(toolID int64, amount int) error
	DecreaseStock(toolID int64, amount int) error
}
