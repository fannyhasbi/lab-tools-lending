package types

import "database/sql"

type (
	BorrowStatus string

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

var (
	BorrowTimeRangeMap map[string]int = map[string]int{
		"oneweek":  7,
		"twoweek":  14,
		"onemonth": 30,
		"twomonth": 60,
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
