package postgres

import (
	"database/sql"
	"testing"
	"time"

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
		BorrowID:       111,
		Status:         types.GetToolReturningStatus("request"),
		CreatedAt:      timeNowString(),
		AdditionalInfo: "test additional info",
		Borrow: types.Borrow{
			UserID:   321,
			Amount:   3,
			Duration: 18,
			ToolID:   999,
			ConfirmedAt: sql.NullTime{
				Valid: true,
				Time:  time.Now(),
			},
			Tool: types.Tool{
				Name: "Test Tool Name 1",
			},
			User: types.User{
				NIM:  "21120XXXXXXXXX",
				Name: "Test Name 1",
			},
		},
	}

	rows := sqlmock.NewRows([]string{"id", "borrow_id", "status", "created_at", "additional_info", "amount", "duration", "tool_id", "borrow_confirmed_at", "tool_name", "user_id", "user_name", "nim"}).
		AddRow(toolReturning.ID, toolReturning.BorrowID, toolReturning.Status, toolReturning.CreatedAt, toolReturning.AdditionalInfo, toolReturning.Borrow.Amount, toolReturning.Borrow.Duration, toolReturning.Borrow.ToolID, toolReturning.Borrow.ConfirmedAt, toolReturning.Borrow.Tool.Name, toolReturning.Borrow.UserID, toolReturning.Borrow.User.Name, toolReturning.Borrow.User.NIM)

	mock.ExpectQuery("^SELECT .+ FROM tool_returning tr INNER JOIN borrows b .+ INNER JOIN tools t .+ INNER JOIN users u .+ WHERE tr.id = .+").
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
		ID:       123,
		BorrowID: 111,
		Borrow: types.Borrow{
			UserID: userID,
		},
		Status:         types.GetToolReturningStatus("request"),
		CreatedAt:      timeNowString(),
		AdditionalInfo: "test additional info",
	}

	rows := sqlmock.NewRows([]string{"id", "borrow_id", "status", "created_at", "additional_info", "user_id"}).
		AddRow(toolReturning.ID, toolReturning.BorrowID, toolReturning.Status, toolReturning.CreatedAt, toolReturning.AdditionalInfo, toolReturning.Borrow.UserID)

	mock.ExpectQuery("^SELECT .+ FROM tool_returning tr INNER JOIN borrows b .+ WHERE b.user_id = .+ AND tr.status = .+ ORDER BY tr.id DESC").
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
			BorrowID:       111,
			Status:         types.GetToolReturningStatus("request"),
			CreatedAt:      timeNowString(),
			AdditionalInfo: "test additional info",
			Borrow: types.Borrow{
				Tool: types.Tool{
					Name: "Test Tool Name 1",
				},
				User: types.User{
					Name: "Test Name 1",
				},
			},
		},
		{
			ID:             124,
			BorrowID:       211,
			Status:         types.GetToolReturningStatus("request"),
			CreatedAt:      timeNowString(),
			AdditionalInfo: "test additional info",
			Borrow: types.Borrow{
				Tool: types.Tool{
					Name: "Test Tool Name 2",
				},
				User: types.User{
					Name: "Test Name 2",
				},
			},
		},
	}

	rows := sqlmock.NewRows([]string{"id", "borrow_id", "status", "created_at", "additional_info", "tool_name", "user_name"})
	for _, v := range toolRets {
		rows.AddRow(v.ID, v.BorrowID, v.Status, v.CreatedAt, v.AdditionalInfo, v.Borrow.Tool.Name, v.Borrow.User.Name)
	}

	mock.ExpectQuery("^SELECT .+ FROM tool_returning tr INNER JOIN borrows b .+ INNER JOIN tools t .+ INNER JOIN users u .+ WHERE tr.status = .+ ORDER BY tr.id ASC").
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

func TestCangGetToolReturningReport(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	query := NewToolReturningQueryPostgres(db)

	year := 2021
	month := 8
	toolRets := []types.ToolReturning{
		{
			ID:          123,
			BorrowID:    111,
			Status:      types.GetToolReturningStatus("request"),
			CreatedAt:   timeNowString(),
			ConfirmedAt: sql.NullTime{Valid: true, Time: time.Now()},
			ConfirmedBy: sql.NullString{Valid: true, String: "Test Confirmed By 1"},
			Borrow: types.Borrow{
				Tool: types.Tool{
					Name: "Test Tool Name 1",
				},
				User: types.User{
					Name: "Test Name 1",
				},
			},
		},
		{
			ID:          124,
			BorrowID:    211,
			Status:      types.GetToolReturningStatus("request"),
			CreatedAt:   timeNowString(),
			ConfirmedAt: sql.NullTime{Valid: true, Time: time.Now()},
			ConfirmedBy: sql.NullString{Valid: true, String: "Test Confirmed By 2"},
			Borrow: types.Borrow{
				Tool: types.Tool{
					Name: "Test Tool Name 2",
				},
				User: types.User{
					Name: "Test Name 2",
				},
			},
		},
	}

	rows := sqlmock.NewRows([]string{"id", "borrow_id", "status", "created_at", "confirmed_at", "confirmed_by", "tool_name", "user_name"})
	for _, v := range toolRets {
		rows.AddRow(v.ID, v.BorrowID, v.Status, v.CreatedAt, v.ConfirmedAt, v.ConfirmedBy, v.Borrow.Tool.Name, v.Borrow.User.Name)
	}

	mock.ExpectQuery(`^SELECT .+ FROM tool_returning tr INNER JOIN borrows b .+ INNER JOIN tools t .+ INNER JOIN users u .+ WHERE tr.status = .+ AND DATE_PART\('year', tr.confirmed_at\) = .+ AND DATE_PART\('month', tr.confirmed_at\) = .+ ORDER BY tr.id ASC`).
		WithArgs(types.GetToolReturningStatus("complete"), year, month).
		WillReturnRows(rows)

	result := query.GetReport(year, month)
	assert.NoError(t, result.Error)
	assert.NotEmpty(t, result.Result)
	assert.NotPanics(t, func() {
		r := result.Result.([]types.ToolReturning)
		assert.Equal(t, toolRets, r)
	})
}
