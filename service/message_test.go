package service

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/fannyhasbi/lab-tools-lending/types"
	"github.com/stretchr/testify/assert"
)

func TestCanChangeChatSessionDetails(t *testing.T) {
	ms := &MessageService{}

	firstDetails := []types.ChatSessionDetail{
		{
			ID:    1,
			Topic: types.Topic["register_init"],
		},
		{
			ID:    2,
			Topic: types.Topic["register_confirm"],
		},
	}

	secondDetail := []types.ChatSessionDetail{
		{
			ID:    1,
			Topic: types.Topic["borrow_init"],
		},
		{
			ID:    2,
			Topic: types.Topic["borrow_confirm"],
		},
	}

	ms.chatSessionDetails = firstDetails
	ms.ChangeChatSessionDetails(secondDetail)

	assert.Equal(t, secondDetail, ms.chatSessionDetails)
}

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
		currentYear := time.Now().Year()
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

	t.Run("error in batch", func(t *testing.T) {
		r := types.QuestionRegistration{
			Name:    testname,
			NIM:     testnim,
			Batch:   1234,
			Address: testaddress,
		}

		err := validateRegisterConfirmation(r)
		assert.Error(t, err)
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

func TestIsIDWithinCommand(t *testing.T) {
	t.Run("return true and the int64", func(t *testing.T) {
		m := "/pinjam 321"
		id, ok := isIDWithinCommand(m)

		assert.True(t, ok)
		assert.Equal(t, int64(321), id)
	})

	t.Run("not affected by the command type", func(t *testing.T) {
		m := "/yoyoyoyoyoy 321"
		id, ok := isIDWithinCommand(m)

		assert.True(t, ok)
		assert.Equal(t, int64(321), id)
	})

	t.Run("exceed split length", func(t *testing.T) {
		m := "/pinjam 321 1 1 1 1"
		id, ok := isIDWithinCommand(m)

		assert.False(t, ok)
		assert.Equal(t, int64(0), id)
	})

	t.Run("invalid id format", func(t *testing.T) {
		m := "/pinjam hello"
		id, ok := isIDWithinCommand(m)

		assert.False(t, ok)
		assert.Equal(t, int64(0), id)
	})
}

func TestIsFlagWithinReturningCommand(t *testing.T) {
	t.Run("basic true", func(t *testing.T) {
		m := fmt.Sprintf("/%s %s", types.CommandReturn, types.ToolReturningFlag)
		ok := isFlagWithinReturningCommand(m)

		assert.True(t, ok)
	})

	t.Run("correct command wrong flag", func(t *testing.T) {
		m := fmt.Sprintf("/%s testwrong", types.CommandReturn)
		ok := isFlagWithinReturningCommand(m)

		assert.False(t, ok)
	})

	t.Run("correct flag but exceed split", func(t *testing.T) {
		m := fmt.Sprintf("/%s %s test correct exceed", types.CommandReturn, types.ToolReturningFlag)
		ok := isFlagWithinReturningCommand(m)

		assert.False(t, ok)
	})
}
