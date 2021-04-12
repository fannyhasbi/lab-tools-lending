package repository

import "github.com/fannyhasbi/lab-tools-lending/types"

type ChatSessionQuery interface {
	Get(user types.User) QueryResult
	GetDetail(chatSession types.ChatSession) QueryResult
}

type ChatSessionRepository interface {
	Save(chatSession *types.ChatSession) error
	Update(chatSession *types.ChatSession) error
}
