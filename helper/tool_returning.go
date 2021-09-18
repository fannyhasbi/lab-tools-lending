package helper

import (
	"fmt"

	"github.com/fannyhasbi/lab-tools-lending/types"
)

func BuildToolReturningReportMessage(rets []types.ToolReturning) string {
	var message string
	for _, ret := range rets {
		message = fmt.Sprintf(
			"%s[%d] %s - %s, %d buah %s (dikonfirmasi oleh: %s)\n",
			message, ret.ID, TranslateDateToBahasa(ret.ConfirmedAt.Time), ret.Borrow.User.Name, ret.Borrow.Amount, ret.Borrow.Tool.Name, ret.ConfirmedBy.String)
	}
	return message
}
