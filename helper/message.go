package helper

import (
	"fmt"
	"strings"

	"github.com/fannyhasbi/lab-tools-lending/types"
)

func SplitNewLine(s string) []string {
	return strings.Split(s, "\n")
}

func RemoveTab(s string) string {
	return strings.Replace(s, "\t", "", -1)
}

func BuildMessageRequest(data *types.MessageRequest) {
	if len(data.ReplyMarkup.InlineKeyboard) == 0 {
		inlineKeyboard := make([][]types.InlineKeyboardButton, 0)
		data.ReplyMarkup.InlineKeyboard = inlineKeyboard
	}
}

func BuildToolListMessage(l []types.Tool) string {
	m := ""
	for _, t := range l {
		m = fmt.Sprintf("%s[%d] %s", m, t.ID, t.Name)
		if t.Stock < 1 {
			m += " (stok kosong)"
		}
		m += "\n"
	}
	return m
}
