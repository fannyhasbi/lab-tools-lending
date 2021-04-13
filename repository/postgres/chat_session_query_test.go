package postgres

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fannyhasbi/lab-tools-lending/types"
	"github.com/stretchr/testify/assert"
)

func timeNow() string {
	return time.Now().Format(time.RFC3339)
}

func TestCanGetChatSession(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	user := types.User{
		ID:     123,
		ChatID: 321,
	}

	query := NewChatSessionQueryPostgres(db)

	rows := sqlmock.NewRows([]string{"id", "status", "created_at"}).
		AddRow(1, types.ChatSessionStatus["progress"], timeNow())

	mock.ExpectQuery("^SELECT(.*)FROM chat_sessions(.*)WHERE user_id = (.*) ORDER BY id DESC").
		WillReturnRows(rows)

	result := query.Get(user)
	assert.NoError(t, result.Error)
	assert.NotEmpty(t, result.Result)
}

func TestCanGetChatSessionDetail(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	chatSession := types.ChatSession{
		ID: 123,
	}

	query := NewChatSessionQueryPostgres(db)

	rows := sqlmock.NewRows([]string{"id", "topic", "created_at"}).
		AddRow(1, types.ChatSessionTopic["register"], timeNow())

	mock.ExpectQuery("^SELECT(.*)FROM chat_session_details(.*)WHERE chat_session_id = (.*) ORDER BY id DESC").
		WillReturnRows(rows)

	result := query.GetDetail(chatSession)
	assert.NoError(t, result.Error)
	assert.NotEmpty(t, result.Result)
}
