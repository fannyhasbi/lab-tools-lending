package service

import (
	"github.com/fannyhasbi/lab-tools-lending/config"
	"github.com/fannyhasbi/lab-tools-lending/repository"
	"github.com/fannyhasbi/lab-tools-lending/repository/postgres"
	"github.com/fannyhasbi/lab-tools-lending/types"
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

func (cs ChatSessionService) GetChatSessions(user types.User) ([]types.ChatSession, error) {
	result := cs.Query.Get(user)

	if result.Error != nil {
		return []types.ChatSession{}, result.Error
	}

	return result.Result.([]types.ChatSession), nil
}

func (cs ChatSessionService) GetChatSessionDetails(chatSession types.ChatSession) ([]types.ChatSessionDetail, error) {
	result := cs.Query.GetDetail(chatSession)

	if result.Error != nil {
		return []types.ChatSessionDetail{}, result.Error
	}

	return result.Result.([]types.ChatSessionDetail), nil
}
