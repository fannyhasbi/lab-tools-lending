package postgres

import "time"

func timeNowString() string {
	return time.Now().Format(time.RFC3339)
}
