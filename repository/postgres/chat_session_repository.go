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

func (csr *ChatSessionRepositoryPostgres) Save(chatSession *types.ChatSession) (types.ChatSession, error) {
	row := csr.DB.QueryRow(`INSERT INTO chat_sessions (status, user_id) VALUES ($1, $2)
		RETURNING id, status, user_id, created_at, updated_at`,
		chatSession.Status,
		chatSession.UserID)

	cs := types.ChatSession{}
	err := row.Scan(
		&cs.ID,
		&cs.Status,
		&cs.UserID,
		&cs.CreatedAt,
		&cs.UpdatedAt,
	)
	if err != nil {
		return types.ChatSession{}, err
	}

	return cs, nil
}

func (csr *ChatSessionRepositoryPostgres) UpdateStatus(id int64, status types.ChatSessionStatusType) error {
	_, err := csr.DB.Exec(`UPDATE chat_sessions SET status = $1 WHERE id = $2`, status, id)
	return err
}

func (csr *ChatSessionRepositoryPostgres) Delete(id int64) error {
	_, err := csr.DB.Exec(`DELETE FROM chat_sessions WHERE id = $1`, id)
	return err
}

func (csr *ChatSessionRepositoryPostgres) SaveDetail(chatSessionDetail *types.ChatSessionDetail) (types.ChatSessionDetail, error) {
	row := csr.DB.QueryRow(`INSERT INTO chat_session_details (topic, chat_session_id, data) VALUES ($1, $2, $3)
		RETURNING id, topic, chat_session_id, created_at, data`,
		chatSessionDetail.Topic,
		chatSessionDetail.ChatSessionID,
		chatSessionDetail.Data)

	csd := types.ChatSessionDetail{}
	err := row.Scan(
		&csd.ID,
		&csd.Topic,
		&csd.ChatSessionID,
		&csd.CreatedAt,
		&csd.Data,
	)
	if err != nil {
		return types.ChatSessionDetail{}, err
	}

	return csd, nil
}

func (csr *ChatSessionRepositoryPostgres) DeleteDetailByChatSessionID(id int64) error {
	_, err := csr.DB.Exec(`DELETE FROM chat_session_details WHERE chat_session_id = $1`, id)
	return err
}
