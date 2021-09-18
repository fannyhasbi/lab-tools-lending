package postgres

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fannyhasbi/lab-tools-lending/types"
	"github.com/stretchr/testify/assert"
)

func TestCanSaveChatSession(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	requestType := types.RequestTypePrivate
	chatSession := types.ChatSession{
		ID:          123,
		Status:      types.ChatSessionStatus["progress"],
		UserID:      321,
		CreatedAt:   timeNowString(),
		UpdatedAt:   timeNowString(),
		RequestType: requestType,
	}

	repository := NewChatSessionRepositoryPostgres(db)

	rows := sqlmock.NewRows([]string{"id", "status", "user_id", "created_at", "updated_at", "request_type"}).
		AddRow(chatSession.ID, chatSession.Status, chatSession.UserID, chatSession.CreatedAt, chatSession.UpdatedAt, chatSession.RequestType)

	mock.ExpectQuery("^INSERT INTO chat_sessions (.+) VALUES (.+) RETURNING (.+)").
		WithArgs(chatSession.Status, chatSession.UserID, requestType).
		WillReturnRows(rows)

	result, err := repository.Save(&chatSession, requestType)
	assert.NoError(t, err)
	assert.Equal(t, chatSession, result)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCanUpdateChatSession(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	var id int64 = 123
	status := types.ChatSessionStatus["complete"]

	repository := NewChatSessionRepositoryPostgres(db)

	mock.ExpectExec("^UPDATE chat_sessions SET status = (.+) WHERE id = (.+)").
		WithArgs(status, id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repository.UpdateStatus(id, status)
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCanDeleteChatSession(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	var id int64 = 123

	repository := NewChatSessionRepositoryPostgres(db)

	mock.ExpectExec("^DELETE FROM chat_sessions WHERE id = (.+)").
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repository.Delete(id)
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCanSaveChatSessionDetail(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	detail := types.ChatSessionDetail{
		ID:            123,
		Topic:         types.Topic["register_init"],
		ChatSessionID: 321,
		CreatedAt:     timeNowString(),
		Data:          `{"type": "typetest"}`,
	}

	repository := NewChatSessionRepositoryPostgres(db)

	rows := sqlmock.NewRows([]string{"id", "topic", "chat_session_id", "created_at", "data"}).
		AddRow(detail.ID, detail.Topic, detail.ChatSessionID, detail.CreatedAt, detail.Data)

	mock.ExpectQuery("^INSERT INTO chat_session_details (.+) VALUES (.+) RETURNING (.+)").
		WithArgs(detail.Topic, detail.ChatSessionID, detail.Data).
		WillReturnRows(rows)

	result, err := repository.SaveDetail(&detail)
	assert.NoError(t, err)
	assert.Equal(t, detail, result)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCanDeleteChatSessionDetailByChatSessionID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	var id int64 = 123

	repository := NewChatSessionRepositoryPostgres(db)

	mock.ExpectExec("^DELETE FROM chat_session_details WHERE chat_session_id = (.+)").
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repository.DeleteDetailByChatSessionID(id)
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
