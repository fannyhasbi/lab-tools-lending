package helper

import (
	"database/sql"
	"fmt"

	"github.com/Jeffail/gabs"
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
	layout := "02/01/2006"
	for _, borrow := range borrows {
		since := borrow.ConfirmedAt.Time.Format(layout)
		until := borrow.ConfirmedAt.Time.AddDate(0, 0, borrow.Duration).Format(layout)

		message = fmt.Sprintf("%s[%d] %s (%s - %s)\n", message, borrow.ID, borrow.Tool.Name, since, until)
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
		message = fmt.Sprintf("%s[%d] %s - %s\n", message, ret.ID, ret.Borrow.User.Name, ret.Borrow.Tool.Name)
	}
	return message
}

func GetBorrowFromChatSessionDetail(details []types.ChatSessionDetail) types.Borrow {
	var borrow types.Borrow

	for _, detail := range details {
		dataParsed, err := gabs.ParseJSON([]byte(detail.Data))
		if err != nil {
			return borrow
		}

		switch detail.Topic {
		case types.Topic["borrow_init"]:
			toolID, _ := dataParsed.Path("tool_id").Data().(float64)
			borrow.ToolID = int64(toolID)
		case types.Topic["borrow_amount"]:
			amount, _ := dataParsed.Path("amount").Data().(float64)
			borrow.Amount = int(amount)
		case types.Topic["borrow_date"]:
			duration, _ := dataParsed.Path("duration").Data().(float64)
			borrow.Duration = int(duration)
		case types.Topic["borrow_reason"]:
			reason, _ := dataParsed.Path("reason").Data().(string)
			borrow.Reason = sql.NullString{Valid: true, String: reason}
		}
	}

	return borrow
}

func GetSameBorrow(borrows []types.Borrow, toolID int64) (types.BorrowStatus, bool) {
	for _, b := range borrows {
		if b.ToolID == toolID {
			return b.Status, true
		}
	}

	return "", false
}

func BuildBorrowReportMessage(borrows []types.Borrow) string {
	var message string
	for _, borrow := range borrows {
		message = fmt.Sprintf(
			"%s[%d] %s - %s, %d buah %s (dikonfirmasi oleh: %s)\n",
			message, borrow.ID, TranslateDateToBahasa(borrow.ConfirmedAt.Time), borrow.User.Name, borrow.Amount, borrow.Tool.Name, borrow.ConfirmedBy.String)
	}
	return message
}
