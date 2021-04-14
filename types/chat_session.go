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

type ChatSessionTopicType string
type ChatSessionDetail struct {
	ID            int64                `json:"id"`
	Topic         ChatSessionTopicType `json:"topic"`
	ChatSessionID int64                `json:"chat_session_id"`
	CreatedAt     string               `json:"created_at"`
}

var ChatSessionTopic map[string]ChatSessionTopicType = map[string]ChatSessionTopicType{
	"register_init":     "RGR_init",
	"register_confirm":  "RGR_confirm",
	"register_complete": "RGR_complete",
}
