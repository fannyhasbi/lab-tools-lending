package config

import (
	"os"
)

func WebhookUrl() string {
	var url string = "https://api.telegram.org/" + os.Getenv("BOT_TOKEN")
	return url
}
