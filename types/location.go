package types

type LocationRequest struct {
	ChatID    int64   `json:"chat_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Heading   uint8   `json:"heading"`
}
