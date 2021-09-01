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

func (us UserService) SaveUser(user types.User) (types.User, error) {
	result, err := us.Repository.Save(&user)
	if err != nil {
		return types.User{}, err
	}

	return result, nil
}

func (us UserService) UpdateUser(user types.User) (types.User, error) {
	result, err := us.Repository.Update(&user)
	if err != nil {
		return types.User{}, err
	}

	return result, nil
}

func (us UserService) DeleteUser(id int64) error {
	return us.Repository.Delete(id)
}

func (us UserService) FindByID(id int64) (types.User, error) {
	result := us.Query.FindByID(id)
	if result.Error == sql.ErrNoRows {
		return types.User{ID: id}, result.Error
	}

	if result.Error != nil {
		return types.User{ID: id}, result.Error
	}

	return result.Result.(types.User), nil
}
