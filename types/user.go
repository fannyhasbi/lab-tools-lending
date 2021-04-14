package types

type User struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	NIM       string `json:"nim"`
	Batch     uint16 `json:"batch"`
	Address   string `json:"address"`
	CreatedAt string `json:"created_at"`
}

func (u *User) IsRegistered() bool {
	return u.ID > 0 && u.Name != "" && u.NIM != "" && u.CreatedAt != ""
}
