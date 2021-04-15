package helper

import (
	"fmt"
	"strings"

	"github.com/fannyhasbi/lab-tools-lending/types"
)

func SplitNewLine(s string) []string {
	return strings.Split(s, "\n")
}

func RemoveTab(s string) string {
	return strings.Replace(s, "\t", "", -1)
}

func BuildToolListMessage(l []types.Tool) string {
	m := ""
	for _, t := range l {
		m = fmt.Sprintf("%s[%d] %s\n", m, t.ID, t.Name)
	}
	return m
}
