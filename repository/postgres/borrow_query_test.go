package postgres

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fannyhasbi/lab-tools-lending/types"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestCanFindBorrowByID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	query := NewBorrowQueryPostgres(db)

	var id int64 = 123
	borrow := types.Borrow{
		ID:         123,
		Amount:     1,
		ReturnDate: sql.NullString{Valid: false, String: ""},
		Status:     types.GetBorrowStatus("request"),
		UserID:     111,
		ToolID:     222,
		CreatedAt:  timeNowString(),
		Reason:     sql.NullString{Valid: true, String: "test reason"},
		Tool: types.Tool{
			Name:  "Test Tool Name 1",
			Stock: 10,
		},
		User: types.User{
			NIM:  "21120XXXXXXXXX",
			Name: "Test Name",
		},
	}

	rows := sqlmock.NewRows([]string{"id", "amount", "return_date", "status", "user_id", "tool_id", "created_at", "reason", "tool_name", "tool_stock", "user_name", "nim"}).
		AddRow(borrow.ID, borrow.Amount, borrow.ReturnDate, borrow.Status, borrow.UserID, borrow.ToolID, borrow.CreatedAt, borrow.Reason, borrow.Tool.Name, borrow.Tool.Stock, borrow.User.Name, borrow.User.NIM)

	mock.ExpectQuery("^SELECT (.+) FROM borrows .+ INNER JOIN tools .+ INNER JOIN users .+ WHERE .+id = .+").WithArgs(id).WillReturnRows(rows)

	result := query.FindByID(id)
	assert.NoError(t, result.Error)
	assert.NotEmpty(t, result.Result)
	assert.NotPanics(t, func() {
		r := result.Result.(types.Borrow)
		assert.Equal(t, borrow, r)
	})

}

func TestCanFindBorrowByUserIDAndStatus(t *testing.T) {
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
		Tool: types.Tool{
			Name: "Test Tool Name 1",
		},
	}

	rows := sqlmock.NewRows([]string{"id", "amount", "return_date", "status", "user_id", "tool_id", "created_at", "tool_name"}).
		AddRow(tt.ID, tt.Amount, tt.ReturnDate, tt.Status, tt.UserID, tt.ToolID, tt.CreatedAt, tt.Tool.Name)

	mock.ExpectQuery("^SELECT (.+) FROM borrows .+ LEFT JOIN tools .+ WHERE .+user_id = (.+) AND .+status (.+) ORDER BY .+id DESC").
		WithArgs(userID, types.GetBorrowStatus("init")).
		WillReturnRows(rows)

	result := query.FindByUserIDAndStatus(userID, types.GetBorrowStatus("init"))
	assert.NoError(t, result.Error)
	assert.NotEmpty(t, result.Result)
	assert.NotPanics(t, func() {
		r := result.Result.(types.Borrow)
		assert.Equal(t, tt, r)
	})
}

func TestCanFindBorrowByUserID(t *testing.T) {
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

func TestCanGetBorrowsByStatus(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	query := NewBorrowQueryPostgres(db)

	status := types.GetBorrowStatus("init")
	tt := []types.Borrow{
		{
			ID:         123,
			Amount:     1,
			ReturnDate: sql.NullString{Valid: false, String: ""},
			Status:     types.GetBorrowStatus("progress"),
			UserID:     111,
			ToolID:     222,
			CreatedAt:  timeNowString(),
			Tool: types.Tool{
				Name: "Tool Name Test 1",
			},
			User: types.User{
				Name: "Test Name 1",
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
			User: types.User{
				Name: "Test Name 2",
			},
		},
	}

	rows := sqlmock.NewRows([]string{"id", "amount", "return_date", "status", "user_id", "tool_id", "created_at", "tool_name", "user_name"})
	for _, v := range tt {
		rows.AddRow(v.ID, v.Amount, v.ReturnDate, v.Status, v.UserID, v.ToolID, v.CreatedAt, v.Tool.Name, v.User.Name)
	}

	mock.ExpectQuery("^SELECT .+ FROM borrows b INNER JOIN tools t .+ INNER JOIN users u .+ WHERE b.status = .+ ORDER BY b.id ASC").
		WithArgs(status).
		WillReturnRows(rows)

	result := query.GetByStatus(status)
	assert.NoError(t, result.Error)
	assert.NotEmpty(t, result.Result)
	assert.NotPanics(t, func() {
		r := result.Result.([]types.Borrow)
		assert.Equal(t, tt, r)
	})
}

func TestCanGetBorrowByMultipleStatus(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	query := NewBorrowQueryPostgres(db)

	var userID int64 = 111
	status := []types.BorrowStatus{
		types.GetBorrowStatus("request"),
		types.GetBorrowStatus("progress"),
	}
	tt := []types.Borrow{
		{
			ID:         123,
			Amount:     1,
			ReturnDate: sql.NullString{Valid: false, String: ""},
			Status:     types.GetBorrowStatus("progress"),
			UserID:     userID,
			ToolID:     222,
			CreatedAt:  timeNowString(),
			Tool: types.Tool{
				Name: "Tool Name Test 1",
			},
			User: types.User{
				Name: "Test Name 1",
			},
		},
		{
			ID:         124,
			Amount:     1,
			ReturnDate: sql.NullString{Valid: false, String: ""},
			Status:     types.GetBorrowStatus("progress"),
			UserID:     userID,
			ToolID:     223,
			CreatedAt:  timeNowString(),
			Tool: types.Tool{
				Name: "Tool Name Test 2",
			},
			User: types.User{
				Name: "Test Name 2",
			},
		},
	}

	rows := sqlmock.NewRows([]string{"id", "amount", "return_date", "status", "user_id", "tool_id", "created_at", "tool_name", "user_name"})
	for _, v := range tt {
		rows.AddRow(v.ID, v.Amount, v.ReturnDate, v.Status, v.UserID, v.ToolID, v.CreatedAt, v.Tool.Name, v.User.Name)
	}

	mock.ExpectQuery("^SELECT .+ FROM borrows b INNER JOIN tools t .+ INNER JOIN users u .+ WHERE b.user_id = .+ AND b.status = ANY.+ ORDER BY b.id ASC").
		WithArgs(userID, pq.Array(status)).
		WillReturnRows(rows)

	result := query.GetByUserIDAndMultipleStatus(userID, status)
	assert.NoError(t, result.Error)
	assert.NotEmpty(t, result.Result)
	assert.NotPanics(t, func() {
		r := result.Result.([]types.Borrow)
		assert.Equal(t, tt, r)
	})
}
