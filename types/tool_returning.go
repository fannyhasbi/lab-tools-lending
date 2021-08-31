package types

type ToolReturning struct {
	ID             int64  `json:"id"`
	ReturnedAt     string `json:"returned_at"`
	UserID         int64  `json:"user_id"`
	ToolID         int64  `json:"tool_id"`
	AdditionalInfo string `json:"additional_info"`
}

var ToolReturningFlag string = "1"
