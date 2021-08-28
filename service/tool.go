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
