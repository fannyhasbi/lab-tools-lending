package helper

import (
	"fmt"
	"testing"

	"github.com/fannyhasbi/lab-tools-lending/types"
	"github.com/stretchr/testify/assert"
)

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
