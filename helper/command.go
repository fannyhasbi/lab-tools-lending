package helper

import (
	"regexp"
	"strings"
)

func GetCommand(message string) string {
	match, err := regexp.MatchString("^/", message)
	if err != nil {
		return ""
	}

	if !match {
		return ""
	}

	spaceIndex := strings.Index(message, " ")
	if spaceIndex > -1 {
		return message[1:spaceIndex]
	}

	return message[1:]
}
