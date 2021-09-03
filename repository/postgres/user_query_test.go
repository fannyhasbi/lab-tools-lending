package postgres

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCanFindUserByID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	var id int64 = 123
	query := NewUserQueryPostgres(db)

	rows := sqlmock.NewRows([]string{"id", "name", "nim", "batch", "address", "created_at"}).
		AddRow(id, "testname", "2111", 2016, "testaddress", timeNowString())

	mock.ExpectQuery("^SELECT(.+)FROM users(.+)WHERE id = (.+)").
		WithArgs(id).
		WillReturnRows(rows)

	result := query.FindByID(id)
	assert.NoError(t, result.Error)
	assert.NotEmpty(t, result.Result)
}
