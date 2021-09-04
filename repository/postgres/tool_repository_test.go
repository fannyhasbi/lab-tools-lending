package postgres

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fannyhasbi/lab-tools-lending/types"
	"github.com/stretchr/testify/assert"
)

func TestCanSaveTool(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	var id int64 = 123
	tool := types.Tool{
		Name:                  "Test Name",
		Brand:                 "Test Brand",
		ProductType:           "Test Product Type",
		Weight:                120,
		Stock:                 23,
		AdditionalInformation: "test additional info",
	}

	repository := NewToolRepositoryPostgres(db)

	rows := sqlmock.NewRows([]string{"id"}).AddRow(id)

	mock.ExpectPrepare("^INSERT INTO tools .+ VALUES .+ RETURNING id").
		ExpectQuery().
		WithArgs(tool.Name, tool.Brand, tool.ProductType, tool.Weight, tool.Stock, tool.AdditionalInformation).
		WillReturnRows(rows)

	result, err := repository.Save(&tool)
	assert.NoError(t, err)
	assert.Equal(t, id, result)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCanSaveToolPhotos(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	var toolID int64 = 123
	photos := []types.TelePhotoSize{
		{
			FileID:       "abc123",
			FileUniqueID: "123abc",
		},
		{
			FileID:       "xyz456",
			FileUniqueID: "456xyz",
		},
	}

	repository := NewToolRepositoryPostgres(db)

	mock.ExpectExec(`^INSERT INTO tool_photos \(tool_id,file_id,file_unique_id\) VALUES \(.+,.+,.+\),\(.+,.+,.+\)`).
		WithArgs(toolID, photos[0].FileID, photos[0].FileUniqueID, toolID, photos[1].FileID, photos[1].FileUniqueID).
		WillReturnResult(sqlmock.NewResult(2, 2))

	err := repository.SavePhotos(toolID, photos)
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCanIncreaseStock(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	var id int64 = 123
	repository := NewToolRepositoryPostgres(db)

	mock.ExpectExec("^UPDATE tools SET stock = stock \\+ 1 WHERE id = .+").
		WithArgs(id).WillReturnResult(sqlmock.NewResult(1, 1))

	err := repository.IncreaseStock(id)
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCanDecreaseStock(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	var id int64 = 123
	repository := NewToolRepositoryPostgres(db)

	mock.ExpectExec("^UPDATE tools SET stock = stock - 1 WHERE id = .+").
		WithArgs(id).WillReturnResult(sqlmock.NewResult(1, 1))

	err := repository.DecreaseStock(id)
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
