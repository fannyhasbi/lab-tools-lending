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
		r := IsRespondTypeExists(types.RespondTypeBorrow)
		assert.True(t, r)
	})
	t.Run("nope", func(t *testing.T) {
		r := IsRespondTypeExists("testdoesnotexists")
		assert.False(t, r)
	})
}
