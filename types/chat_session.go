package types

import "time"

type ChatSession struct {
	ID        int64     `json:"id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ChatSessionDetail struct {
	ID        int64     `json:"id"`
	Topic     int64     `json:"topic"`
	CreatedAt time.Time `json:"created_at"`
}
