package service

import (
	"github.com/fannyhasbi/lab-tools-lending/config"
	"github.com/fannyhasbi/lab-tools-lending/repository"
	"github.com/fannyhasbi/lab-tools-lending/repository/postgres"
	"github.com/fannyhasbi/lab-tools-lending/types"
)

type ToolService struct {
	Query      repository.ToolQuery
	Repository repository.ToolRepository
}

func NewToolService() *ToolService {
	var toolQuery repository.ToolQuery
	var toolRepository repository.ToolRepository

	db := config.InitPostgresDB()
	toolQuery = postgres.NewToolQueryPostgres(db)
	toolRepository = postgres.NewToolRepositoryPostgres(db)

	return &ToolService{
		Query:      toolQuery,
		Repository: toolRepository,
	}
}

func (ts ToolService) SaveTool(tool types.Tool) (int64, error) {
	result, err := ts.Repository.Save(&tool)
	if err != nil {
		return int64(0), err
	}

	return result, nil
}

func (ts ToolService) SaveToolPhotos(toolID int64, photos []types.TelePhotoSize) error {
	return ts.Repository.SavePhotos(toolID, photos)
}

func (ts ToolService) IncreaseStock(id int64) error {
	return ts.Repository.IncreaseStock(id)
}

func (ts ToolService) DecreaseStock(id int64) error {
	return ts.Repository.DecreaseStock(id)
}

func (ts ToolService) FindByID(id int64) (types.Tool, error) {
	result := ts.Query.FindByID(id)

	if result.Error != nil {
		return types.Tool{}, result.Error
	}

	return result.Result.(types.Tool), nil
}

func (ts ToolService) GetAvailableTools() ([]types.Tool, error) {
	result := ts.Query.GetAvailableTools()

	if result.Error != nil {
		return []types.Tool{}, result.Error
	}

	return result.Result.([]types.Tool), nil
}
