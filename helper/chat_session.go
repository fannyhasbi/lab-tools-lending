package helper

import (
	"github.com/Jeffail/gabs"
	"github.com/fannyhasbi/lab-tools-lending/types"
)

type SessionDataContainer struct {
	container *gabs.Container
}

func NewSessionDataGenerator() SessionDataContainer {
	return SessionDataContainer{
		container: gabs.New(),
	}
}

func (sdc SessionDataContainer) BorrowInit(toolID int64) string {
	sdc.container.Set(types.Topic["borrow_init"], "type")
	sdc.container.Set(toolID, "tool_id")
	return sdc.container.String()
}

func (sdc SessionDataContainer) BorrowDateRange(dateDuration int) string {
	sdc.container.Set(types.Topic["borrow_date"], "type")
	sdc.container.Set(dateDuration, "duration")
	return sdc.container.String()
}

func (sdc SessionDataContainer) BorrowConfirmation(userResponse bool) string {
	sdc.container.Set(types.Topic["borrow_confirm"], "type")
	sdc.container.Set(userResponse, "user_response")
	return sdc.container.String()
}

func (sdc SessionDataContainer) ToolReturningComplete(userResponse bool) string {
	sdc.container.Set(types.Topic["tool_returning_confirm"], "type")
	sdc.container.Set(userResponse, "user_response")
	return sdc.container.String()
}
