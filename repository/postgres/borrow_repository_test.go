package postgres

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fannyhasbi/lab-tools-lending/types"
	"github.com/stretchr/testify/assert"
)

func TestCanSaveBorrow(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	borrow := types.Borrow{
		ID:         123,
		Amount:     5,
		ReturnDate: "2016-01-01",
		UserID:     111,
		ToolID:     222,
		CreatedAt:  timeNowString(),
	}

	repository := NewBorrowRepositoryPostgres(db)

	rows := sqlmock.NewRows([]string{"id", "amount", "return_date", "user_id", "tool_id", "created_at"}).
		AddRow(borrow.ID, borrow.Amount, borrow.ReturnDate, borrow.UserID, borrow.ToolID, borrow.CreatedAt)

	mock.ExpectQuery("^INSERT INTO borrows (.*) VALUES (.*) RETURNING (.*)").WillReturnRows(rows)

	result, err := repository.Save(&borrow)
	assert.NoError(t, err)
	assert.Equal(t, borrow, result)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
