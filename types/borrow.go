package types

import "database/sql"

type (
	BorrowStatus    string
	BorrowTimeRange string

	Borrow struct {
		ID         int64          `json:"id"`
		Amount     int            `json:"amount"`
		ReturnDate sql.NullString `json:"return_date"`
		Status     BorrowStatus   `json:"status"`
		UserID     int64          `json:"user_id"`
		ToolID     int64          `json:"tool_id"`
		CreatedAt  string         `json:"created_at"`
	}
)

func GetBorrowStatus(s string) BorrowStatus {
	status := map[string]BorrowStatus{
		"init":     "INIT",
		"progress": "PROGRESS",
		"returned": "RETURNED",
	}

	return status[s]
}

func GetBorrowTimeRange(tr string) BorrowTimeRange {
	timeRange := map[string]BorrowTimeRange{
		"1week":  "1WEEK",
		"2week":  "2WEEK",
		"1month": "1MONTH",
		"2month": "2MONTH",
	}

	return timeRange[tr]
}
