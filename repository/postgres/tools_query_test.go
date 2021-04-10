package postgres

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCanGetTools(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	query := NewToolsQueryPostgres(db)

	rows := sqlmock.NewRows([]string{"id", "name", "brand", "product_type", "weight", "additional_info"}).
		AddRow(1, "nametest", "brandtest", "producttypetest", 99.0, "additionaltest")

	mock.ExpectQuery("^SELECT(.*)FROM tools(.*)").
		WillReturnRows(rows)

	result := query.GetTools()
	assert.NoError(t, result.Error)
	assert.NotEmpty(t, result.Result)
}
