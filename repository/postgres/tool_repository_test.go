package postgres

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

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
