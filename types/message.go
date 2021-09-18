package types

type (
	MessageRequest struct {
		ChatID      int64                `json:"chat_id"`
		Text        string               `json:"text"`
		ParseMode   string               `json:"parse_mode"`
		ReplyMarkup InlineKeyboardMarkup `json:"reply_markup"`
	}

	PhotoRequest struct {
		ChatID int64  `json:"chat_id"`
		Photo  string `json:"photo"`
	}

	PhotoGroupRequest struct {
		ChatID int64             `json:"chat_id"`
		Media  []InputMediaPhoto `json:"media"`
	}

	InputMediaPhoto struct {
		Type  string `json:"type"`
		Media string `json:"media"`
	}

	InlineKeyboardMarkup struct {
		InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
	}

	InlineKeyboardButton struct {
		Text         string `json:"text"`
		CallbackData string `json:"callback_data"`
	}
)
