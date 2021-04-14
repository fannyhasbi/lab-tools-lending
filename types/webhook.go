package types

const (
	RequestTypeCommon         = "common"
	RequestTypeInlineCallback = "inline_callback_query"
)

type WebhookRequest struct {
	Message struct {
		Text string `json:"text"`
		Chat struct {
			ID int64 `json:"id"`
		} `json:"chat"`
	} `json:"message"`
}

type InlineCallbackQuery struct {
	CallbackQuery struct {
		ID   string `json:"id"`
		From struct {
			ID int64 `json:"id"`
		} `json:"from"`
		Data string `json:"data"`
	} `json:"callback_query"`
}
