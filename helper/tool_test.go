package helper

import (
	"testing"

	"github.com/fannyhasbi/lab-tools-lending/types"
	"github.com/stretchr/testify/assert"
)

func TestGetToolFromChatSessionDetail_Add(t *testing.T) {
	tool := types.Tool{
		Name:                  "Test Name",
		Brand:                 "Test Brand",
		ProductType:           "TestPr0duc7Typ3",
		Weight:                173,
		Stock:                 15,
		AdditionalInformation: "test additional info",
	}

	t.Run("full tool session", func(t *testing.T) {
		tools := []types.ChatSessionDetail{
			{
				Topic: types.Topic["manage_add_name"],
				Data:  NewSessionDataGenerator().ManageAddName(tool.Name),
			},
			{
				Topic: types.Topic["manage_add_brand"],
				Data:  NewSessionDataGenerator().ManageAddBrand(tool.Brand),
			},
			{
				Topic: types.Topic["manage_add_type"],
				Data:  NewSessionDataGenerator().ManageAddProductType(tool.ProductType),
			},
			{
				Topic: types.Topic["manage_add_weight"],
				Data:  NewSessionDataGenerator().ManageAddWeight(tool.Weight),
			},
			{
				Topic: types.Topic["manage_add_stock"],
				Data:  NewSessionDataGenerator().ManageAddStock(tool.Stock),
			},
			{
				Topic: types.Topic["manage_add_info"],
				Data:  NewSessionDataGenerator().ManageAddInfo(tool.AdditionalInformation),
			},
		}

		r := GetToolFromChatSessionDetail(types.ManageTypeAdd, tools)

		assert.Equal(t, tool, r)
	})

	t.Run("not full session", func(t *testing.T) {
		tools := []types.ChatSessionDetail{
			{
				Topic: types.Topic["manage_add_name"],
				Data:  NewSessionDataGenerator().ManageAddName(tool.Name),
			},
			{
				Topic: types.Topic["manage_add_brand"],
				Data:  NewSessionDataGenerator().ManageAddBrand(tool.Brand),
			},
		}

		r := GetToolFromChatSessionDetail(types.ManageTypeAdd, tools)
		expected := types.Tool{
			Name:  tool.Name,
			Brand: tool.Brand,
		}

		assert.Equal(t, expected, r)
	})
}

func TestGetToolFromChatSessionDetail_Edit(t *testing.T) {
	tool := types.Tool{
		ID:                    111,
		Name:                  "Test Name Edit",
		Brand:                 "Test Brand Edit",
		ProductType:           "TestPr0duc7Typ3Edit",
		Weight:                333,
		Stock:                 32,
		AdditionalInformation: "test additional info edit",
	}

	t.Run("full tool session", func(t *testing.T) {
		tools := []types.ChatSessionDetail{
			{
				Topic: types.Topic["manage_edit_init"],
				Data:  NewSessionDataGenerator().ManageEditInit(tool.ID),
			},
			{
				Topic: types.Topic["manage_edit_name"],
				Data:  NewSessionDataGenerator().ManageEditName(tool.Name),
			},
			{
				Topic: types.Topic["manage_edit_brand"],
				Data:  NewSessionDataGenerator().ManageEditBrand(tool.Brand),
			},
			{
				Topic: types.Topic["manage_edit_type"],
				Data:  NewSessionDataGenerator().ManageEditProductType(tool.ProductType),
			},
			{
				Topic: types.Topic["manage_edit_weight"],
				Data:  NewSessionDataGenerator().ManageEditWeight(tool.Weight),
			},
			{
				Topic: types.Topic["manage_edit_stock"],
				Data:  NewSessionDataGenerator().ManageEditStock(tool.Stock),
			},
			{
				Topic: types.Topic["manage_edit_info"],
				Data:  NewSessionDataGenerator().ManageEditInfo(tool.AdditionalInformation),
			},
		}

		r := GetToolFromChatSessionDetail(types.ManageTypeEdit, tools)

		assert.Equal(t, tool, r)
	})

	t.Run("not full session", func(t *testing.T) {
		tools := []types.ChatSessionDetail{
			{
				Topic: types.Topic["manage_edit_init"],
				Data:  NewSessionDataGenerator().ManageEditInit(tool.ID),
			},
			{
				Topic: types.Topic["manage_edit_name"],
				Data:  NewSessionDataGenerator().ManageAddName(tool.Name),
			},
			{
				Topic: types.Topic["manage_edit_brand"],
				Data:  NewSessionDataGenerator().ManageAddBrand(tool.Brand),
			},
		}

		r := GetToolFromChatSessionDetail(types.ManageTypeEdit, tools)
		expected := types.Tool{
			ID:    tool.ID,
			Name:  tool.Name,
			Brand: tool.Brand,
		}

		assert.Equal(t, expected, r)
	})
}

func TestGetToolFromChatSessionDetail_Photo(t *testing.T) {
	tool := types.Tool{
		ID: 123,
	}

	t.Run("full session", func(t *testing.T) {
		sessions := []types.ChatSessionDetail{
			{
				Topic: types.Topic["manage_photo_init"],
				Data:  NewSessionDataGenerator().ManagePhotoInit(tool.ID),
			},
		}

		r := GetToolFromChatSessionDetail(types.ManageTypePhoto, sessions)

		assert.Equal(t, types.Tool{ID: tool.ID}, r)
	})
}

func TestCanGetToolPhotosFromChatSessionDetails(t *testing.T) {
	mediaGroupID := "123"
	photos := []types.TelePhotoSize{
		{
			FileID:       "testfileid1",
			FileUniqueID: "testfileuniqueid1",
		},
		{
			FileID:       "testfileid2",
			FileUniqueID: "testfileuniqueid2",
		},
	}

	t.Run("add photo session", func(t *testing.T) {
		sessions := []types.ChatSessionDetail{
			{
				Topic: types.Topic["manage_add_brand"],
				Data:  NewSessionDataGenerator().ManageAddBrand("test brand"),
			},
			{
				Topic: types.Topic["manage_add_photo"],
				Data:  NewSessionDataGenerator().ManageAddPhoto(mediaGroupID, photos[0].FileID, photos[0].FileUniqueID),
			},
			{
				Topic: types.Topic["manage_add_photo"],
				Data:  NewSessionDataGenerator().ManageAddPhoto(mediaGroupID, photos[1].FileID, photos[1].FileUniqueID),
			},
		}

		r := GetToolPhotosFromChatSessionDetails(sessions)

		assert.Equal(t, photos, r)
	})

	t.Run("edit photo session", func(t *testing.T) {
		var toolID int64 = 555
		sessions := []types.ChatSessionDetail{
			{
				Topic: types.Topic["manage_photo_init"],
				Data:  NewSessionDataGenerator().ManagePhotoInit(toolID),
			},
			{
				Topic: types.Topic["manage_photo_upload"],
				Data:  NewSessionDataGenerator().ManagePhotoUpload(mediaGroupID, photos[0].FileID, photos[0].FileUniqueID),
			},
			{
				Topic: types.Topic["manage_photo_upload"],
				Data:  NewSessionDataGenerator().ManagePhotoUpload(mediaGroupID, photos[1].FileID, photos[1].FileUniqueID),
			},
		}

		r := GetToolPhotosFromChatSessionDetails(sessions)

		assert.Equal(t, photos, r)
	})
}

func TestCanPickPhoto(t *testing.T) {
	t.Run("best quality", func(t *testing.T) {
		photos := []types.TelePhotoSize{
			{
				FileID:   "file1",
				FileSize: 10,
			},
			{
				FileID:   "file2",
				FileSize: 100,
			},
			{
				FileID:   "file3",
				FileSize: 50,
			},
		}

		r := PickPhoto(photos)
		assert.Equal(t, photos[1], r)
	})

	t.Run("first index", func(t *testing.T) {
		photos := []types.TelePhotoSize{
			{
				FileID:   "file1",
				FileSize: 10,
			},
			{
				FileID:   "file2",
				FileSize: 10,
			},
			{
				FileID:   "file3",
				FileSize: 10,
			},
		}

		r := PickPhoto(photos)
		assert.Equal(t, photos[0], r)
	})
}
