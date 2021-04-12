package service

import (
	"github.com/fannyhasbi/lab-tools-lending/config"
	"github.com/fannyhasbi/lab-tools-lending/repository"
	"github.com/fannyhasbi/lab-tools-lending/repository/postgres"
)

type ChatSessionService struct {
	Query      repository.ChatSessionQuery
	Repository repository.ChatSessionRepository
}

func NewChatSessionService() *ChatSessionService {
	var chatSessionQuery repository.ChatSessionQuery
	var chatSessionRepository repository.ChatSessionRepository

	db := config.InitPostgresDB()
	chatSessionQuery = postgres.NewChatSessionQueryPostgres(db)
	chatSessionRepository = postgres.NewChatSessionRepositoryPostgres(db)

	return &ChatSessionService{
		Query:      chatSessionQuery,
		Repository: chatSessionRepository,
	}
}
