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
		Tool       Tool           `json:"tool"`
		User       User           `json:"user"`
	}
)

var (
	BorrowTimeRangeMap map[string]int = map[string]int{
		"oneweek":  7,
		"twoweek":  14,
		"onemonth": 30,
		"twomonth": 60,
	}

	borrowStatusMap = map[string]BorrowStatus{
		"init":     "INIT",
		"request":  "REQUEST",
		"progress": "PROGRESS",
		"returned": "RETURNED",
		"cancel":   "CANCEL",
	}
)

func GetBorrowStatus(s string) BorrowStatus {
	return borrowStatusMap[s]
}
