package helper

import (
	"testing"

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
