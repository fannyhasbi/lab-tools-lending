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

func (uq UserQueryPostgres) FindByChatID(chatID int64) repository.QueryResult {
	row := uq.DB.QueryRow(`
		SELECT id, chat_id, name, nim, batch, address, created_at
		FROM users
		WHERE chat_id = $1
	`, chatID)

	user := types.User{}
	result := repository.QueryResult{}

	err := row.Scan(
		&user.ID,
		&user.ChatID,
		&user.Name,
		&user.NIM,
		&user.Batch,
		&user.Address,
		&user.CreatedAt,
	)
	if err != nil {
		result.Error = err
		return result
	}

	result.Result = user
	return result
}
