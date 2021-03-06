package helper

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/fannyhasbi/lab-tools-lending/types"
)

func GetCommand(msg string) string {
	message := msg
	match, _ := regexp.MatchString("^/", message)

	if !match {
		return ""
	}

	mentionIndex := strings.Index(message, "@")
	if mentionIndex > -1 {
		message = message[:mentionIndex]
	}

	spaceIndex := strings.Index(message, " ")
	if spaceIndex > -1 {
		return message[1:spaceIndex]
	}

	return message[1:]
}

func GetRespondCommandOrder(s string) (types.RespondCommandOrder, bool) {
	ss := strings.Split(s, " ")
	if len(ss) < 3 || len(ss) > 4 {
		return types.RespondCommandOrder{}, false
	}

	ss[1] = strings.ToLower(ss[1])
	resType := types.RespondType(ss[1])
	isExist := isRespondTypeExists(resType)
	if !isExist {
		return types.RespondCommandOrder{}, false
	}

	id, err := strconv.ParseInt(ss[2], 10, 64)
	if err != nil {
		return types.RespondCommandOrder{}, false
	}

	text := ""
	if len(ss) > 3 {
		text = ss[3]
	}

	result := types.RespondCommandOrder{
		Type: resType,
		ID:   id,
		Text: text,
	}
	return result, true
}

func isRespondTypeExists(c types.RespondType) bool {
	if c == types.RespondTypeBorrow || c == types.RespondTypeToolReturning {
		return true
	}
	return false
}

func isManageTypeExists(c types.ManageType) bool {
	if c == types.ManageTypeAdd || c == types.ManageTypeEdit || c == types.ManageTypeDelete || c == types.ManageTypePhoto {
		return true
	}
	return false
}

func isReportTypeExists(c types.ReportType) bool {
	if c == types.ReportTypeBorrow || c == types.ReportTypeToolReturning {
		return true
	}
	return false
}

func GetManageCommandOrder(s string) (types.ManageCommandOrder, bool) {
	ss := strings.Split(s, " ")
	if len(ss) < 2 || len(ss) > 3 {
		return types.ManageCommandOrder{}, false
	}

	ss[1] = strings.ToLower(ss[1])
	manageType := types.ManageType(ss[1])
	isExist := isManageTypeExists(manageType)
	if !isExist {
		return types.ManageCommandOrder{}, false
	}

	if len(ss) == 2 {
		if manageType == types.ManageTypeAdd || manageType == types.ManageTypeEdit || manageType == types.ManageTypeDelete {
			return types.ManageCommandOrder{Type: manageType}, true
		}
		return types.ManageCommandOrder{}, false
	}

	id, err := strconv.ParseInt(ss[2], 10, 64)
	if err != nil {
		return types.ManageCommandOrder{}, false
	}

	result := types.ManageCommandOrder{
		Type: manageType,
		ID:   id,
	}
	return result, true
}

func GetCheckCommandOrder(s string) (types.CheckCommandOrder, bool) {
	ss := strings.Split(s, " ")
	if len(ss) < 2 || len(ss) > 3 {
		return types.CheckCommandOrder{}, false
	}

	i, err := strconv.ParseInt(ss[1], 10, 64)
	if err != nil {
		return types.CheckCommandOrder{}, false
	}

	if len(ss) == 2 {
		return types.CheckCommandOrder{ID: i}, true
	}

	return types.CheckCommandOrder{ID: i, Text: ss[2]}, true
}

func GetReportCommandOrder(s string) (types.ReportCommandOrder, bool) {
	ss := strings.Split(s, " ")
	if len(ss) < 2 || len(ss) > 3 {
		return types.ReportCommandOrder{}, false
	}

	ss[1] = strings.ToLower(ss[1])
	reportType := types.ReportType(ss[1])
	if isExist := isReportTypeExists(reportType); !isExist {
		return types.ReportCommandOrder{}, false
	}

	if len(ss) == 2 {
		return types.ReportCommandOrder{Type: reportType}, true
	}

	return types.ReportCommandOrder{Type: reportType, Text: ss[2]}, true
}
