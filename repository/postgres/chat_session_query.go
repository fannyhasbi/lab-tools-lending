package postgres

import (
	"database/sql"

	"github.com/fannyhasbi/lab-tools-lending/repository"
	"github.com/fannyhasbi/lab-tools-lending/types"
)

type ChatSessionQueryPostgres struct {
	DB *sql.DB
}

func NewChatSessionQueryPostgres(DB *sql.DB) repository.ChatSessionQuery {
	return &ChatSessionQueryPostgres{
		DB: DB,
	}
}

func (csq ChatSessionQueryPostgres) Get(user types.User) repository.QueryResult {
	return repository.QueryResult{}
}
func (csq ChatSessionQueryPostgres) GetDetail(chatSession types.ChatSession) repository.QueryResult {
	return repository.QueryResult{}
}
