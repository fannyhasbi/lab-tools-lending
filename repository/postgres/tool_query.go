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

func (tq ToolQueryPostgres) GetTool() repository.QueryResult {
	rows, err := tq.DB.Query(`SELECT * FROM tools`)

	articles := []types.Tool{}
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
				&temp.AdditionalInformation,
			)

			articles = append(articles, temp)
		}
		result.Result = articles
	}
	return result
}
