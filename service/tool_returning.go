package service

import (
	"fmt"
	"time"

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

func (trs ToolReturningService) UpdateToolReturningConfirm(id int64, datetime time.Time, firstName, lastName string) error {
	confirmedBy := firstName
	if len(lastName) > 0 {
		confirmedBy = fmt.Sprintf("%s %s", firstName, lastName)
	}

	return trs.Repository.UpdateConfirm(id, datetime, confirmedBy)
}

func (trs ToolReturningService) FindToolReturningByID(id int64) (types.ToolReturning, error) {
	result := trs.Query.FindByID(id)
	if result.Error != nil {
		return types.ToolReturning{}, result.Error
	}

	return result.Result.(types.ToolReturning), nil
}

func (trs ToolReturningService) FindOnProgressByUserID(id int64) (types.ToolReturning, error) {
	result := trs.Query.FindByUserIDAndStatus(id, types.GetToolReturningStatus("request"))
	if result.Error != nil {
		return types.ToolReturning{}, result.Error
	}

	return result.Result.(types.ToolReturning), nil
}

func (trs ToolReturningService) GetToolReturningRequests() ([]types.ToolReturning, error) {
	result := trs.Query.GetByStatus(types.GetToolReturningStatus("request"))
	if result.Error != nil {
		return []types.ToolReturning{}, result.Error
	}

	return result.Result.([]types.ToolReturning), nil
}
