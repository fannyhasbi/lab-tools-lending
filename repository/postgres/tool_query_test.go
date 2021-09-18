package postgres

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fannyhasbi/lab-tools-lending/types"
	"github.com/stretchr/testify/assert"
)

func TestCanFindToolByID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	query := NewToolQueryPostgres(db)

	tt := types.Tool{
		ID:                    1,
		Name:                  "nametest",
		Brand:                 "brandtest",
		ProductType:           "producttypetest",
		Weight:                99.0,
		Stock:                 10,
		AdditionalInformation: "additionaltest",
		CreatedAt:             timeNowString(),
		UpdatedAt:             timeNowString(),
	}

	rows := sqlmock.NewRows([]string{"id", "name", "brand", "product_type", "weight", "stock", "additional_info", "created_at", "updated_at"}).
		AddRow(tt.ID, tt.Name, tt.Brand, tt.ProductType, tt.Weight, tt.Stock, tt.AdditionalInformation, tt.CreatedAt, tt.UpdatedAt)

	mock.ExpectQuery("^SELECT (.+) FROM tools WHERE id = (.+) AND deleted_at IS NULL").
		WithArgs(tt.ID).
		WillReturnRows(rows)

	result := query.FindByID(tt.ID)
	assert.NoError(t, result.Error)
	assert.NotEmpty(t, result.Result)
	assert.NotPanics(t, func() {
		r := result.Result.(types.Tool)
		assert.Equal(t, tt, r)
	})
}

func TestCanGetTools(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	query := NewToolQueryPostgres(db)

	tools := []types.Tool{
		{
			ID:                    1,
			Name:                  "Test Name 1",
			Brand:                 "Test Brand 1",
			ProductType:           "test type 1",
			Weight:                123,
			Stock:                 10,
			AdditionalInformation: "test additional info 1",
			CreatedAt:             timeNowString(),
			UpdatedAt:             timeNowString(),
		},
		{
			ID:                    2,
			Name:                  "Test Name 2",
			Brand:                 "Test Brand 2",
			ProductType:           "test type 2",
			Weight:                321,
			Stock:                 100,
			AdditionalInformation: "test additional info 2",
			CreatedAt:             timeNowString(),
			UpdatedAt:             timeNowString(),
		},
	}

	rows := sqlmock.NewRows([]string{"id", "name", "brand", "product_type", "weight", "stock", "additional_info", "created_at", "updated_at"})
	for _, v := range tools {
		rows.AddRow(v.ID, v.Name, v.Brand, v.ProductType, v.Weight, v.Stock, v.AdditionalInformation, v.CreatedAt, v.UpdatedAt)
	}

	mock.ExpectQuery("^SELECT .+ FROM tools WHERE deleted_at IS NULL ORDER BY id ASC").WillReturnRows(rows)

	result := query.Get()
	assert.NoError(t, result.Error)
	assert.NotEmpty(t, result.Result)
	assert.NotPanics(t, func() {
		r := result.Result.([]types.Tool)
		assert.Equal(t, tools, r)
	})
}

func TestCanGetAvailableTool(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	query := NewToolQueryPostgres(db)

	rows := sqlmock.NewRows([]string{"id", "name", "brand", "product_type", "weight", "stock", "additional_info", "created_at", "updated_at"}).
		AddRow(1, "nametest", "brandtest", "producttypetest", 99.0, 10, "additionaltest", timeNowString(), timeNowString())

	mock.ExpectQuery("^SELECT (.+) FROM tools WHERE stock > 0 AND deleted_at IS NULL ORDER BY id ASC").
		WillReturnRows(rows)

	result := query.GetAvailableTools()
	assert.NoError(t, result.Error)
	assert.NotEmpty(t, result.Result)
}

func TestCanGetPhotos(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	query := NewToolQueryPostgres(db)

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

	rows := sqlmock.NewRows([]string{"file_id", "file_unique_id"}).AddRow(photos[0].FileID, photos[0].FileUniqueID).AddRow(photos[1].FileID, photos[1].FileUniqueID)

	mock.ExpectQuery("^SELECT p.file_id, p.file_unique_id FROM tool_photos p INNER JOIN tools t .+ WHERE p.tool_id = .+ AND t.deleted_at IS NULL").WithArgs(toolID).WillReturnRows(rows)

	result := query.GetPhotos(toolID)
	assert.NoError(t, result.Error)
	assert.NotEmpty(t, result.Result)
	assert.NotPanics(t, func() {
		r := result.Result.([]types.TelePhotoSize)
		assert.Equal(t, photos, r)
	})
}
