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

func (cs ChatSessionService) GetChatSession(user types.User, requestType types.RequestType) (types.ChatSession, error) {
	result := cs.Query.Get(user, requestType)

	if result.Error != nil {
		return types.ChatSession{}, result.Error
	}

	return result.Result.(types.ChatSession), nil
}

func (cs ChatSessionService) GetChatSessionDetails(chatSession types.ChatSession) ([]types.ChatSessionDetail, error) {
	result := cs.Query.GetDetail(chatSession)

	if result.Error != nil {
		return []types.ChatSessionDetail{}, result.Error
	}

	return result.Result.([]types.ChatSessionDetail), nil
}

func (cs ChatSessionService) SaveChatSession(chatSession types.ChatSession, requestType types.RequestType) (types.ChatSession, error) {
	result, err := cs.Repository.Save(&chatSession, requestType)
	if err != nil {
		return types.ChatSession{}, err
	}

	return result, nil
}

func (cs ChatSessionService) UpdateChatSessionStatus(id int64, status types.ChatSessionStatusType) error {
	return cs.Repository.UpdateStatus(id, status)
}

func (cs ChatSessionService) DeleteChatSession(id int64) error {
	return cs.Repository.Delete(id)
}

func (cs ChatSessionService) SaveChatSessionDetail(chatSessionDetail types.ChatSessionDetail) (types.ChatSessionDetail, error) {
	result, err := cs.Repository.SaveDetail(&chatSessionDetail)
	if err != nil {
		return types.ChatSessionDetail{}, err
	}

	return result, nil
}

func (cs ChatSessionService) DeleteChatSessionDetailByChatSessionID(id int64) error {
	return cs.Repository.DeleteDetailByChatSessionID(id)
}
