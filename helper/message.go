package helper

import "strings"

func SplitNewLine(s string) []string {
	return strings.Split(s, "\n")
}

func RemoveTab(s string) string {
	return strings.Replace(s, "\t", "", -1)
}
