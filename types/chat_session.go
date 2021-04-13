package types

import "time"

type ChatSession struct {
	ID        int64     `json:"id"`
	Status    string    `json:"status"`
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ChatSessionDetail struct {
	ID            int64     `json:"id"`
	Topic         string    `json:"topic"`
	ChatSessionID int64     `json:"chat_session_id"`
	CreatedAt     time.Time `json:"created_at"`
}

var ChatSessionStatus map[string]string = map[string]string{
	"progress": "PROGRESS",
	"complete": "COMPLETE",
}

var ChatSessionTopic map[string]string = map[string]string{
	"register": "REGISTER",
}
