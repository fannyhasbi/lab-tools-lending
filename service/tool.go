package service

import (
	"github.com/fannyhasbi/lab-tools-lending/config"
	"github.com/fannyhasbi/lab-tools-lending/repository"
	"github.com/fannyhasbi/lab-tools-lending/repository/postgres"
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
