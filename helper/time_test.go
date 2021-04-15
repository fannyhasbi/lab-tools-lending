package helper

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetDateFromTimestamp(t *testing.T) {
	tm := time.Now()
	tms := tm.Format(time.RFC3339)

	expected := tm.Format("2006-01-02")

	r := GetDateFromTimestamp(tms)

	assert.Equal(t, expected, r)
}
