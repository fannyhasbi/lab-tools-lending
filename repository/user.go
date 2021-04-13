package repository

import "github.com/fannyhasbi/lab-tools-lending/types"

type UserQuery interface {
	FindByChatID(chatID int64) QueryResult
}

type UserRepository interface {
	Update(user *types.User) error
}
