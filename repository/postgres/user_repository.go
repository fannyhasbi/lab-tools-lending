package postgres

import (
	"database/sql"

	"github.com/fannyhasbi/lab-tools-lending/repository"
	"github.com/fannyhasbi/lab-tools-lending/types"
)

type UserRepositoryPostgres struct {
	DB *sql.DB
}

func NewUserRepositoryPostgres(DB *sql.DB) repository.UserRepository {
	return &UserRepositoryPostgres{
		DB: DB,
	}
}

func (ur *UserRepositoryPostgres) Update(user *types.User) error {
	return nil
}
