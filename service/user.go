package service

import (
	"database/sql"

	"github.com/fannyhasbi/lab-tools-lending/config"
	"github.com/fannyhasbi/lab-tools-lending/repository"
	"github.com/fannyhasbi/lab-tools-lending/repository/postgres"
	"github.com/fannyhasbi/lab-tools-lending/types"
)

type UserService struct {
	Query      repository.UserQuery
	Repository repository.UserRepository
}

func NewUserService() *UserService {
	var userQuery repository.UserQuery
	var userRepository repository.UserRepository

	db := config.InitPostgresDB()
	userQuery = postgres.NewUserQueryPostgres(db)
	userRepository = postgres.NewUserRepositoryPostgres(db)

	return &UserService{
		Query:      userQuery,
		Repository: userRepository,
	}
}

func (us UserService) FindByChatID(chatID int64) (types.User, error) {
	result := us.Query.FindByChatID(chatID)

	if result.Error == sql.ErrNoRows {
		return types.User{ChatID: chatID}, result.Error
	}

	if result.Error != nil {
		return types.User{}, result.Error
	}

	return result.Result.(types.User), nil
}
