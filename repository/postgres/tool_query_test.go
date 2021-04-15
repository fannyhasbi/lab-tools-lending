package postgres

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCanGetAvailableTool(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	query := NewToolQueryPostgres(db)

	rows := sqlmock.NewRows([]string{"id", "name", "brand", "product_type", "weight", "stock", "additional_info", "created_at", "updated_at"}).
		AddRow(1, "nametest", "brandtest", "producttypetest", 99.0, 10, "additionaltest", timeNowString(), timeNowString())

	mock.ExpectQuery("^SELECT(.*) FROM tools(.*) WHERE stock > 0").
		WillReturnRows(rows)

	result := query.GetAvailableTools()
	assert.NoError(t, result.Error)
	assert.NotEmpty(t, result.Result)
}
