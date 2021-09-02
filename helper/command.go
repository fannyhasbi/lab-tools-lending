package helper

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/fannyhasbi/lab-tools-lending/types"
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

func GetRespondCommands(s string) (types.RespondCommands, bool) {
	ss := strings.Split(s, " ")
	if len(ss) < 3 || len(ss) > 4 {
		return types.RespondCommands{}, false
	}

	ss[1] = strings.ToLower(ss[1])
	resType := types.RespondType(ss[1])
	isExist := IsRespondTypeExists(resType)
	if !isExist {
		return types.RespondCommands{}, false
	}

	id, err := strconv.ParseInt(ss[2], 10, 64)
	if err != nil {
		return types.RespondCommands{}, false
	}

	text := ""
	if len(ss) > 3 {
		text = ss[3]
	}

	result := types.RespondCommands{
		Type: resType,
		ID:   id,
		Text: text,
	}
	return result, true
}

func IsRespondTypeExists(c types.RespondType) bool {
	if c == types.RespondTypeBorrow || c == types.RespondTypeToolReturning {
		return true
	}
	return false
}
