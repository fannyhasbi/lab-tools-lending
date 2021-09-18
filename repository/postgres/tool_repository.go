package postgres

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

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
	stmt, err := tr.DB.Prepare(`UPDATE tools SET name = $1, brand = $2, product_type = $3, weight = $4, stock = $5, additional_info = $6
		WHERE id = $7`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(tool.Name, tool.Brand, tool.ProductType, tool.Weight, tool.Stock, tool.AdditionalInformation, tool.ID)
	return err
}

func (tr *ToolRepositoryPostgres) Delete(toolID int64, deletedAt time.Time) error {
	stmt, err := tr.DB.Prepare(`UPDATE tools SET deleted_at = $1 WHERE id = $2`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(deletedAt, toolID)
	return err
}

func (tr *ToolRepositoryPostgres) SavePhotos(toolID int64, photos []types.TelePhotoSize) error {
	columns := []string{"tool_id", "file_id", "file_unique_id"}

	columnStr := ""
	for i := range columns {
		columnStr += columns[i] + ","
	}
	columnStr = columnStr[:len(columnStr)-1]

	query := fmt.Sprintf("INSERT INTO tool_photos (%s) VALUES ", columnStr)

	values := []interface{}{}
	for i, s := range photos {
		values = append(values, toolID, s.FileID, s.FileUniqueID)

		numFields := len(columns)
		n := i * numFields

		query += `(`
		for j := 0; j < numFields; j++ {
			query += `$` + strconv.Itoa(n+j+1) + `,`
		}
		query = query[:len(query)-1] + `),`
	}
	query = query[:len(query)-1] // remove the trailing comma

	_, err := tr.DB.Exec(query, values...)
	return err
}

func (tr *ToolRepositoryPostgres) DeletePhotos(toolID int64) error {
	_, err := tr.DB.Exec(`DELETE FROM tool_photos WHERE tool_id = $1`, toolID)
	return err
}

func (tr *ToolRepositoryPostgres) IncreaseStock(toolID int64, amount int) error {
	_, err := tr.DB.Exec(`UPDATE tools SET stock = stock + $1 WHERE id = $2`, amount, toolID)
	return err
}

func (tr *ToolRepositoryPostgres) DecreaseStock(toolID int64, amount int) error {
	_, err := tr.DB.Exec(`UPDATE tools SET stock = stock - $1 WHERE id = $2`, amount, toolID)
	return err
}
