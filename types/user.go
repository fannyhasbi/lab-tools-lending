package types

type User struct {
	ID        int64  `json:"id"`
	ChatID    int64  `json:"chat_id"`
	Name      string `json:"name"`
	NIM       string `json:"nim"`
	Batch     uint16 `json:"batch"`
	Address   string `json:"address"`
	CreatedAt string `json:"created_at"`
}
