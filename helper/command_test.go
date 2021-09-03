package helper

import (
	"fmt"
	"testing"

	"github.com/fannyhasbi/lab-tools-lending/types"
	"github.com/stretchr/testify/assert"
)

func TestGetCommand(t *testing.T) {
	t.Run("return the help command", func(t *testing.T) {
		m := "/help"
		r := GetCommand(m)

		assert.Equal(t, "help", r)
	})

	t.Run("return empty string if no '/' symbol", func(t *testing.T) {
		m := "haztest"
		r := GetCommand(m)

		assert.Empty(t, r)
	})

	t.Run("return empty string if '/' symbol is not in the first char", func(t *testing.T) {
		m := "haztest/"
		r := GetCommand(m)

		assert.Empty(t, r)
	})

	t.Run("return the command, separated with space char", func(t *testing.T) {
		m := "/help haztest"
		r := GetCommand(m)

		assert.Equal(t, "help", r)
	})
}

func TestGetRespondCommands(t *testing.T) {
	t.Run("correct", func(t *testing.T) {
		s := fmt.Sprintf("/%s %s %d yes", types.CommandRespond, types.RespondTypeBorrow, 123)
		r, ok := GetRespondCommands(s)

		expected := types.RespondCommands{
			Type: types.RespondTypeBorrow,
			ID:   123,
			Text: "yes",
		}

		assert.True(t, ok)
		assert.Equal(t, expected, r)
	})

	t.Run("length less than 3", func(t *testing.T) {
		s := fmt.Sprintf("/%s %s", types.CommandRespond, types.RespondTypeBorrow)
		r, ok := GetRespondCommands(s)

		assert.False(t, ok)
		assert.Equal(t, types.RespondCommands{}, r)
	})

	t.Run("length equal 3", func(t *testing.T) {
		s := fmt.Sprintf("/%s %s %d", types.CommandRespond, types.RespondTypeBorrow, 123)
		r, ok := GetRespondCommands(s)

		expected := types.RespondCommands{
			Type: types.RespondTypeBorrow,
			ID:   123,
			Text: "",
		}

		assert.True(t, ok)
		assert.Equal(t, expected, r)
	})

	t.Run("exceed length 4", func(t *testing.T) {
		s := fmt.Sprintf("/%s %s %d yes oke nais 123", types.CommandRespond, types.RespondTypeBorrow, 123)
		r, ok := GetRespondCommands(s)

		assert.False(t, ok)
		assert.Equal(t, types.RespondCommands{}, r)
	})

	t.Run("not in category", func(t *testing.T) {
		s := fmt.Sprintf("/%s testnotincategory %d yes", types.CommandRespond, 123)
		r, ok := GetRespondCommands(s)

		assert.False(t, ok)
		assert.Equal(t, types.RespondCommands{}, r)
	})

	t.Run("wrong id", func(t *testing.T) {
		s := fmt.Sprintf("/%s %s testwrongid yes", types.CommandRespond, types.RespondTypeBorrow)
		r, ok := GetRespondCommands(s)

		assert.False(t, ok)
		assert.Equal(t, types.RespondCommands{}, r)
	})
}

func TestIsRespondTypeExists(t *testing.T) {
	t.Run("exist", func(t *testing.T) {
		r := isRespondTypeExists(types.RespondTypeBorrow)
		assert.True(t, r)
	})
	t.Run("nope", func(t *testing.T) {
		r := isRespondTypeExists("testdoesnotexists")
		assert.False(t, r)
	})
}

func TestIsCommandTypeExists(t *testing.T) {
	t.Run("exist", func(t *testing.T) {
		r := isManageTypeExists(types.ManageTypeAdd)
		assert.True(t, r)
	})
	t.Run("nope", func(t *testing.T) {
		r := isManageTypeExists("testdoesnotexists")
		assert.False(t, r)
	})
}

func TestGetManageCommands(t *testing.T) {
	t.Run("full value", func(t *testing.T) {
		s := fmt.Sprintf("/%s %s %d", types.CommandManage, types.ManageTypeEdit, 123)
		r, ok := GetManageCommands(s)

		expected := types.ManageCommands{
			Type: types.ManageTypeEdit,
			ID:   123,
		}

		assert.True(t, ok)
		assert.Equal(t, expected, r)
	})

	t.Run("length less than 2", func(t *testing.T) {
		s := fmt.Sprintf("/%s", types.CommandManage)
		r, ok := GetManageCommands(s)

		assert.False(t, ok)
		assert.Equal(t, types.ManageCommands{}, r)
	})

	t.Run("length equal 2 with correct type", func(t *testing.T) {
		s := fmt.Sprintf("/%s %s", types.CommandManage, types.ManageTypeAdd)
		r, ok := GetManageCommands(s)

		expected := types.ManageCommands{
			Type: types.ManageTypeAdd,
		}
		assert.True(t, ok)
		assert.Equal(t, expected, r)
	})

	t.Run("length equal 2 with incorrect type", func(t *testing.T) {
		s := fmt.Sprintf("/%s %s", types.CommandManage, types.ManageTypeEdit)
		r, ok := GetManageCommands(s)

		assert.False(t, ok)
		assert.Equal(t, types.ManageCommands{}, r)
	})

	t.Run("length exceed 3", func(t *testing.T) {
		s := fmt.Sprintf("/%s %s %d %s", types.CommandManage, types.ManageTypeEdit, 123, "testexceed")
		r, ok := GetManageCommands(s)

		assert.False(t, ok)
		assert.Equal(t, types.ManageCommands{}, r)
	})

	t.Run("wrong type", func(t *testing.T) {
		s := fmt.Sprintf("/%s testwrongtype %d", types.CommandManage, 123)
		r, ok := GetManageCommands(s)

		assert.False(t, ok)
		assert.Equal(t, types.ManageCommands{}, r)
	})

	t.Run("wrong id", func(t *testing.T) {
		s := fmt.Sprintf("/%s %s testwrongid", types.CommandManage, types.ManageTypeEdit)
		r, ok := GetManageCommands(s)

		assert.False(t, ok)
		assert.Equal(t, types.ManageCommands{}, r)
	})
}
