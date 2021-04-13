package postgres

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCanFindChatByID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	var chatID int64 = 123
	query := NewUserQueryPostgres(db)

	rows := sqlmock.NewRows([]string{"id", "chat_id", "name", "nim", "batch", "address", "created_at"}).
		AddRow(1, chatID, "testname", "2111", 2016, "testaddress", timeNowString())

	mock.ExpectQuery("^SELECT(.*)FROM users(.*)WHERE chat_id = (.*)").
		WillReturnRows(rows)

	result := query.FindByChatID(chatID)
	assert.NoError(t, result.Error)
	assert.NotEmpty(t, result.Result)
}
