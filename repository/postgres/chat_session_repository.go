package postgres

import (
	"database/sql"

	"github.com/fannyhasbi/lab-tools-lending/repository"
	"github.com/fannyhasbi/lab-tools-lending/types"
)

type ChatSessionRepositoryPostgres struct {
	DB *sql.DB
}

func NewChatSessionRepositoryPostgres(DB *sql.DB) repository.ChatSessionRepository {
	return &ChatSessionRepositoryPostgres{
		DB: DB,
	}
}

func (csr *ChatSessionRepositoryPostgres) Save(chatSession *types.ChatSession) error {
	return nil
}
func (csr *ChatSessionRepositoryPostgres) Update(chatSession *types.ChatSession) error {
	return nil
}
