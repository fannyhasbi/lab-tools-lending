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

func GetBorrowsByStatus(borrows []types.Borrow, status types.BorrowStatus) []types.Borrow {
	result := []types.Borrow{}
	for _, borrow := range borrows {
		if borrow.Status == status {
			result = append(result, borrow)
		}
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

func BuildBorrowRequestListMessage(borrows []types.Borrow) string {
	var message string
	for _, borrow := range borrows {
		message = fmt.Sprintf("%s[%d] %s - %s\n", message, borrow.ID, borrow.User.Name, borrow.Tool.Name)
	}
	return message
}

func BuildToolReturningRequestListMessage(rets []types.ToolReturning) string {
	var message string
	for _, ret := range rets {
		message = fmt.Sprintf("%s[%d] %s - %s\n", message, ret.ID, ret.User.Name, ret.Tool.Name)
	}
	return message
}
