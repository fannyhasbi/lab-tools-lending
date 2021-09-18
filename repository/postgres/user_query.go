package postgres

import (
	"database/sql"

	"github.com/fannyhasbi/lab-tools-lending/repository"
	"github.com/fannyhasbi/lab-tools-lending/types"
)

type UserQueryPostgres struct {
	DB *sql.DB
}

func NewUserQueryPostgres(DB *sql.DB) repository.UserQuery {
	return &UserQueryPostgres{
		DB: DB,
	}
}

func (uq UserQueryPostgres) FindByID(chatID int64) repository.QueryResult {
	row := uq.DB.QueryRow(`
		SELECT id, name, nim, batch, address, created_at, user_type
		FROM users
		WHERE id = $1
	`, chatID)

	user := types.User{}
	result := repository.QueryResult{}

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.NIM,
		&user.Batch,
		&user.Address,
		&user.CreatedAt,
		&user.UserType,
	)

	if err != nil {
		result.Error = err
		return result
	}

	result.Result = user
	return result
}
