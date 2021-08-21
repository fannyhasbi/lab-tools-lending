package helper

import (
	"strconv"
	"testing"
	"time"

	"github.com/fannyhasbi/lab-tools-lending/types"
	"github.com/stretchr/testify/assert"
)

func TestGetDateFromTimestamp(t *testing.T) {
	tm := time.Now()
	tms := tm.Format(time.RFC3339)

	expected := tm.Format("2006-01-02")

	r := GetDateFromTimestamp(tms)

	assert.Equal(t, expected, r)
}

func TestTranslateDateToBahasa(t *testing.T) {
	date_str := "2021-08-21"
	tm, err := time.Parse("2006-01-02", date_str)

	r := TranslateDateToBahasa(tm)

	expected := "21 Agustus 2021"

	assert.Equal(t, expected, r)
	assert.NoError(t, err)
}

func TestMonthNameSwitcher(t *testing.T) {
	t.Run("switch februari", func(t *testing.T) {
		monthInt := 2
		expected := "Februari"

		r := monthNameSwitcher(monthInt)

		assert.Equal(t, expected, r)
	})

	t.Run("switch desember", func(t *testing.T) {
		monthInt := 12
		expected := "Desember"

		r := monthNameSwitcher(monthInt)

		assert.Equal(t, expected, r)
	})

	t.Run("empty string", func(t *testing.T) {
		monthInt := 999
		expected := ""

		r := monthNameSwitcher(monthInt)

		assert.Equal(t, expected, r)
	})
}

func TestGetBorrowTimeRangeValue(t *testing.T) {
	t.Run("borrow time range map one week", func(t *testing.T) {
		btrm := types.BorrowTimeRangeMap["oneweek"]
		r, err := GetBorrowTimeRangeValue(strconv.Itoa(btrm))
		assert.NoError(t, err)
		assert.Equal(t, 7, r)
	})

	t.Run("borrow time range map one month", func(t *testing.T) {
		btrm := types.BorrowTimeRangeMap["onemonth"]
		r, err := GetBorrowTimeRangeValue(strconv.Itoa(btrm))
		assert.NoError(t, err)
		assert.Equal(t, 30, r)
	})

	t.Run("borrow time range map two month", func(t *testing.T) {
		btrm := types.BorrowTimeRangeMap["twomonth"]
		r, err := GetBorrowTimeRangeValue(strconv.Itoa(btrm))
		assert.NoError(t, err)
		assert.Equal(t, 60, r)
	})

	t.Run("custom time 83", func(t *testing.T) {
		r, err := GetBorrowTimeRangeValue("83")
		assert.NoError(t, err)
		assert.Equal(t, 83, r)
	})

	t.Run("not a number", func(t *testing.T) {
		r, err := GetBorrowTimeRangeValue("abcdefg")
		assert.Error(t, err)
		assert.Equal(t, 0, r)
	})
}
