package types

import "database/sql"

type (
	BorrowStatus    string
	borrowTimeRange struct {
		OneWeek  string
		TwoWeek  string
		OneMonth string
		TwoMonth string
	}

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

func GetBorrowTimeRange() borrowTimeRange {
	return borrowTimeRange{
		OneWeek:  "1WEEK",
		TwoWeek:  "2WEEK",
		OneMonth: "1MONTH",
		TwoMonth: "2MONTH",
	}
}
