package types

type (
	ChatSessionStatusType string
	TopicType             string

	ChatSession struct {
		ID        int64                 `json:"id"`
		Status    ChatSessionStatusType `json:"status"`
		UserID    int64                 `json:"user_id"`
		CreatedAt string                `json:"created_at"`
		UpdatedAt string                `json:"updated_at"`
	}

	ChatSessionDetail struct {
		ID            int64     `json:"id"`
		Topic         TopicType `json:"topic"`
		ChatSessionID int64     `json:"chat_session_id"`
		Data          string    `json:"data"`
		CreatedAt     string    `json:"created_at"`
	}
)

var (
	ChatSessionStatus map[string]ChatSessionStatusType = map[string]ChatSessionStatusType{
		"progress": "PROGRESS",
		"complete": "COMPLETE",
	}

	Topic map[string]TopicType = map[string]TopicType{
		"register_init":     "RGR_init",
		"register_confirm":  "RGR_confirm",
		"register_complete": "RGR_complete",

		"borrow_init":    "BRW_init",
		"borrow_date":    "BRW_date",
		"borrow_confirm": "BRW_confirm",

		"tool_returning_init":     "RET_init",
		"tool_returning_confirm":  "RET_confim",
		"tool_returning_complete": "RET_complete",
	}
)
