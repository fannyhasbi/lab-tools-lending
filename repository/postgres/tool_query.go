package postgres

import (
	"database/sql"

	"github.com/fannyhasbi/lab-tools-lending/repository"
	"github.com/fannyhasbi/lab-tools-lending/types"
)

type ToolQueryPostgres struct {
	DB *sql.DB
}

func NewToolQueryPostgres(DB *sql.DB) repository.ToolQuery {
	return &ToolQueryPostgres{
		DB: DB,
	}
}

func (tq ToolQueryPostgres) GetAvailableTools() repository.QueryResult {
	rows, err := tq.DB.Query(`SELECT id, name, brand, product_type, weight, stock, additional_info, created_at, updated_at FROM tools WHERE stock > 0`)

	tools := []types.Tool{}
	result := repository.QueryResult{}

	if err != nil {
		result.Error = err
	} else {
		for rows.Next() {
			temp := types.Tool{}
			rows.Scan(
				&temp.ID,
				&temp.Name,
				&temp.Brand,
				&temp.ProductType,
				&temp.Weight,
				&temp.Stock,
				&temp.AdditionalInformation,
				&temp.CreatedAt,
				&temp.UpdatedAt,
			)

			tools = append(tools, temp)
		}
		result.Result = tools
	}
	return result
}
