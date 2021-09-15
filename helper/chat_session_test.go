package helper

import (
	"fmt"
	"testing"

	"github.com/fannyhasbi/lab-tools-lending/types"
	"github.com/stretchr/testify/assert"
)

func TestCanGetChatSessionDetailByTopic(t *testing.T) {
	details := []types.ChatSessionDetail{
		{
			ID:    1,
			Topic: types.Topic["tool_returning_init"],
			Data:  `{"test":"hello"}`,
		},
		{
			ID:    2,
			Topic: types.Topic["tool_returning_confirm"],
			Data:  `{"test":"case"}`,
		},
		{
			ID:    3,
			Topic: types.Topic["tool_returning_complete"],
			Data:  `{"test":"complete"}`,
		},
	}

	t.Run("return the detail", func(t *testing.T) {
		detail, found := GetChatSessionDetailByTopic(details, types.Topic["tool_returning_confirm"])

		assert.True(t, found)
		assert.Equal(t, details[1], detail)
	})

	t.Run("value not found", func(t *testing.T) {
		detail, found := GetChatSessionDetailByTopic(details, types.Topic["borrow_init"])

		assert.False(t, found)
		assert.Equal(t, types.ChatSessionDetail{}, detail)
	})
}

func TestSessionGeneratorRegisterComplete(t *testing.T) {
	resp := true
	gen := NewSessionDataGenerator()
	r := gen.RegisterComplete(resp)

	expected := fmt.Sprintf(`{"type":"%s","user_response":%t}`, string(types.Topic["register_complete"]), resp)

	assert.JSONEq(t, expected, r)
}

func TestSessionGeneratorBorrowInit(t *testing.T) {
	var id int64 = 123
	gen := NewSessionDataGenerator()
	r := gen.BorrowInit(id)

	expected := fmt.Sprintf(`{"type":"%s","tool_id":%d}`, string(types.Topic["borrow_init"]), id)

	assert.JSONEq(t, expected, r)
}

func TestSessionGeneratorBorrowAmount(t *testing.T) {
	amount := 3
	gen := NewSessionDataGenerator()
	r := gen.BorrowAmount(amount)

	expected := fmt.Sprintf(`{"type":"%s","amount":%d}`, string(types.Topic["borrow_amount"]), amount)
	assert.JSONEq(t, expected, r)
}

func TestSessionGeneratorBorrowDuration(t *testing.T) {
	duration := 30
	gen := NewSessionDataGenerator()
	r := gen.BorrowDuration(duration)

	expected := fmt.Sprintf(`{"type":"%s","duration":%d}`, string(types.Topic["borrow_date"]), duration)

	assert.JSONEq(t, expected, r)
}

func TestSessionGeneratorBorrowConfirmation(t *testing.T) {
	resp := true
	gen := NewSessionDataGenerator()
	r := gen.BorrowConfirmation(resp)

	expected := fmt.Sprintf(`{"type":"%s","user_response":%t}`, string(types.Topic["borrow_confirm"]), resp)

	assert.JSONEq(t, expected, r)
}

func TestSessionGeneratorToolReturnConfirm(t *testing.T) {
	additionalInfo := "test keterangan tambahan"
	gen := NewSessionDataGenerator()
	r := gen.ToolReturningConfirm(additionalInfo)

	expected := fmt.Sprintf(`{"type":"%s","additional_info":"%s"}`, string(types.Topic["tool_returning_confirm"]), additionalInfo)

	assert.JSONEq(t, expected, r)
}

func TestSessionGeneratorToolReturningComplete(t *testing.T) {
	resp := true
	gen := NewSessionDataGenerator()
	r := gen.ToolReturningComplete(resp)

	expected := fmt.Sprintf(`{"type":"%s","user_response":%t}`, string(types.Topic["tool_returning_complete"]), resp)

	assert.JSONEq(t, expected, r)
}

func TestSessionGeneratorRespondBorrowInit(t *testing.T) {
	borrowID := int64(123)
	userResponse := "yes"
	gen := NewSessionDataGenerator()
	r := gen.RespondBorrowInit(borrowID, userResponse)

	expected := fmt.Sprintf(`{"type":"%s","borrow_id":%d,"user_response":"%s"}`, string(types.Topic["respond_borrow_init"]), borrowID, userResponse)

	assert.JSONEq(t, expected, r)
}

func TestSessionGeneratorRespondBorrowComplete(t *testing.T) {
	description := "test description"
	gen := NewSessionDataGenerator()
	r := gen.RespondBorrowComplete(description)

	expected := fmt.Sprintf(`{"type":"%s","description":"%s"}`, string(types.Topic["respond_borrow_complete"]), description)
	assert.JSONEq(t, expected, r)
}

func TestSessionGeneratorRespondToolReturningInit(t *testing.T) {
	toolReturningID := int64(123)
	userResponse := "yes"
	gen := NewSessionDataGenerator()
	r := gen.RespondToolReturningInit(toolReturningID, userResponse)

	expected := fmt.Sprintf(`{"type":"%s","tool_returning_id":%d,"user_response":"%s"}`, string(types.Topic["respond_tool_returning_init"]), toolReturningID, userResponse)

	assert.JSONEq(t, expected, r)
}

func TestSessionGeneratorRespondToolReturningComplete(t *testing.T) {
	description := "test description"
	gen := NewSessionDataGenerator()
	r := gen.RespondToolReturningComplete(description)

	expected := fmt.Sprintf(`{"type":"%s","description":"%s"}`, string(types.Topic["respond_tool_returning_complete"]), description)
	assert.JSONEq(t, expected, r)
}

func TestSessionGeneratorManageAdd(t *testing.T) {
	tool := types.Tool{
		Name:                  "Test Tool Name",
		Brand:                 "Test Brand",
		ProductType:           "testPr0duc7Typ3",
		Weight:                123.96,
		Stock:                 32,
		AdditionalInformation: "test additional information",
	}

	t.Run("name", func(t *testing.T) {
		gen := NewSessionDataGenerator()
		r := gen.ManageAddName(tool.Name)
		expected := fmt.Sprintf(`{"type":"%s","name":"%s"}`, types.Topic["manage_add_name"], tool.Name)
		assert.JSONEq(t, expected, r)
	})
	t.Run("brand", func(t *testing.T) {
		gen := NewSessionDataGenerator()
		r := gen.ManageAddBrand(tool.Brand)
		expected := fmt.Sprintf(`{"type":"%s","brand":"%s"}`, types.Topic["manage_add_brand"], tool.Brand)
		assert.JSONEq(t, expected, r)
	})
	t.Run("product type", func(t *testing.T) {
		gen := NewSessionDataGenerator()
		r := gen.ManageAddProductType(tool.ProductType)
		expected := fmt.Sprintf(`{"type":"%s","product_type":"%s"}`, types.Topic["manage_add_type"], tool.ProductType)
		assert.JSONEq(t, expected, r)
	})
	t.Run("weight", func(t *testing.T) {
		gen := NewSessionDataGenerator()
		r := gen.ManageAddWeight(tool.Weight)
		expected := fmt.Sprintf(`{"type":"%s","weight":%.2f}`, types.Topic["manage_add_weight"], tool.Weight)
		assert.JSONEq(t, expected, r)
	})
	t.Run("stock", func(t *testing.T) {
		gen := NewSessionDataGenerator()
		r := gen.ManageAddStock(tool.Stock)
		expected := fmt.Sprintf(`{"type":"%s","stock":%d}`, types.Topic["manage_add_stock"], tool.Stock)
		assert.JSONEq(t, expected, r)
	})
	t.Run("info", func(t *testing.T) {
		gen := NewSessionDataGenerator()
		r := gen.ManageAddInfo(tool.AdditionalInformation)
		expected := fmt.Sprintf(`{"type":"%s","info":"%s"}`, types.Topic["manage_add_info"], tool.AdditionalInformation)
		assert.JSONEq(t, expected, r)
	})
	t.Run("photo", func(t *testing.T) {
		mediaGroupdID := "123"
		fileID := "testFileID1234"
		fileUniqueID := "testFileUniqueID4321"

		gen := NewSessionDataGenerator()
		r := gen.ManageAddPhoto(mediaGroupdID, fileID, fileUniqueID)
		expected := fmt.Sprintf(`{"type":"%s","media_group_id":"%s","file_id":"%s","file_unique_id":"%s"}`, types.Topic["manage_add_photo"], mediaGroupdID, fileID, fileUniqueID)
		assert.JSONEq(t, expected, r)
	})
	t.Run("confirm", func(t *testing.T) {
		userResponse := true
		gen := NewSessionDataGenerator()
		r := gen.ManageAddConfirm(userResponse)
		expected := fmt.Sprintf(`{"type":"%s","user_response":%t}`, types.Topic["manage_add_confirm"], userResponse)
		assert.JSONEq(t, expected, r)
	})
}

func TestSessionGeneratorManageEdit(t *testing.T) {
	t.Run("init", func(t *testing.T) {
		var toolID int64 = 123
		gen := NewSessionDataGenerator()
		r := gen.ManageEditInit(toolID)
		expected := fmt.Sprintf(`{"type":"%s","tool_id":%d}`, types.Topic["manage_edit_init"], toolID)
		assert.JSONEq(t, expected, r)
	})
	t.Run("field", func(t *testing.T) {
		field := "brand"
		gen := NewSessionDataGenerator()
		r := gen.ManageEditField(field)
		expected := fmt.Sprintf(`{"type":"%s","field":"%s"}`, types.Topic["manage_edit_field"], field)
		assert.JSONEq(t, expected, r)
	})
	t.Run("complete", func(t *testing.T) {
		newData := "test new data"
		gen := NewSessionDataGenerator()
		r := gen.ManageEditComplete(newData)
		expected := fmt.Sprintf(`{"type":"%s","new_data":"%s"}`, types.Topic["manage_edit_complete"], newData)
		assert.JSONEq(t, expected, r)
	})
}

func TestSessionGeneratorManagePhoto(t *testing.T) {
	t.Run("init", func(t *testing.T) {
		var toolID int64 = 321
		gen := NewSessionDataGenerator()
		r := gen.ManagePhotoInit(toolID)
		expected := fmt.Sprintf(`{"type":"%s","tool_id":%d}`, types.Topic["manage_photo_init"], toolID)
		assert.JSONEq(t, expected, r)
	})
	t.Run("upload", func(t *testing.T) {
		mediaGroupdID := "123"
		fileID := "testFileID1234"
		fileUniqueID := "testFileUniqueID4321"

		gen := NewSessionDataGenerator()
		r := gen.ManagePhotoUpload(mediaGroupdID, fileID, fileUniqueID)
		expected := fmt.Sprintf(`{"type":"%s","media_group_id":"%s","file_id":"%s","file_unique_id":"%s"}`, types.Topic["manage_photo_upload"], mediaGroupdID, fileID, fileUniqueID)
		assert.JSONEq(t, expected, r)
	})
	t.Run("confirm", func(t *testing.T) {
		userResponse := true
		gen := NewSessionDataGenerator()
		r := gen.ManagePhotoConfirm(userResponse)
		expected := fmt.Sprintf(`{"type":"%s","user_response":%t}`, types.Topic["manage_photo_confirm"], userResponse)
		assert.JSONEq(t, expected, r)
	})
}
