package types

type (
	MessageRequest struct {
		ChatID      int64                `json:"chat_id"`
		Text        string               `json:"text"`
		ParseMode   string               `json:"parse_mode"`
		ReplyMarkup InlineKeyboardMarkup `json:"reply_markup"`
	}

	InlineKeyboardMarkup struct {
		InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
	}

	InlineKeyboardButton struct {
		Text         string `json:"text"`
		CallbackData string `json:"callback_data"`
	}
)

var (
	// testing group
	AdminGroupIDs = []int64{-502157840}
)
