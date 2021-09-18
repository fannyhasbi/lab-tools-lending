package types

type (
	UserType string

	User struct {
		ID        int64    `json:"id"`
		Name      string   `json:"name"`
		NIM       string   `json:"nim"`
		Batch     uint16   `json:"batch"`
		Address   string   `json:"address"`
		CreatedAt string   `json:"created_at"`
		UserType  UserType `json:"user_type"`
	}
)

const (
	UserTypeStudent UserType = "student"
	UserTypeAdmin   UserType = "admin"
	UserTypeBoth    UserType = "both"
)
