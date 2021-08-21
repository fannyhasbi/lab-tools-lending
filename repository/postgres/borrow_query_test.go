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
		WithArgs(tt.ID, types.GetBorrowStatus("init")).
		WillReturnRows(rows)

	result := query.FindInitialByUserID(tt.ID)
	assert.NoError(t, result.Error)
	assert.NotEmpty(t, result.Result)
	assert.NotPanics(t, func() {
		r := result.Result.(types.Borrow)
		assert.Equal(t, tt, r)
	})
}
