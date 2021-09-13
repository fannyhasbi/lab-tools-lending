package postgres

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fannyhasbi/lab-tools-lending/types"
	"github.com/stretchr/testify/assert"
)

func TestCanSaveBorrow(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	borrow := types.Borrow{
		ID:        123,
		Amount:    5,
		Duration:  7,
		Status:    types.GetBorrowStatus("progress"),
		Reason:    sql.NullString{Valid: true, String: "test borrow reason"},
		UserID:    111,
		ToolID:    222,
		CreatedAt: timeNowString(),
	}

	repository := NewBorrowRepositoryPostgres(db)

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(borrow.ID)

	mock.ExpectQuery("^INSERT INTO borrows (.+) VALUES (.+) RETURNING id").
		WithArgs(borrow.Amount, borrow.Status, borrow.UserID, borrow.ToolID, borrow.Reason, borrow.Duration).
		WillReturnRows(rows)

	result, err := repository.Save(&borrow)
	assert.NoError(t, err)
	assert.Equal(t, borrow.ID, result)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCanUpdateBorrowStatus(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	var id int64 = 123
	status := types.GetBorrowStatus("request")

	repository := NewBorrowRepositoryPostgres(db)

	mock.ExpectExec("^UPDATE borrows SET status = .+ WHERE id = .+").
		WithArgs(status, id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repository.UpdateStatus(id, status)
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCanUpdateBorrowConfirm(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	var id int64 = 123
	confirmedAt := time.Now()
	confirmedBy := "Test Confirmed By"

	repository := NewBorrowRepositoryPostgres(db)

	mock.ExpectExec("^UPDATE borrows SET confirmed_at = .+ confirmed_by = .+ WHERE id = .+").
		WithArgs(confirmedAt, confirmedBy, id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repository.UpdateConfirm(id, confirmedAt, confirmedBy)
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
