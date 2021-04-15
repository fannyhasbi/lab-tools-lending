package types

type Borrow struct {
	ID         int64  `json:"id"`
	Amount     int    `json:"amount"`
	ReturnDate string `json:"return_date"`
	UserID     int64  `json:"user_id"`
	ToolID     int64  `json:"tool_id"`
	CreatedAt  string `json:"created_at"`
}
