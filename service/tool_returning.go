package service

import (
	"github.com/fannyhasbi/lab-tools-lending/config"
	"github.com/fannyhasbi/lab-tools-lending/repository"
	"github.com/fannyhasbi/lab-tools-lending/repository/postgres"
	"github.com/fannyhasbi/lab-tools-lending/types"
)

type ToolReturningService struct {
	Query      repository.ToolReturningQuery
	Repository repository.ToolReturningRepository
}

func NewToolReturningService() *ToolReturningService {
	var toolReturningQuery repository.ToolReturningQuery
	var ToolReturningRepository repository.ToolReturningRepository

	db := config.InitPostgresDB()
	toolReturningQuery = postgres.NewToolReturningQueryPostgres(db)
	ToolReturningRepository = postgres.NewToolReturningRepositoryPostgres(db)

	return &ToolReturningService{
		Query:      toolReturningQuery,
		Repository: ToolReturningRepository,
	}
}

func (trs ToolReturningService) SaveToolReturning(toolReturning types.ToolReturning) (types.ToolReturning, error) {
	result, err := trs.Repository.Save(&toolReturning)
	if err != nil {
		return types.ToolReturning{}, err
	}

	return result, nil
}

func (trs ToolReturningService) UpdateToolReturningStatus(id int64, status types.ToolReturningStatus) error {
	return trs.Repository.UpdateStatus(id, status)
}

func (trs ToolReturningService) FindOnProgressByUserID(id int64) (types.ToolReturning, error) {
	result := trs.Query.FindByUserIDAndStatus(id, types.GetToolReturningStatus("progress"))
	if result.Error != nil {
		return types.ToolReturning{}, result.Error
	}

	return result.Result.(types.ToolReturning), nil
}
