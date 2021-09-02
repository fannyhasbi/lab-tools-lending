package helper

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	BasicDateLayout = "2006-01-02"
)

func GetDateFromTimestamp(s string) string {
	return s[:10]
}

func TranslateDateToBahasa(date time.Time) string {
	month := monthNameSwitcher(int(date.Month()))
	return fmt.Sprintf("%d %s %d", date.Day(), month, date.Year())
}

// changeDateStringFormat change into "YYYY-MM-DD" format
func ChangeDateStringFormat(date string) string {
	// anticipate if there is pattern "YYYY-MM-DD HH:mm:ss" and "YYYY-MM-DDTHH:mm:ss"
	dateStr := strings.Split(date, " ")
	dateStr = strings.Split(dateStr[0], "T")
	return dateStr[0]
}

// TranslateDateStringToBahasa parameter date in "YYYY-MM-DD"
func TranslateDateStringToBahasa(date string) string {
	result := ChangeDateStringFormat(date)
	dateStr := strings.Split(result, "-")

	year, _ := strconv.Atoi(dateStr[0])
	month, _ := strconv.Atoi(dateStr[1])
	day, _ := strconv.Atoi(dateStr[2])
	return fmt.Sprintf("%d %s %d", day, monthNameSwitcher(month), year)
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
