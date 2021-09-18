package repository

import "github.com/fannyhasbi/lab-tools-lending/types"

type UserQuery interface {
	FindByID(chatID int64) QueryResult
}

type UserRepository interface {
	Save(user *types.User) (types.User, error)
	Update(user *types.User) (types.User, error)
	Delete(id int64) error
	UpdateUserType(id int64, userType types.UserType) error
}
