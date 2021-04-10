package postgres

import (
	"database/sql"

	"github.com/fannyhasbi/lab-tools-lending/repository"
	"github.com/fannyhasbi/lab-tools-lending/types"
)

type ToolsQueryPostgres struct {
	DB *sql.DB
}

func NewToolsQueryPostgres(DB *sql.DB) repository.ToolsQuery {
	return &ToolsQueryPostgres{
		DB: DB,
	}
}

func (tq ToolsQueryPostgres) GetTools() repository.QueryResult {
	rows, err := tq.DB.Query(`SELECT * FROM tools`)

	articles := []types.Tools{}
	result := repository.QueryResult{}

	if err != nil {
		result.Error = err
	} else {
		for rows.Next() {
			temp := types.Tools{}
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
