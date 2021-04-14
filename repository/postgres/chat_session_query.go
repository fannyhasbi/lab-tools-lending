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
	rows, err := csq.DB.Query(`
		SELECT id, status, created_at
		FROM chat_sessions
		WHERE user_id = $1
			AND status = $2
		ORDER BY id DESC
	`, user.ID, types.ChatSessionStatus["progress"])

	chatSessions := []types.ChatSession{}
	result := repository.QueryResult{}

	if err != nil {
		result.Error = err
	} else {
		for rows.Next() {
			temp := types.ChatSession{}
			rows.Scan(
				&temp.ID,
				&temp.Status,
				&temp.CreatedAt,
			)

			chatSessions = append(chatSessions, temp)
		}
		result.Result = chatSessions
	}
	return result
}

func (csq ChatSessionQueryPostgres) GetDetail(chatSession types.ChatSession) repository.QueryResult {
	rows, err := csq.DB.Query(`
		SELECT id, topic, chat_session_id, created_at
		FROM chat_session_details
		WHERE chat_session_id = $1
		ORDER BY id DESC
	`, chatSession.ID)

	chatSessionDetails := []types.ChatSessionDetail{}
	result := repository.QueryResult{}

	if err != nil {
		result.Error = err
	} else {
		for rows.Next() {
			temp := types.ChatSessionDetail{}
			rows.Scan(
				&temp.ID,
				&temp.Topic,
				&temp.ChatSessionID,
				&temp.CreatedAt,
			)

			chatSessionDetails = append(chatSessionDetails, temp)
		}
		result.Result = chatSessionDetails
	}
	return result
}
