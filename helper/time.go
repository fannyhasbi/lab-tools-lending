package helper

import (
	"fmt"
	"strconv"
	"time"
)

func GetDateFromTimestamp(s string) string {
	return s[:10]
}

func TranslateDateToBahasa(date time.Time) string {
	month := monthNameSwitcher(int(date.Month()))
	return fmt.Sprintf("%d %s %d", date.Day(), month, date.Year())
}

func monthNameSwitcher(month int) (m string) {
	switch month {
	case 1:
		m = "Januari"
	case 2:
		m = "Februari"
	case 3:
		m = "Maret"
	case 4:
		m = "April"
	case 5:
		m = "Mei"
	case 6:
		m = "Juni"
	case 7:
		m = "Juli"
	case 8:
		m = "Agustus"
	case 9:
		m = "September"
	case 10:
		m = "Oktober"
	case 11:
		m = "November"
	case 12:
		m = "Desember"
	}

	return
}

func GetBorrowTimeRangeValue(message string) (r int, err error) {
	return strconv.Atoi(message)
}
