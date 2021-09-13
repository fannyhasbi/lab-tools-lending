package helper

import (
	"fmt"
	"testing"

	"github.com/fannyhasbi/lab-tools-lending/types"
	"github.com/stretchr/testify/assert"
)

func TestSplitNewLine(t *testing.T) {
	s := "hello\nworld\ntest"
	r := SplitNewLine(s)

	assert.NotZero(t, len(r))
	assert.Equal(t, []string{"hello", "world", "test"}, r)
}

func TestRemoveTab(t *testing.T) {
	s := `hello
	world
	test`
	r := RemoveTab(s)

	assert.Equal(t, "hello\nworld\ntest", r)
}

func TestCanBuildMessageRequest(t *testing.T) {
	var id int64 = 123
	text := "hello"

	t.Run("InlineKeyboard is not nil", func(t *testing.T) {
		req := types.MessageRequest{
			ChatID: id,
			Text:   text,
		}

		assert.Nil(t, req.ReplyMarkup.InlineKeyboard)
		assert.Equal(t, 0, len(req.ReplyMarkup.InlineKeyboard))
		assert.Equal(t, 0, cap(req.ReplyMarkup.InlineKeyboard))
		assert.Empty(t, req.ParseMode)

		BuildMessageRequest(&req)

		assert.NotNil(t, req.ReplyMarkup.InlineKeyboard)
		assert.Equal(t, 0, len(req.ReplyMarkup.InlineKeyboard))
		assert.Equal(t, 0, cap(req.ReplyMarkup.InlineKeyboard))
		assert.Equal(t, id, req.ChatID)
		assert.Equal(t, text, req.Text)
		assert.Empty(t, req.ParseMode)
	})

	t.Run("prefilled InlineKeyboard", func(t *testing.T) {
		ikb := [][]types.InlineKeyboardButton{
			{
				{
					Text:         "yes",
					CallbackData: "yes",
				},
				{
					Text:         "no",
					CallbackData: "no",
				},
			},
		}

		replyMarkup := types.InlineKeyboardMarkup{
			InlineKeyboard: ikb,
		}

		req := types.MessageRequest{
			ChatID:      id,
			Text:        text,
			ReplyMarkup: replyMarkup,
		}

		assert.NotNil(t, req.ReplyMarkup.InlineKeyboard)

		BuildMessageRequest(&req)

		assert.NotNil(t, req.ReplyMarkup.InlineKeyboard)
		assert.Equal(t, len(ikb), len(req.ReplyMarkup.InlineKeyboard))
		assert.Equal(t, len(ikb[0]), len(req.ReplyMarkup.InlineKeyboard[0]))
		assert.Equal(t, replyMarkup, req.ReplyMarkup)
	})
}

func TestCanBuildToolListMessage(t *testing.T) {
	tools := []types.Tool{
		{
			ID:    123,
			Name:  "hello1",
			Stock: 10,
		},
		{
			ID:    321,
			Name:  "hello2",
			Stock: 0,
		},
	}

	r := BuildToolListMessage(tools)

	expected := fmt.Sprintf("[%d] %s\n[%d] %s (stok kosong)\n", tools[0].ID, tools[0].Name, tools[1].ID, tools[1].Name)

	assert.Equal(t, expected, r)
}

func TestGetReportTimeFromCommand(t *testing.T) {
	t.Run("correct", func(t *testing.T) {
		s := "2021-8"
		y, m, ok := GetReportTimeFromCommand(s)

		assert.True(t, ok)
		assert.Equal(t, 2021, y)
		assert.Equal(t, 8, m)
	})
	t.Run("zero filled", func(t *testing.T) {
		s := "00002021-00008"
		y, m, ok := GetReportTimeFromCommand(s)

		assert.True(t, ok)
		assert.Equal(t, 2021, y)
		assert.Equal(t, 8, m)
	})
	t.Run("length exceed 2", func(t *testing.T) {
		s := "2021-08-30"
		y, m, ok := GetReportTimeFromCommand(s)

		assert.True(t, ok)
		assert.Equal(t, 2021, y)
		assert.Equal(t, 8, m)
	})
	t.Run("just year", func(t *testing.T) {
		s := "2021"
		y, m, ok := GetReportTimeFromCommand(s)

		assert.False(t, ok)
		assert.Zero(t, y)
		assert.Zero(t, m)
	})
	t.Run("random words", func(t *testing.T) {
		s := "hello-world"
		y, m, ok := GetReportTimeFromCommand(s)

		assert.False(t, ok)
		assert.Zero(t, y)
		assert.Zero(t, m)
	})
}
