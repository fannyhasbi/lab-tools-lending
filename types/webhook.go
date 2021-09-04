package types

type (
	RequestType string

	teleMessageFrom struct {
		ID int64 `json:"id"`
	}

	teleMessageChat struct {
		ID   int64  `json:"id"`
		Type string `json:"type"`
	}

	teleMessage struct {
		MessageID int64           `json:"message_id"`
		From      teleMessageFrom `json:"from"`
		Text      string          `json:"text"`
		Chat      teleMessageChat `json:"chat"`
	}

	WebhookRequest struct {
		Message teleMessage `json:"message"`
	}

	teleCallbackQuery struct {
		ID      string          `json:"id"`
		From    teleMessageFrom `json:"from"`
		Message teleMessage     `json:"message"`
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
