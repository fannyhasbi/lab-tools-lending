package postgres

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fannyhasbi/lab-tools-lending/types"
	"github.com/stretchr/testify/assert"
)

func TestCanGetChatSession(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	requestType := types.RequestTypePrivate
	user := types.User{
		ID: 123,
	}

	query := NewChatSessionQueryPostgres(db)

	rows := sqlmock.NewRows([]string{"id", "status", "created_at", "request_type"}).
		AddRow(1, types.ChatSessionStatus["progress"], timeNowString(), requestType)

	mock.ExpectQuery("^SELECT(.+)FROM chat_sessions(.+)WHERE user_id = (.+) ORDER BY id DESC").
		WithArgs(user.ID, types.ChatSessionStatus["progress"], requestType).
		WillReturnRows(rows)

	result := query.Get(user, types.RequestTypePrivate)
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

	rows := sqlmock.NewRows([]string{"id", "topic", "chat_session_id", "created_at", "data"}).
		AddRow(1, types.Topic["register"], chatSession.ID, timeNowString(), `{"type": "something"}`)

	mock.ExpectQuery("^SELECT(.+)FROM chat_session_details(.+)WHERE chat_session_id = (.+) ORDER BY id DESC").
		WithArgs(chatSession.ID).
		WillReturnRows(rows)

	result := query.GetDetail(chatSession)
	assert.NoError(t, result.Error)
	assert.NotEmpty(t, result.Result)
	assert.NotPanics(t, func() {
		r := result.Result.([]types.ChatSessionDetail)
		assert.Equal(t, len(r), 1)
	})
}
