package helper

import (
	"regexp"
	"strings"
)

func GetCommand(message string) string {
	match, _ := regexp.MatchString("^/", message)

	if !match {
		return ""
	}

	spaceIndex := strings.Index(message, " ")
	if spaceIndex > -1 {
		return message[1:spaceIndex]
	}

	return message[1:]
}
