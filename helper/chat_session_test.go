package helper

import (
	"fmt"
	"testing"

	"github.com/fannyhasbi/lab-tools-lending/types"
	"github.com/stretchr/testify/assert"
)

func TestCanGetChatSessionDetailByTopic(t *testing.T) {
	details := []types.ChatSessionDetail{
		{
			ID:    1,
			Topic: types.Topic["tool_returning_init"],
			Data:  `{"test":"hello"}`,
		},
		{
			ID:    2,
			Topic: types.Topic["tool_returning_confirm"],
			Data:  `{"test":"case"}`,
		},
		{
			ID:    3,
			Topic: types.Topic["tool_returning_complete"],
			Data:  `{"test":"complete"}`,
		},
	}

	t.Run("return the detail", func(t *testing.T) {
		detail, found := GetChatSessionDetailByTopic(details, types.Topic["tool_returning_confirm"])

		assert.True(t, found)
		assert.Equal(t, details[1], detail)
	})

	t.Run("value not found", func(t *testing.T) {
		detail, found := GetChatSessionDetailByTopic(details, types.Topic["borrow_init"])

		assert.False(t, found)
		assert.Equal(t, types.ChatSessionDetail{}, detail)
	})
}

func TestSessionGeneratorRegisterComplete(t *testing.T) {
	resp := true
	gen := NewSessionDataGenerator()
	r := gen.RegisterComplete(resp)

	expected := fmt.Sprintf(`{"type":"%s","user_response":%t}`, string(types.Topic["register_complete"]), resp)

	assert.JSONEq(t, expected, r)
}

func TestSessionGeneratorBorrowInit(t *testing.T) {
	var id int64 = 123
	gen := NewSessionDataGenerator()
	r := gen.BorrowInit(id)

	expected := fmt.Sprintf(`{"type":"%s","tool_id":%d}`, string(types.Topic["borrow_init"]), id)

	assert.JSONEq(t, expected, r)
}

func TestSessionGeneratorBorrowDateRange(t *testing.T) {
	duration := 30
	gen := NewSessionDataGenerator()
	r := gen.BorrowDateRange(duration)

	expected := fmt.Sprintf(`{"type":"%s","duration":%d}`, string(types.Topic["borrow_date"]), duration)

	assert.JSONEq(t, expected, r)
}

func TestSessionGeneratorBorrowConfirmation(t *testing.T) {
	resp := true
	gen := NewSessionDataGenerator()
	r := gen.BorrowConfirmation(resp)

	expected := fmt.Sprintf(`{"type":"%s","user_response":%t}`, string(types.Topic["borrow_confirm"]), resp)

	assert.JSONEq(t, expected, r)
}

func TestSessionGeneratorToolReturnConfirm(t *testing.T) {
	additionalInfo := "test keterangan tambahan"
	gen := NewSessionDataGenerator()
	r := gen.ToolReturningConfirm(additionalInfo)

	expected := fmt.Sprintf(`{"type":"%s","additional_info":"%s"}`, string(types.Topic["tool_returning_confirm"]), additionalInfo)

	assert.JSONEq(t, expected, r)
}

func TestSessionGeneratorToolReturningComplete(t *testing.T) {
	resp := true
	gen := NewSessionDataGenerator()
	r := gen.ToolReturningComplete(resp)

	expected := fmt.Sprintf(`{"type":"%s","user_response":%t}`, string(types.Topic["tool_returning_complete"]), resp)

	assert.JSONEq(t, expected, r)
}
