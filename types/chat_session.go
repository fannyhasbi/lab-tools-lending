package types

type ChatSessionStatusType string
type ChatSession struct {
	ID        int64                 `json:"id"`
	Status    ChatSessionStatusType `json:"status"`
	UserID    int64                 `json:"user_id"`
	CreatedAt string                `json:"created_at"`
	UpdatedAt string                `json:"updated_at"`
}

var ChatSessionStatus map[string]ChatSessionStatusType = map[string]ChatSessionStatusType{
	"progress": "PROGRESS",
	"complete": "COMPLETE",
}

type TopicType string
type ChatSessionDetail struct {
	ID            int64     `json:"id"`
	Topic         TopicType `json:"topic"`
	ChatSessionID int64     `json:"chat_session_id"`
	CreatedAt     string    `json:"created_at"`
}

var Topic map[string]TopicType = map[string]TopicType{
	"register_init":     "RGR_init",
	"register_confirm":  "RGR_confirm",
	"register_complete": "RGR_complete",
}
