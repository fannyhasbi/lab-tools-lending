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
	row := csq.DB.QueryRow(`
		SELECT id, status, created_at
		FROM chat_sessions
		WHERE user_id = $1
			AND status = $2
		ORDER BY id DESC
	`, user.ID, types.ChatSessionStatus["progress"])

	chatSession := types.ChatSession{}
	result := repository.QueryResult{}

	err := row.Scan(
		&chatSession.ID,
		&chatSession.Status,
		&chatSession.CreatedAt,
	)

	if err != nil {
		result.Error = err
		return result
	}

	result.Result = chatSession
	return result
}

func (csq ChatSessionQueryPostgres) GetDetail(chatSession types.ChatSession) repository.QueryResult {
	rows, err := csq.DB.Query(`
		SELECT id, topic, chat_session_id, created_at, data
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
				&temp.Data,
			)

			chatSessionDetails = append(chatSessionDetails, temp)
		}
		result.Result = chatSessionDetails
	}
	return result
}
