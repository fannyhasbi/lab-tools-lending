package repository

import "github.com/fannyhasbi/lab-tools-lending/types"

type ChatSessionQuery interface {
	Get(user types.User) QueryResult
	GetDetail(chatSession types.ChatSession) QueryResult
}

type ChatSessionRepository interface {
	Save(chatSession *types.ChatSession) (types.ChatSession, error)
	UpdateStatus(id int64, status types.ChatSessionStatusType) error
	Delete(id int64) error
	SaveDetail(chatSessionDetail *types.ChatSessionDetail) (types.ChatSessionDetail, error)
	DeleteDetailByChatSessionID(id int64) error
}
