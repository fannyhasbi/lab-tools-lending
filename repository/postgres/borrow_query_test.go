package postgres

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fannyhasbi/lab-tools-lending/types"
	"github.com/stretchr/testify/assert"
)

func TestCanFindInitialByUserID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	query := NewBorrowQueryPostgres(db)

	var userID int64 = 555
	tt := types.Borrow{
		ID:         123,
		Amount:     1,
		ReturnDate: sql.NullString{Valid: false, String: ""},
		Status:     types.GetBorrowStatus("init"),
		UserID:     111,
		ToolID:     222,
		CreatedAt:  timeNowString(),
	}

	rows := sqlmock.NewRows([]string{"id", "amount", "return_date", "status", "user_id", "tool_id", "created_at"}).
		AddRow(tt.ID, tt.Amount, tt.ReturnDate, tt.Status, tt.UserID, tt.ToolID, tt.CreatedAt)

	mock.ExpectQuery("^SELECT (.+) FROM borrows WHERE user_id = (.+) AND status (.+) ORDER BY id DESC").
		WithArgs(userID, types.GetBorrowStatus("init")).
		WillReturnRows(rows)

	result := query.FindInitialByUserID(userID)
	assert.NoError(t, result.Error)
	assert.NotEmpty(t, result.Result)
	assert.NotPanics(t, func() {
		r := result.Result.(types.Borrow)
		assert.Equal(t, tt, r)
	})
}

func TestCanFindByUserID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	query := NewBorrowQueryPostgres(db)

	var userID int64 = 555
	tt := []types.Borrow{
		{
			ID:         123,
			Amount:     1,
			ReturnDate: sql.NullString{Valid: false, String: ""},
			Status:     types.GetBorrowStatus("init"),
			UserID:     111,
			ToolID:     222,
			CreatedAt:  timeNowString(),
			Tool: types.Tool{
				Name: "Tool Name Test 1",
			},
		},
		{
			ID:         124,
			Amount:     1,
			ReturnDate: sql.NullString{Valid: false, String: ""},
			Status:     types.GetBorrowStatus("progress"),
			UserID:     111,
			ToolID:     223,
			CreatedAt:  timeNowString(),
			Tool: types.Tool{
				Name: "Tool Name Test 2",
			},
		},
	}

	rows := sqlmock.NewRows([]string{"id", "amount", "return_date", "status", "user_id", "tool_id", "created_at", "tool_name"})
	for _, v := range tt {
		rows.AddRow(v.ID, v.Amount, v.ReturnDate, v.Status, v.UserID, v.ToolID, v.CreatedAt, v.Tool.Name)
	}

	mock.ExpectQuery("^SELECT .+ FROM borrows .+ INNER JOIN tools .+ WHERE .+user_id = .+ ORDER BY .+id DESC").
		WithArgs(userID).
		WillReturnRows(rows)

	result := query.FindByUserID(userID)
	assert.NoError(t, result.Error)
	assert.NotEmpty(t, result.Result)
	assert.NotPanics(t, func() {
		r := result.Result.([]types.Borrow)
		assert.Equal(t, tt, r)
	})
}
