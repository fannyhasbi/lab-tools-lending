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

func (tq ToolQueryPostgres) FindByID(id int64) repository.QueryResult {
	row := tq.DB.QueryRow(`SELECT id, name, brand, product_type, weight, stock, additional_info, created_at, updated_at FROM tools WHERE id = $1`, id)

	tool := types.Tool{}
	result := repository.QueryResult{}

	err := row.Scan(
		&tool.ID,
		&tool.Name,
		&tool.Brand,
		&tool.ProductType,
		&tool.Weight,
		&tool.Stock,
		&tool.AdditionalInformation,
		&tool.CreatedAt,
		&tool.UpdatedAt,
	)

	if err != nil {
		result.Error = err
		return result
	}

	result.Result = tool
	return result
}

func (tq ToolQueryPostgres) GetAvailableTools() repository.QueryResult {
	rows, err := tq.DB.Query(`SELECT id, name, brand, product_type, weight, stock, additional_info, created_at, updated_at FROM tools WHERE stock > 0 ORDER BY id ASC`)

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

func (tq ToolQueryPostgres) GetPhotos(toolID int64) repository.QueryResult {
	rows, err := tq.DB.Query(`SELECT file_id, file_unique_id FROM tool_photos WHERE tool_id = $1 ORDER BY id DESC`, toolID)

	photos := []types.TelePhotoSize{}
	result := repository.QueryResult{}

	if err != nil {
		result.Error = err
	} else {
		for rows.Next() {
			temp := types.TelePhotoSize{}
			rows.Scan(
				&temp.FileID,
				&temp.FileUniqueID,
			)

			photos = append(photos, temp)
		}
		result.Result = photos
	}
	return result
}
