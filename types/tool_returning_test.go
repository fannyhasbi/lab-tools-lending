package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetToolReturningStatus(t *testing.T) {
	t.Run("correct", func(t *testing.T) {
		r := GetToolReturningStatus("request")

		assert.Equal(t, toolReturningStatusMap["request"], r)
	})

	t.Run("empty", func(t *testing.T) {
		r := GetToolReturningStatus("testwrong")

		assert.Empty(t, r)
		assert.Equal(t, ToolReturningStatus(""), r)
	})
}
