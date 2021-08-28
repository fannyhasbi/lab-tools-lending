package helper

import (
	"fmt"

	"github.com/fannyhasbi/lab-tools-lending/types"
)

func GroupBorrowStatus(borrows []types.Borrow) map[types.BorrowStatus][]types.Borrow {
	result := make(map[types.BorrowStatus][]types.Borrow)

	for _, borrow := range borrows {
		result[borrow.Status] = append(result[borrow.Status], borrow)
	}

	return result
}

func BuildBorrowedMessage(borrows []types.Borrow) string {
	var message string
	for _, borrow := range borrows {
		message = fmt.Sprintf("%s* %s (%s)\n", message, borrow.Tool.Name, TranslateDateStringToBahasa(borrow.CreatedAt))
	}
	return message
}
