package helper

import (
	"github.com/Jeffail/gabs"
	"github.com/fannyhasbi/lab-tools-lending/types"
)

type SessionDataContainer struct {
	container *gabs.Container
}

func GetChatSessionDetailByTopic(details []types.ChatSessionDetail, topic types.TopicType) (types.ChatSessionDetail, bool) {
	for _, detail := range details {
		if detail.Topic == topic {
			return detail, true
		}
	}
	return types.ChatSessionDetail{}, false
}

func NewSessionDataGenerator() SessionDataContainer {
	return SessionDataContainer{
		container: gabs.New(),
	}
}

func (sdc SessionDataContainer) RegisterComplete(userResponse bool) string {
	sdc.container.Set(types.Topic["register_complete"], "type")
	sdc.container.Set(userResponse, "user_response")
	return sdc.container.String()
}

func (sdc SessionDataContainer) BorrowInit(toolID int64) string {
	sdc.container.Set(types.Topic["borrow_init"], "type")
	sdc.container.Set(toolID, "tool_id")
	return sdc.container.String()
}

func (sdc SessionDataContainer) BorrowAmount(amount int) string {
	sdc.container.Set(types.Topic["borrow_amount"], "type")
	sdc.container.Set(amount, "amount")
	return sdc.container.String()
}

func (sdc SessionDataContainer) BorrowDuration(dateDuration int) string {
	sdc.container.Set(types.Topic["borrow_date"], "type")
	sdc.container.Set(dateDuration, "duration")
	return sdc.container.String()
}

func (sdc SessionDataContainer) BorrowReason(reason string) string {
	sdc.container.Set(types.Topic["borrow_reason"], "type")
	sdc.container.Set(reason, "reason")
	return sdc.container.String()
}

func (sdc SessionDataContainer) BorrowConfirmation(userResponse bool) string {
	sdc.container.Set(types.Topic["borrow_confirm"], "type")
	sdc.container.Set(userResponse, "user_response")
	return sdc.container.String()
}

func (sdc SessionDataContainer) ToolReturningInit(borrowID int64) string {
	sdc.container.Set(types.Topic["tool_returning_init"], "type")
	sdc.container.Set(borrowID, "borrow_id")
	return sdc.container.String()
}

func (sdc SessionDataContainer) ToolReturningConfirm(additionalInfo string) string {
	sdc.container.Set(types.Topic["tool_returning_confirm"], "type")
	sdc.container.Set(additionalInfo, "additional_info")
	return sdc.container.String()
}

func (sdc SessionDataContainer) ToolReturningComplete(userResponse bool) string {
	sdc.container.Set(types.Topic["tool_returning_complete"], "type")
	sdc.container.Set(userResponse, "user_response")
	return sdc.container.String()
}

func (sdc SessionDataContainer) RespondBorrowInit(borrowID int64, userResponse string) string {
	sdc.container.Set(types.Topic["respond_borrow_init"], "type")
	sdc.container.Set(borrowID, "borrow_id")
	sdc.container.Set(userResponse, "user_response")
	return sdc.container.String()
}

func (sdc SessionDataContainer) RespondBorrowComplete(description string) string {
	sdc.container.Set(types.Topic["respond_borrow_complete"], "type")
	sdc.container.Set(description, "description")
	return sdc.container.String()
}

func (sdc SessionDataContainer) RespondToolReturningInit(toolReturningID int64, userResponse string) string {
	sdc.container.Set(types.Topic["respond_tool_returning_init"], "type")
	sdc.container.Set(toolReturningID, "tool_returning_id")
	sdc.container.Set(userResponse, "user_response")
	return sdc.container.String()
}

func (sdc SessionDataContainer) RespondToolReturningComplete(description string) string {
	sdc.container.Set(types.Topic["respond_tool_returning_complete"], "type")
	sdc.container.Set(description, "description")
	return sdc.container.String()
}

func (sdc SessionDataContainer) managePhoto(topic types.TopicType, mediaGroupID, fileID, fileUniqueID string) {
	sdc.container.Set(topic, "type")
	sdc.container.Set(mediaGroupID, "media_group_id")
	sdc.container.Set(fileID, "file_id")
	sdc.container.Set(fileUniqueID, "file_unique_id")
}

func (sdc SessionDataContainer) ManageAddName(name string) string {
	sdc.container.Set(types.Topic["manage_add_name"], "type")
	sdc.container.Set(name, "name")
	return sdc.container.String()
}

func (sdc SessionDataContainer) ManageAddBrand(brand string) string {
	sdc.container.Set(types.Topic["manage_add_brand"], "type")
	sdc.container.Set(brand, "brand")
	return sdc.container.String()
}

func (sdc SessionDataContainer) ManageAddProductType(productType string) string {
	sdc.container.Set(types.Topic["manage_add_type"], "type")
	sdc.container.Set(productType, "product_type")
	return sdc.container.String()
}

func (sdc SessionDataContainer) ManageAddWeight(weight float32) string {
	sdc.container.Set(types.Topic["manage_add_weight"], "type")
	sdc.container.Set(weight, "weight")
	return sdc.container.String()
}

func (sdc SessionDataContainer) ManageAddStock(stock int64) string {
	sdc.container.Set(types.Topic["manage_add_stock"], "type")
	sdc.container.Set(stock, "stock")
	return sdc.container.String()
}

func (sdc SessionDataContainer) ManageAddInfo(info string) string {
	sdc.container.Set(types.Topic["manage_add_info"], "type")
	sdc.container.Set(info, "info")
	return sdc.container.String()
}

func (sdc SessionDataContainer) ManageAddPhoto(mediaGroupID, fileID, fileUniqueID string) string {
	sdc.managePhoto(types.Topic["manage_add_photo"], mediaGroupID, fileID, fileUniqueID)
	return sdc.container.String()
}

func (sdc SessionDataContainer) ManageAddConfirm(userResponse bool) string {
	sdc.container.Set(types.Topic["manage_add_confirm"], "type")
	sdc.container.Set(userResponse, "user_response")
	return sdc.container.String()
}

func (sdc SessionDataContainer) ManageEditInit(toolID int64) string {
	sdc.container.Set(types.Topic["manage_edit_init"], "type")
	sdc.container.Set(toolID, "tool_id")
	return sdc.container.String()
}

func (sdc SessionDataContainer) ManageEditField(field string) string {
	sdc.container.Set(types.Topic["manage_edit_field"], "type")
	sdc.container.Set(field, "field")
	return sdc.container.String()
}

func (sdc SessionDataContainer) ManageEditComplete(newData string) string {
	sdc.container.Set(types.Topic["manage_edit_complete"], "type")
	sdc.container.Set(newData, "new_data")
	return sdc.container.String()
}

func (sdc SessionDataContainer) ManageDeleteInit(toolID int64) string {
	sdc.container.Set(types.Topic["manage_delete_init"], "type")
	sdc.container.Set(toolID, "tool_id")
	return sdc.container.String()
}

func (sdc SessionDataContainer) ManageDeleteComplete(userResponse bool) string {
	sdc.container.Set(types.Topic["manage_delete_complete"], "type")
	sdc.container.Set(userResponse, "user_response")
	return sdc.container.String()
}

func (sdc SessionDataContainer) ManagePhotoInit(toolID int64) string {
	sdc.container.Set(types.Topic["manage_photo_init"], "type")
	sdc.container.Set(toolID, "tool_id")
	return sdc.container.String()
}

func (sdc SessionDataContainer) ManagePhotoUpload(mediaGroupID, fileID, fileUniqueID string) string {
	sdc.managePhoto(types.Topic["manage_photo_upload"], mediaGroupID, fileID, fileUniqueID)
	return sdc.container.String()
}

func (sdc SessionDataContainer) ManagePhotoConfirm(userResponse bool) string {
	sdc.container.Set(types.Topic["manage_photo_confirm"], "type")
	sdc.container.Set(userResponse, "user_response")
	return sdc.container.String()
}
