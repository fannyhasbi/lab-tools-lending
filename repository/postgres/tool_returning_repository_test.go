package postgres

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fannyhasbi/lab-tools-lending/types"
	"github.com/stretchr/testify/assert"
)

func TestCanSaveToolReturning(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	toolReturning := types.ToolReturning{
		ID:             123,
		UserID:         111,
		ToolID:         222,
		Status:         types.GetToolReturningStatus("request"),
		CreatedAt:      timeNowString(),
		AdditionalInfo: "Test additional info.",
	}

	repository := NewToolReturningRepositoryPostgres(db)

	rows := sqlmock.NewRows([]string{"id", "user_id", "tool_id", "status", "created_at", "additional_info"}).
		AddRow(toolReturning.ID, toolReturning.UserID, toolReturning.ToolID, toolReturning.Status, toolReturning.CreatedAt, toolReturning.AdditionalInfo)

	mock.ExpectPrepare("^INSERT INTO tool_returning .+ VALUES .+ RETURNING .+").
		ExpectQuery().
		WithArgs(toolReturning.UserID, toolReturning.ToolID, toolReturning.Status, toolReturning.AdditionalInfo).
		WillReturnRows(rows)

	result, err := repository.Save(&toolReturning)
	assert.NoError(t, err)
	assert.Equal(t, toolReturning, result)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCanUpdateStatusToolReturning(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	var id int64 = 123
	status := types.GetToolReturningStatus("complete")

	repository := NewToolReturningRepositoryPostgres(db)

	mock.ExpectExec("^UPDATE tool_returning SET status = .+ WHERE id = .+").
		WithArgs(status, id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repository.UpdateStatus(id, status)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCanUpdateConfirmToolReturning(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	var id int64 = 123
	now := time.Now()
	confirmedBy := "Test Confirmed By"

	repository := NewToolReturningRepositoryPostgres(db)

	mock.ExpectExec("^UPDATE tool_returning SET confirmed_at = .+ confirmed_by = .+ WHERE id = .+").
		WithArgs(now, confirmedBy, id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repository.UpdateConfirm(id, now, confirmedBy)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
