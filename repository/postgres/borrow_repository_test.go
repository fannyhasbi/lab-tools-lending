package postgres

import (
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
		UserID:    111,
		ToolID:    222,
		CreatedAt: timeNowString(),
	}

	repository := NewBorrowRepositoryPostgres(db)

	rows := sqlmock.NewRows([]string{"id", "amount", "duration", "status", "user_id", "tool_id", "created_at", "confirmed_at"}).
		AddRow(borrow.ID, borrow.Amount, borrow.Duration, borrow.Status, borrow.UserID, borrow.ToolID, borrow.CreatedAt, borrow.ConfirmedAt)

	mock.ExpectQuery("^INSERT INTO borrows (.+) VALUES (.+) RETURNING (.+)").WillReturnRows(rows)

	result, err := repository.Save(&borrow)
	assert.NoError(t, err)
	assert.Equal(t, borrow, result)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCanUpdateBorrow(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	borrow := types.Borrow{
		ID:        123,
		Amount:    5,
		Duration:  7,
		Status:    types.GetBorrowStatus("progress"),
		UserID:    111,
		ToolID:    222,
		CreatedAt: timeNowString(),
	}

	repository := NewBorrowRepositoryPostgres(db)

	rows := sqlmock.NewRows([]string{"id", "amount", "duration", "status", "user_id", "tool_id", "created_at", "confirmed_at"}).
		AddRow(borrow.ID, borrow.Amount, borrow.Duration, borrow.Status, borrow.UserID, borrow.ToolID, borrow.CreatedAt, borrow.ConfirmedAt)

	mock.ExpectQuery("^UPDATE borrows SET .+ WHERE id = .+ RETURNING .+").
		WillReturnRows(rows)

	result, err := repository.Update(&borrow)
	assert.NoError(t, err)
	assert.Equal(t, borrow, result)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCanUpdateBorrowReason(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	var id int64 = 123
	var reason string = "test reason text"

	repository := NewBorrowRepositoryPostgres(db)

	mock.ExpectPrepare("^UPDATE borrows SET reason = .+ WHERE id = .+").
		ExpectExec().
		WithArgs(reason, id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repository.UpdateReason(id, reason)
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCanUpdateBorrowConfirmedAt(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	var id int64 = 123
	var confirmedAt = time.Now()

	repository := NewBorrowRepositoryPostgres(db)

	mock.ExpectExec("^UPDATE borrows SET confirmed_at = .+ WHERE id = .+").
		WithArgs(confirmedAt, id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repository.UpdateConfirmedAt(id, confirmedAt)
	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
