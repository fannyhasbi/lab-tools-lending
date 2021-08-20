package helper

import (
	"testing"

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
