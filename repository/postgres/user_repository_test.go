package postgres

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fannyhasbi/lab-tools-lending/types"
	"github.com/stretchr/testify/assert"
)

func TestCanSaveUser(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	user := types.User{
		ID:        123,
		Name:      "testname",
		NIM:       "2112xxxxxxxxxx",
		Batch:     2016,
		Address:   "jalan test message",
		CreatedAt: timeNowString(),
	}

	repository := NewUserRepositoryPostgres(db)

	rows := sqlmock.NewRows([]string{"id", "name", "nim", "batch", "address", "created_at"}).
		AddRow(user.ID, user.Name, user.NIM, user.Batch, user.Address, user.CreatedAt)

	mock.ExpectQuery("^INSERT INTO users (.*) VALUES (.*) RETURNING (.*)").
		WillReturnRows(rows)

	result, err := repository.Save(&user)
	assert.NoError(t, err)
	assert.Equal(t, user, result)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCanUpdateUser(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	user := types.User{
		ID:        123,
		Name:      "testname",
		NIM:       "2112xxxxxxxxxx",
		Batch:     2016,
		Address:   "jalan test message",
		CreatedAt: timeNowString(),
	}

	repository := NewUserRepositoryPostgres(db)

	rows := sqlmock.NewRows([]string{"id", "name", "nim", "batch", "address", "created_at"}).
		AddRow(user.ID, user.Name, user.NIM, user.Batch, user.Address, user.CreatedAt)

	mock.ExpectQuery("^UPDATE users SET (.*) WHERE id = (.*)").
		WillReturnRows(rows)

	result, err := repository.Update(&user)
	assert.NoError(t, err)
	assert.Equal(t, user, result)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCanDeleteUser(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	var id int64 = 123

	repository := NewUserRepositoryPostgres(db)

	mock.ExpectExec("^DELETE FROM users WHERE id = (.*)").
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repository.Delete(id)
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
