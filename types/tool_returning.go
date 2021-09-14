package types

import "database/sql"

type (
	ToolReturningStatus string
	ToolReturning       struct {
		ID             int64               `json:"id"`
		CreatedAt      string              `json:"created_at"`
		ConfirmedAt    sql.NullTime        `json:"confirmed_at"`
		ConfirmedBy    sql.NullString      `json:"confirmed_by"`
		BorrowID       int64               `json:"borrow_id"`
		Status         ToolReturningStatus `json:"status"`
		AdditionalInfo string              `json:"additional_info"`
		Borrow         Borrow              `json:"borrow"`
	}
)

var (
	ToolReturningFlag string = "1"

	toolReturningStatusMap = map[string]ToolReturningStatus{
		"request":  "REQUEST",
		"reject":   "REJECT",
		"complete": "COMPLETE",
	}
)

func GetToolReturningStatus(s string) ToolReturningStatus {
	return toolReturningStatusMap[s]
}
