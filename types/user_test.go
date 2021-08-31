package types

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsRegistered(t *testing.T) {
	t.Run("registered", func(t *testing.T) {
		user := &User{
			ID:        1,
			Name:      "test name",
			NIM:       "211201XXXXXXXX",
			CreatedAt: time.Now().Format(time.RFC3339),
		}

		r := user.IsRegistered()
		assert.True(t, r)
	})

	t.Run("not registered", func(t *testing.T) {
		user := &User{
			ID: 2,
		}

		r := user.IsRegistered()
		assert.False(t, r)
	})
}
