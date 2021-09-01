package postgres

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fannyhasbi/lab-tools-lending/types"
	"github.com/stretchr/testify/assert"
)

func TestCanFindToolReturningByUserIDAndStatus(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	query := NewToolReturningQueryPostgres(db)

	var userID int64 = 555
	toolReturning := types.ToolReturning{
		ID:             123,
		UserID:         userID,
		ToolID:         111,
		Status:         types.GetToolReturningStatus("progress"),
		ReturnedAt:     timeNowString(),
		AdditionalInfo: "test additional info",
	}

	rows := sqlmock.NewRows([]string{"id", "user_id", "tool_id", "status", "created_at", "additional_info"}).
		AddRow(toolReturning.ID, toolReturning.UserID, toolReturning.ToolID, toolReturning.Status, toolReturning.ReturnedAt, toolReturning.AdditionalInfo)

	mock.ExpectQuery("^SELECT .+ FROM tool_returning WHERE user_id = .+ AND status = .+ ORDER BY id DESC").
		WithArgs(userID, types.GetToolReturningStatus("progress")).
		WillReturnRows(rows)

	result := query.FindByUserIDAndStatus(userID, types.GetToolReturningStatus("progress"))
	assert.NoError(t, result.Error)
	assert.NotEmpty(t, result.Result)
	assert.NotPanics(t, func() {
		r := result.Result.(types.ToolReturning)
		assert.Equal(t, toolReturning, r)
	})
}
