package types

type (
	RequestType string

	TeleMessageFrom struct {
		ID        int64  `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Username  string `json:"username"`
	}

	teleMessageChat struct {
		ID   int64  `json:"id"`
		Type string `json:"type"`
	}

	TelePhotoSize struct {
		FileID       string `json:"file_id"`
		FileUniqueID string `json:"file_unique_id"`
		FileSize     int64  `json:"file_size"`
		Width        int    `json:"width"`
		Height       int    `json:"height"`
	}

	TeleMessage struct {
		MessageID    int64           `json:"message_id"`
		From         TeleMessageFrom `json:"from"`
		Text         string          `json:"text"`
		Chat         teleMessageChat `json:"chat"`
		MediaGroupID string          `json:"media_group_id"`
		Photo        []TelePhotoSize `json:"photo"`
	}

	WebhookRequest struct {
		Message TeleMessage `json:"message"`
	}

	teleCallbackQuery struct {
		ID      string          `json:"id"`
		From    TeleMessageFrom `json:"from"`
		Message TeleMessage     `json:"message"`
		Data    string          `json:"data"`
	}

	InlineCallbackQuery struct {
		CallbackQuery teleCallbackQuery `json:"callback_query"`
	}
)

var (
	RequestTypePrivate RequestType = "private"
	RequestTypeGroup   RequestType = "group"
)
