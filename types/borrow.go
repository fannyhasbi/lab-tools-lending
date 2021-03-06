package types

import "database/sql"

type (
	BorrowStatus string

	Borrow struct {
		ID          int64          `json:"id"`
		Amount      int            `json:"amount"`
		Duration    int            `json:"duration"`
		Status      BorrowStatus   `json:"status"`
		UserID      int64          `json:"user_id"`
		ToolID      int64          `json:"tool_id"`
		CreatedAt   string         `json:"created_at"`
		ConfirmedAt sql.NullTime   `json:"confirmed_at"`
		ConfirmedBy sql.NullString `json:"confirmed_by"`
		Reason      sql.NullString `json:"reason"`
		Tool        Tool           `json:"tool"`
		User        User           `json:"user"`
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
		"request":  "REQUEST",
		"reject":   "REJECT",
		"progress": "PROGRESS",
		"returned": "RETURNED",
	}

	BorrowMinimalDuration = 7
)

func GetBorrowStatus(s string) BorrowStatus {
	return borrowStatusMap[s]
}
