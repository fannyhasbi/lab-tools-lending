package postgres

import (
	"database/sql"

	"github.com/fannyhasbi/lab-tools-lending/repository"
	"github.com/fannyhasbi/lab-tools-lending/types"
)

type ToolRepositoryPostgres struct {
	DB *sql.DB
}

func NewToolRepositoryPostgres(DB *sql.DB) repository.ToolRepository {
	return &ToolRepositoryPostgres{
		DB: DB,
	}
}

func (tr *ToolRepositoryPostgres) Save(tool *types.Tool) error {
	return nil
}

func (tr *ToolRepositoryPostgres) Update(tool *types.Tool) error {
	return nil
}

func (tr *ToolRepositoryPostgres) IncreaseStock(toolID int64) error {
	_, err := tr.DB.Exec(`UPDATE tools SET stock = stock + 1 WHERE id = $1`, toolID)
	return err
}

func (tr *ToolRepositoryPostgres) DecreaseStock(toolID int64) error {
	_, err := tr.DB.Exec(`UPDATE tools SET stock = stock - 1 WHERE id = $1`, toolID)
	return err
}
