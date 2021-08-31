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
	dateStr := "2021-08-03"
	tm, err := time.Parse("2006-01-02", dateStr)

	r := TranslateDateToBahasa(tm)

	expected := "3 Agustus 2021"

	assert.Equal(t, expected, r)
	assert.NoError(t, err)
}

func TestTranslateDateStringToBahasa(t *testing.T) {
	t.Run("basic date", func(t *testing.T) {
		dateStr := "2021-08-03"
		r := TranslateDateStringToBahasa(dateStr)
		expected := "3 Agustus 2021"

		assert.Equal(t, expected, r)
	})

	t.Run("date with hour:minute:second", func(t *testing.T) {
		dateStr := "2021-08-22 11:57:58"
		r := TranslateDateStringToBahasa(dateStr)
		expected := "22 Agustus 2021"

		assert.Equal(t, expected, r)
	})

	t.Run("date with RFC format", func(t *testing.T) {
		dateStr := "2021-08-03T12:00:08.971317Z"
		r := TranslateDateStringToBahasa(dateStr)
		expected := "3 Agustus 2021"

		assert.Equal(t, expected, r)
	})
}

func TestMonthNameSwitcher(t *testing.T) {
	t.Run("januari", func(t *testing.T) {
		assert.Equal(t, "Januari", monthNameSwitcher(1))
	})
	t.Run("februari", func(t *testing.T) {
		assert.Equal(t, "Februari", monthNameSwitcher(2))
	})
	t.Run("Maret", func(t *testing.T) {
		assert.Equal(t, "Maret", monthNameSwitcher(3))
	})
	t.Run("April", func(t *testing.T) {
		assert.Equal(t, "April", monthNameSwitcher(4))
	})
	t.Run("Mei", func(t *testing.T) {
		assert.Equal(t, "Mei", monthNameSwitcher(5))
	})
	t.Run("Juni", func(t *testing.T) {
		assert.Equal(t, "Juni", monthNameSwitcher(6))
	})
	t.Run("Juli", func(t *testing.T) {
		assert.Equal(t, "Juli", monthNameSwitcher(7))
	})
	t.Run("Agustus", func(t *testing.T) {
		assert.Equal(t, "Agustus", monthNameSwitcher(8))
	})
	t.Run("September", func(t *testing.T) {
		assert.Equal(t, "September", monthNameSwitcher(9))
	})
	t.Run("Oktober", func(t *testing.T) {
		assert.Equal(t, "Oktober", monthNameSwitcher(10))
	})
	t.Run("November", func(t *testing.T) {
		assert.Equal(t, "November", monthNameSwitcher(11))
	})
	t.Run("Desember", func(t *testing.T) {
		assert.Equal(t, "Desember", monthNameSwitcher(12))
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
