package helper

import (
	"fmt"
	"testing"
	"time"

	"github.com/fannyhasbi/lab-tools-lending/types"
	"github.com/stretchr/testify/assert"
)

func TestBorrowStatusGrouping(t *testing.T) {
	// testInit := []types.Borrow{{ID: 10, Status: types.GetBorrowStatus("init")}}
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
	// tt = append(tt, testInit...)
	tt = append(tt, testProgress...)
	tt = append(tt, testReturned...)

	r := GroupBorrowStatus(tt)

	assert.Equal(t, testProgress, r[types.GetBorrowStatus("progress")])
	assert.Equal(t, testReturned, r[types.GetBorrowStatus("returned")])
	assert.Empty(t, r[types.GetBorrowStatus("init")])
}

func TestBuildBorrowedMessage(t *testing.T) {
	tt := []types.Borrow{
		{
			CreatedAt: "2021-08-01",
			Tool: types.Tool{
				Name: "Tool Name Test 1",
			},
		},
		{
			CreatedAt: time.Now().Format(time.RFC3339),
			Tool: types.Tool{
				Name: "Tool Name Test 2",
			},
		},
	}

	r := BuildBorrowedMessage(tt)

	expected := fmt.Sprintf("* %s (%s)\n* %s (%s)\n",
		tt[0].Tool.Name, TranslateDateStringToBahasa(tt[0].CreatedAt),
		tt[1].Tool.Name, TranslateDateStringToBahasa(tt[1].CreatedAt))

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
