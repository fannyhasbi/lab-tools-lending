package types

import "database/sql"

type (
	ToolReturningStatus string
	ToolReturning       struct {
		ID             int64               `json:"id"`
		CreatedAt      string              `json:"created_at"`
		ConfirmedAt    sql.NullTime        `json:"confirmed_at"`
		ConfirmedBy    sql.NullString      `json:"confirmed_by"`
		UserID         int64               `json:"user_id"`
		ToolID         int64               `json:"tool_id"`
		Status         ToolReturningStatus `json:"status"`
		AdditionalInfo string              `json:"additional_info"`
		Tool           Tool                `json:"tool"`
		User           User                `json:"user"`
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
