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
