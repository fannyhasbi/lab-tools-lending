package helper

import (
	"fmt"
	"strconv"
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

func GetReportTimeFromCommand(yearmonth string) (int, int, bool) {
	splittedTime := strings.Split(yearmonth, "-")
	if len(splittedTime) < 2 {
		return 0, 0, false
	}

	year, err := strconv.Atoi(splittedTime[0])
	if err != nil {
		return 0, 0, false
	}

	month, err := strconv.Atoi(splittedTime[1])
	if err != nil {
		return 0, 0, false
	}

	if month < 1 || month > 12 {
		return 0, 0, false
	}

	return year, month, true
}
