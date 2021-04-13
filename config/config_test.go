package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPort(t *testing.T) {
	t.Run("use default port", func(t *testing.T) {
		p := GetPort()

		assert.Equal(t, port, p)
	})

	t.Run("can get port using env", func(t *testing.T) {
		os.Setenv("PORT", "1234")

		p := GetPort()

		assert.Equal(t, "1234", p)
	})
}
