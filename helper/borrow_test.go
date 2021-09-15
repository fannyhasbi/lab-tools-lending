package helper

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/fannyhasbi/lab-tools-lending/types"
	"github.com/stretchr/testify/assert"
)

func TestBorrowStatusGrouping(t *testing.T) {
	testProgress := []types.Borrow{
		{
			ID:     20,
			Status: types.GetBorrowStatus("progress"),
		},
		{
			ID:     30,
			Status: types.GetBorrowStatus("progress"),
		},
	}
	testReturned := []types.Borrow{
		{
			ID:     40,
			Status: types.GetBorrowStatus("returned"),
		},
	}

	var tt []types.Borrow
	tt = append(tt, testProgress...)
	tt = append(tt, testReturned...)

	r := GroupBorrowStatus(tt)

	assert.Equal(t, testProgress, r[types.GetBorrowStatus("progress")])
	assert.Equal(t, testReturned, r[types.GetBorrowStatus("returned")])
	assert.Empty(t, r[types.GetBorrowStatus("request")])
}

func TestCanGetBorrowByStatus(t *testing.T) {
	borrows := []types.Borrow{
		{
			ID:     1,
			Status: types.GetBorrowStatus("request"),
		},
		{
			ID:     2,
			Status: types.GetBorrowStatus("request"),
		},
		{
			ID:     3,
			Status: types.GetBorrowStatus("progress"),
		},
		{
			ID:     4,
			Status: types.GetBorrowStatus("returned"),
		},
	}

	t.Run("request", func(t *testing.T) {
		r := GetBorrowsByStatus(borrows, types.GetBorrowStatus("request"))
		expected := []types.Borrow{borrows[0], borrows[1]}
		assert.Equal(t, expected, r)
	})

	t.Run("progress", func(t *testing.T) {
		r := GetBorrowsByStatus(borrows, types.GetBorrowStatus("progress"))
		expected := []types.Borrow{borrows[2]}
		assert.Equal(t, expected, r)
	})
}

func TestBuildBorrowedMessage(t *testing.T) {
	b := []types.Borrow{
		{
			ID:       1,
			Duration: 10,
			ConfirmedAt: sql.NullTime{
				Valid: true,
				Time:  time.Now(),
			},
			Tool: types.Tool{
				Name: "Test Tool Name 1",
			},
		},
		{
			ID:       2,
			Duration: 13,
			ConfirmedAt: sql.NullTime{
				Valid: true,
				Time:  time.Now(),
			},
			Tool: types.Tool{
				Name: "Test Tool Name 2",
			},
		},
	}

	r := BuildBorrowedMessage(b)

	layout := "02/01/2006"
	expected := fmt.Sprintf("[%d] %s (%s - %s)\n[%d] %s (%s - %s)\n",
		b[0].ID, b[0].Tool.Name, b[0].ConfirmedAt.Time.Format(layout), b[0].ConfirmedAt.Time.AddDate(0, 0, b[0].Duration).Format(layout),
		b[1].ID, b[1].Tool.Name, b[1].ConfirmedAt.Time.Format(layout), b[1].ConfirmedAt.Time.AddDate(0, 0, b[1].Duration).Format(layout))

	assert.Equal(t, expected, r)
}

func TestBuildBorrowRequestMessage(t *testing.T) {
	borrows := []types.Borrow{
		{
			ID: 123,
			Tool: types.Tool{
				Name: "Test Tool Name 1",
			},
			User: types.User{
				Name: "Test Name 1",
			},
		},
		{
			ID: 321,
			Tool: types.Tool{
				Name: "Test Tool Name 2",
			},
			User: types.User{
				Name: "Test Name 2",
			},
		},
	}

	r := BuildBorrowRequestListMessage(borrows)

	expected := fmt.Sprintf("[%d] %s - %s\n[%d] %s - %s\n",
		borrows[0].ID, borrows[0].User.Name, borrows[0].Tool.Name,
		borrows[1].ID, borrows[1].User.Name, borrows[1].Tool.Name)

	assert.Equal(t, expected, r)
}

func TestBuildToolReturningRequestMessage(t *testing.T) {
	rets := []types.ToolReturning{
		{
			ID: 123,
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
			ID: 321,
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

	r := BuildToolReturningRequestListMessage(rets)

	expected := fmt.Sprintf("[%d] %s - %s\n[%d] %s - %s\n",
		rets[0].ID, rets[0].Borrow.User.Name, rets[0].Borrow.Tool.Name,
		rets[1].ID, rets[1].Borrow.User.Name, rets[1].Borrow.Tool.Name)

	assert.Equal(t, expected, r)
}

func TestGetBorrowFromChatSessionDetail(t *testing.T) {
	var toolID int64 = 123
	var amount int = 3
	var duration int = 23
	var reason string = "test borrow reason"

	t.Run("full borrow session", func(t *testing.T) {
		borrows := []types.ChatSessionDetail{
			{
				Topic: types.Topic["borrow_init"],
				Data:  NewSessionDataGenerator().BorrowInit(toolID),
			},
			{
				Topic: types.Topic["borrow_amount"],
				Data:  NewSessionDataGenerator().BorrowAmount(amount),
			},
			{
				Topic: types.Topic["borrow_date"],
				Data:  NewSessionDataGenerator().BorrowDuration(duration),
			},
			{
				Topic: types.Topic["borrow_reason"],
				Data:  NewSessionDataGenerator().BorrowReason(reason),
			},
		}

		r := GetBorrowFromChatSessionDetail(borrows)

		expected := types.Borrow{
			ToolID:   toolID,
			Amount:   amount,
			Duration: duration,
			Reason: sql.NullString{
				Valid:  true,
				String: reason,
			},
		}

		assert.Equal(t, expected, r)
	})

	t.Run("not full session", func(t *testing.T) {
		borrows := []types.ChatSessionDetail{
			{
				Topic: types.Topic["borrow_init"],
				Data:  NewSessionDataGenerator().BorrowInit(toolID),
			},
		}

		r := GetBorrowFromChatSessionDetail(borrows)

		expected := types.Borrow{
			ToolID:   toolID,
			Amount:   0,
			Duration: 0,
			Reason: sql.NullString{
				Valid:  false,
				String: "",
			},
		}

		assert.Equal(t, expected, r)
	})
}

func TestCanGetSameBorrow(t *testing.T) {
	t.Run("same", func(t *testing.T) {
		var toolID int64 = 123
		borrows := []types.Borrow{
			{
				ToolID: 333,
				Status: types.GetBorrowStatus("progress"),
			},
			{
				ToolID: toolID,
				Status: types.GetBorrowStatus("request"),
			},
		}

		r, same := GetSameBorrow(borrows, toolID)
		assert.True(t, same)
		assert.Equal(t, types.GetBorrowStatus("request"), r)
	})
	t.Run("not same", func(t *testing.T) {
		var toolID int64 = 123
		borrows := []types.Borrow{
			{
				ToolID: 111,
				Status: types.GetBorrowStatus("progress"),
			},
			{
				ToolID: 999,
				Status: types.GetBorrowStatus("request"),
			},
		}

		r, same := GetSameBorrow(borrows, toolID)
		assert.False(t, same)
		assert.Equal(t, types.BorrowStatus(""), r)
	})
}

func TestCanBuildBorrowReportMessage(t *testing.T) {
	borrows := []types.Borrow{
		{
			ID: 1,
			ConfirmedAt: sql.NullTime{
				Valid: true,
				Time:  time.Now(),
			},
			ConfirmedBy: sql.NullString{
				Valid:  true,
				String: "Test confirmed by 1",
			},
			User: types.User{
				Name: "Test User 1",
			},
			Tool: types.Tool{
				Name: "Test Tool 1",
			},
		},
		{
			ID: 2,
			ConfirmedAt: sql.NullTime{
				Valid: true,
				Time:  time.Now().Add(time.Hour * 24 * 2),
			},
			ConfirmedBy: sql.NullString{
				Valid:  true,
				String: "Test confirmed by 2",
			},
			User: types.User{
				Name: "Test User 2",
			},
			Tool: types.Tool{
				Name: "Test Tool 2",
			},
		},
	}

	r := BuildBorrowReportMessage(borrows)

	// todo: make a better assertion
	assert.Contains(t, r, borrows[0].User.Name)
	assert.Contains(t, r, borrows[1].User.Name)
}
