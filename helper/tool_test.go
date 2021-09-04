package helper

import (
	"testing"

	"github.com/fannyhasbi/lab-tools-lending/types"
	"github.com/stretchr/testify/assert"
)

func TestGetToolFromChatSessionDetail(t *testing.T) {
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
				Data:  NewSessionDataGenerator().ManageAddType(tool.ProductType),
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

		r := GetToolFromChatSessionDetail(tools)

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

		r := GetToolFromChatSessionDetail(tools)
		expected := types.Tool{
			Name:  tool.Name,
			Brand: tool.Brand,
		}

		assert.Equal(t, expected, r)
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
