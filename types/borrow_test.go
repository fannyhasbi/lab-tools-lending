package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBorrowStatus(t *testing.T) {
	t.Run("correct", func(t *testing.T) {
		r := GetBorrowStatus("progress")

		assert.Equal(t, borrowStatusMap["progress"], r)
	})

	t.Run("empty", func(t *testing.T) {
		r := GetBorrowStatus("testwrong")

		assert.Empty(t, r)
		assert.Equal(t, BorrowStatus(""), r)
	})
}
