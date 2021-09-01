package types

type (
	ToolReturningStatus string
	ToolReturning       struct {
		ID             int64               `json:"id"`
		ReturnedAt     string              `json:"returned_at"`
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
		"complete": "COMPLETE",
	}
)

func GetToolReturningStatus(s string) ToolReturningStatus {
	return toolReturningStatusMap[s]
}
