package postgres

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fannyhasbi/lab-tools-lending/types"
	"github.com/stretchr/testify/assert"
)

func TestCanFindToolReturningByID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	query := NewToolReturningQueryPostgres(db)

	var id int64 = 555
	toolReturning := types.ToolReturning{
		ID:             123,
		UserID:         111,
		ToolID:         222,
		Status:         types.GetToolReturningStatus("request"),
		CreatedAt:      timeNowString(),
		AdditionalInfo: "test additional info",
		Tool: types.Tool{
			Name: "Test Tool Name 1",
		},
		User: types.User{
			NIM:  "21120XXXXXXXXX",
			Name: "Test Name 1",
		},
	}

	rows := sqlmock.NewRows([]string{"id", "user_id", "tool_id", "status", "created_at", "additional_info", "tool_name", "user_name", "nim"}).
		AddRow(toolReturning.ID, toolReturning.UserID, toolReturning.ToolID, toolReturning.Status, toolReturning.CreatedAt, toolReturning.AdditionalInfo, toolReturning.Tool.Name, toolReturning.User.Name, toolReturning.User.NIM)

	mock.ExpectQuery("^SELECT .+ FROM tool_returning tr INNER JOIN tools t .+ INNER JOIN users u .+ WHERE tr.id = .+").
		WithArgs(id).
		WillReturnRows(rows)

	result := query.FindByID(id)
	assert.NoError(t, result.Error)
	assert.NotEmpty(t, result.Result)
	assert.NotPanics(t, func() {
		r := result.Result.(types.ToolReturning)
		assert.Equal(t, toolReturning, r)
	})
}

func TestCanFindToolReturningByUserIDAndStatus(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	query := NewToolReturningQueryPostgres(db)

	var userID int64 = 555
	toolReturning := types.ToolReturning{
		ID:             123,
		UserID:         userID,
		ToolID:         111,
		Status:         types.GetToolReturningStatus("request"),
		CreatedAt:      timeNowString(),
		AdditionalInfo: "test additional info",
	}

	rows := sqlmock.NewRows([]string{"id", "user_id", "tool_id", "status", "created_at", "additional_info"}).
		AddRow(toolReturning.ID, toolReturning.UserID, toolReturning.ToolID, toolReturning.Status, toolReturning.CreatedAt, toolReturning.AdditionalInfo)

	mock.ExpectQuery("^SELECT .+ FROM tool_returning WHERE user_id = .+ AND status = .+ ORDER BY id DESC").
		WithArgs(userID, types.GetToolReturningStatus("request")).
		WillReturnRows(rows)

	result := query.FindByUserIDAndStatus(userID, types.GetToolReturningStatus("request"))
	assert.NoError(t, result.Error)
	assert.NotEmpty(t, result.Result)
	assert.NotPanics(t, func() {
		r := result.Result.(types.ToolReturning)
		assert.Equal(t, toolReturning, r)
	})
}

func TestCanGetToolReturningByStatus(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	query := NewToolReturningQueryPostgres(db)

	status := types.GetToolReturningStatus("request")
	toolRets := []types.ToolReturning{
		{
			ID:             123,
			UserID:         111,
			ToolID:         222,
			Status:         types.GetToolReturningStatus("request"),
			CreatedAt:      timeNowString(),
			AdditionalInfo: "test additional info",
			Tool: types.Tool{
				Name: "Test Tool Name 1",
			},
			User: types.User{
				Name: "Test Name 1",
			},
		},
		{
			ID:             124,
			UserID:         211,
			ToolID:         312,
			Status:         types.GetToolReturningStatus("request"),
			CreatedAt:      timeNowString(),
			AdditionalInfo: "test additional info",
			Tool: types.Tool{
				Name: "Test Tool Name 2",
			},
			User: types.User{
				Name: "Test Name 2",
			},
		},
	}

	rows := sqlmock.NewRows([]string{"id", "user_id", "tool_id", "status", "created_at", "additional_info", "tool_name", "user_name"})
	for _, v := range toolRets {
		rows.AddRow(v.ID, v.UserID, v.ToolID, v.Status, v.CreatedAt, v.AdditionalInfo, v.Tool.Name, v.User.Name)
	}

	mock.ExpectQuery("^SELECT .+ FROM tool_returning tr INNER JOIN tools t .+ INNER JOIN users u .+ WHERE tr.status = .+ ORDER BY tr.id ASC").
		WithArgs(status).
		WillReturnRows(rows)

	result := query.GetByStatus(status)
	assert.NoError(t, result.Error)
	assert.NotEmpty(t, result.Result)
	assert.NotPanics(t, func() {
		r := result.Result.([]types.ToolReturning)
		assert.Equal(t, toolRets, r)
	})
}
