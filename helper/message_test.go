package helper

import (
	"fmt"
	"testing"

	"github.com/fannyhasbi/lab-tools-lending/types"
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

func TestCanBuildToolListMessage(t *testing.T) {
	tools := []types.Tool{
		{
			ID:   123,
			Name: "hello1",
		},
		{
			ID:   321,
			Name: "hello2",
		},
	}

	r := BuildToolListMessage(tools)

	expected := fmt.Sprintf("[%d] %s\n[%d] %s\n", tools[0].ID, tools[0].Name, tools[1].ID, tools[1].Name)

	assert.Equal(t, expected, r)
}
