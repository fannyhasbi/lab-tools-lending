package service

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/fannyhasbi/lab-tools-lending/types"
	"github.com/stretchr/testify/assert"
)

func TestGetRegistrationMessage(t *testing.T) {
	n := "testname"
	nim := "211XXXXXXXXXXXX"
	b := 2016
	a := "jalan test message"

	t.Run("can get message", func(t *testing.T) {
		message := fmt.Sprintf(`%s
			%s
			%d
			%s`, n, nim, b, a)
		message = strings.Replace(message, "\t", "", -1)

		r, err := getRegistrationMessage(message)

		expect := types.QuestionRegistration{
			Name:    n,
			NIM:     nim,
			Batch:   b,
			Address: a,
		}

		assert.NoError(t, err)
		assert.Equal(t, expect, r)
	})

	t.Run("error invalid format", func(t *testing.T) {
		message := `hello`
		r, err := getRegistrationMessage(message)
		expect := types.QuestionRegistration{}

		assert.Error(t, err)
		assert.Equal(t, expect, r)
	})

	t.Run("error batch int conversion", func(t *testing.T) {
		message := fmt.Sprintf(`%s
			%s
			thisiswrongintformat
			%s`, n, nim, a)
		message = strings.Replace(message, "\t", "", -1)

		r, err := getRegistrationMessage(message)
		expect := types.QuestionRegistration{}

		assert.Error(t, err)
		assert.Equal(t, expect, r)
	})
}

func TestValidateRegisterMessageBatch(t *testing.T) {
	t.Run("valid and no error", func(t *testing.T) {
		b := 2010
		err := validateRegisterMessageBatch(b)

		assert.NoError(t, err)
	})

	t.Run("error below the limit", func(t *testing.T) {
		b := 1111
		err := validateRegisterMessageBatch(b)

		assert.Error(t, err)
	})

	t.Run("error beyond the limit", func(t *testing.T) {
		currentYear, _, _ := time.Now().Date()
		b := currentYear + 100
		err := validateRegisterMessageBatch(b)

		assert.Error(t, err)
	})
}

func TestValidateRegisterConfirmation(t *testing.T) {
	testname := "testname"
	testnim := "2112xxxxxxxxxx"
	testbatch := 2016
	testaddress := "jalan test message"

	t.Run("success", func(t *testing.T) {
		r := types.QuestionRegistration{
			Name:    testname,
			NIM:     testnim,
			Batch:   testbatch,
			Address: testaddress,
		}

		err := validateRegisterConfirmation(r)
		assert.NoError(t, err)
	})

	t.Run("invalid name length", func(t *testing.T) {
		r := types.QuestionRegistration{
			Name:    "abc",
			NIM:     testnim,
			Batch:   testbatch,
			Address: testaddress,
		}

		err := validateRegisterConfirmation(r)
		assert.Error(t, err)
	})

	t.Run("invalid NIM length", func(t *testing.T) {
		r := types.QuestionRegistration{
			Name:    testname,
			NIM:     "123",
			Batch:   testbatch,
			Address: testaddress,
		}
		r2 := types.QuestionRegistration{
			Name:    testname,
			NIM:     "2112XXXXXXXXXXXXXXXXXXXXXX",
			Batch:   testbatch,
			Address: testaddress,
		}

		err := validateRegisterConfirmation(r)
		err2 := validateRegisterConfirmation(r2)
		assert.Error(t, err)
		assert.Error(t, err2)
	})

	t.Run("invalid address length", func(t *testing.T) {
		r := types.QuestionRegistration{
			Name:    testname,
			NIM:     testnim,
			Batch:   testbatch,
			Address: "jl",
		}

		err := validateRegisterConfirmation(r)
		assert.Error(t, err)
	})
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

		buildMessageRequest(&req)

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

		buildMessageRequest(&req)

		assert.NotNil(t, req.ReplyMarkup.InlineKeyboard)
		assert.Equal(t, len(ikb), len(req.ReplyMarkup.InlineKeyboard))
		assert.Equal(t, len(ikb[0]), len(req.ReplyMarkup.InlineKeyboard[0]))
		assert.Equal(t, replyMarkup, req.ReplyMarkup)
	})
}

func TestIsToolIDWithinBorrowMessage(t *testing.T) {
	t.Run("return true and the int64", func(t *testing.T) {
		m := "/pinjam 321"
		ok, id := isToolIDWithinBorrowMessage(m)

		assert.True(t, ok)
		assert.Equal(t, int64(321), id)
	})

	t.Run("not affected by the command type", func(t *testing.T) {
		m := "/yoyoyoyoyoy 321"
		ok, id := isToolIDWithinBorrowMessage(m)

		assert.True(t, ok)
		assert.Equal(t, int64(321), id)
	})

	t.Run("exceed split length", func(t *testing.T) {
		m := "/pinjam 321 1 1 1 1"
		ok, id := isToolIDWithinBorrowMessage(m)

		assert.False(t, ok)
		assert.Equal(t, int64(0), id)
	})

	t.Run("invalid id format", func(t *testing.T) {
		m := "/pinjam hello"
		ok, id := isToolIDWithinBorrowMessage(m)

		assert.False(t, ok)
		assert.Equal(t, int64(0), id)
	})
}
