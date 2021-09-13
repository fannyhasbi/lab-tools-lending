package helper

import (
	"database/sql"
	"testing"
	"time"

	"github.com/fannyhasbi/lab-tools-lending/types"
	"github.com/stretchr/testify/assert"
)

func TestCanBuildToolReturningReportMessage(t *testing.T) {
	rets := []types.ToolReturning{
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
				Time:  time.Now(),
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

	r := BuildToolReturningReportMessage(rets)

	// todo: make a better assertion
	assert.Contains(t, r, rets[0].User.Name)
	assert.Contains(t, r, rets[1].User.Name)
}
