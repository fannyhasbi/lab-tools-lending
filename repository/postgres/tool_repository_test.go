package postgres

import (
	"testing"
	"time"

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

func TestCanUpdateTool(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	tool := types.Tool{
		ID:                    123,
		Name:                  "Test Name",
		Brand:                 "Test Brand",
		ProductType:           "Test Product Type",
		Weight:                120,
		Stock:                 23,
		AdditionalInformation: "test additional info",
	}

	repository := NewToolRepositoryPostgres(db)

	mock.ExpectPrepare("^UPDATE tools SET .+ WHERE id = .+").
		ExpectExec().
		WithArgs(tool.Name, tool.Brand, tool.ProductType, tool.Weight, tool.Stock, tool.AdditionalInformation, tool.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repository.Update(&tool)
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCanDeleteTool(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	currentTime := time.Now()
	var toolID int64 = 123

	repository := NewToolRepositoryPostgres(db)

	mock.ExpectPrepare("^UPDATE tools SET deleted_at = .* WHERE id = .+").
		ExpectExec().
		WithArgs(currentTime, toolID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repository.Delete(toolID, currentTime)
	assert.NoError(t, err)
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

func TestCanDeleteToolPhotos(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	var id int64 = 555
	repository := NewToolRepositoryPostgres(db)

	mock.ExpectExec("^DELETE FROM tool_photos WHERE tool_id = .+").
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repository.DeletePhotos(id)
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCanIncreaseStock(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	var id int64 = 123
	amount := 3
	repository := NewToolRepositoryPostgres(db)

	mock.ExpectExec("^UPDATE tools SET stock = stock \\+ .+ WHERE id = .+").
		WithArgs(amount, id).WillReturnResult(sqlmock.NewResult(1, 1))

	err := repository.IncreaseStock(id, amount)
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCanDecreaseStock(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	var id int64 = 123
	amount := 3
	repository := NewToolRepositoryPostgres(db)

	mock.ExpectExec("^UPDATE tools SET stock = stock - .+ WHERE id = .+").
		WithArgs(amount, id).WillReturnResult(sqlmock.NewResult(1, 1))

	err := repository.DecreaseStock(id, amount)
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
