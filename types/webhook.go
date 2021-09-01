package types

type (
	RequestType string

	WebhookRequest struct {
		Message struct {
			From struct {
				ID int64 `json:"id"`
			}
			Text string `json:"text"`
			Chat struct {
				ID   int64  `json:"id"`
				Type string `json:"type"`
			} `json:"chat"`
		} `json:"message"`
	}

	InlineCallbackQuery struct {
		CallbackQuery struct {
			ID   string `json:"id"`
			From struct {
				ID int64 `json:"id"`
			} `json:"from"`
			Message struct {
				Chat struct {
					ID   int64  `json:"id"`
					Type string `json:"type"`
				} `json:"chat"`
			} `json:"message"`
			Data string `json:"data"`
		} `json:"callback_query"`
	}
)

var (
	RequestTypePrivate RequestType = "private"
	RequestTypeGroup   RequestType = "group"
)
