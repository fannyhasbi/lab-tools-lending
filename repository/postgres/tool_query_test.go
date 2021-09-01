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

	mock.ExpectQuery("^SELECT (.+) FROM tools WHERE id = (.+)").WillReturnRows(rows)

	result := query.FindByID(tt.ID)
	assert.NoError(t, result.Error)
	assert.NotEmpty(t, result.Result)
	assert.NotPanics(t, func() {
		r := result.Result.(types.Tool)
		assert.Equal(t, tt, r)
	})
}

func TestCanGetAvailableTool(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	query := NewToolQueryPostgres(db)

	rows := sqlmock.NewRows([]string{"id", "name", "brand", "product_type", "weight", "stock", "additional_info", "created_at", "updated_at"}).
		AddRow(1, "nametest", "brandtest", "producttypetest", 99.0, 10, "additionaltest", timeNowString(), timeNowString())

	mock.ExpectQuery("^SELECT (.+) FROM tools WHERE stock > 0 ORDER BY id ASC").
		WillReturnRows(rows)

	result := query.GetAvailableTools()
	assert.NoError(t, result.Error)
	assert.NotEmpty(t, result.Result)
}
