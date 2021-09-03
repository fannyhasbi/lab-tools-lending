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

func (tr *ToolRepositoryPostgres) Save(tool *types.Tool) (int64, error) {
	stmt, err := tr.DB.Prepare(`INSERT INTO tools (name, brand, product_type, weight, stock, additional_info)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`)

	if err != nil {
		return int64(0), err
	}

	row := stmt.QueryRow(tool.Name, tool.Brand, tool.ProductType, tool.Weight, tool.Stock, tool.AdditionalInformation)

	var id int64
	err = row.Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
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
