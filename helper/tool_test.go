package helper

import (
	"fmt"
	"strconv"
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
		ID: 111,
	}

	t.Run("full session", func(t *testing.T) {
		sessions := []types.ChatSessionDetail{
			{
				Topic: types.Topic["manage_edit_init"],
				Data:  NewSessionDataGenerator().ManageEditInit(tool.ID),
			},
		}

		r := GetToolFromChatSessionDetail(types.ManageTypeEdit, sessions)

		assert.Equal(t, tool, r)
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

func TestIsToolFieldExists(t *testing.T) {
	t.Run("correct", func(t *testing.T) {
		for _, field := range toolFields() {
			r := IsToolFieldExists(string(field))
			assert.True(t, r)
		}
	})
	t.Run("incorrect", func(t *testing.T) {
		f := "testincorrectfield"
		r := IsToolFieldExists(f)
		assert.False(t, r)
	})
}

func TestGetToolValueByField(t *testing.T) {
	tool := types.Tool{
		Name:                  "Test Tool Name",
		Brand:                 "Test Tool Brand",
		ProductType:           "Test product",
		Weight:                123,
		Stock:                 321,
		AdditionalInformation: "test additional info",
	}

	t.Run("name", func(t *testing.T) {
		r := GetToolValueByField(tool, string(types.ToolFieldName))
		assert.Equal(t, tool.Name, r)
	})
	t.Run("brand", func(t *testing.T) {
		r := GetToolValueByField(tool, string(types.ToolFieldBrand))
		assert.Equal(t, tool.Brand, r)
	})
	t.Run("type", func(t *testing.T) {
		r := GetToolValueByField(tool, string(types.ToolFieldProductType))
		assert.Equal(t, tool.ProductType, r)
	})
	t.Run("weight", func(t *testing.T) {
		r := GetToolValueByField(tool, string(types.ToolFieldWeight))
		assert.Equal(t, fmt.Sprintf("%.2f", tool.Weight), r)
	})
	t.Run("stock", func(t *testing.T) {
		r := GetToolValueByField(tool, string(types.ToolFieldStock))
		assert.Equal(t, strconv.FormatInt(tool.Stock, 10), r)
	})
	t.Run("additional info", func(t *testing.T) {
		r := GetToolValueByField(tool, string(types.ToolFieldAdditionalInfo))
		assert.Equal(t, tool.AdditionalInformation, r)
	})
	t.Run("no case", func(t *testing.T) {
		r := GetToolValueByField(tool, "testnocase")
		assert.Equal(t, "", r)
	})
}

func TestCanChangeToolValueByField(t *testing.T) {
	tool := types.Tool{
		Name:                  "Test Tool Name",
		Brand:                 "Test Brand",
		ProductType:           "Test product type",
		Weight:                123,
		Stock:                 321,
		AdditionalInformation: "test additional info",
	}

	t.Run("name", func(t *testing.T) {
		newTool := tool
		newTool.Name = "New Test Tool Name"
		r, err := ChangeToolValueByField(tool, string(types.ToolFieldName), newTool.Name)
		assert.NoError(t, err)
		assert.Equal(t, newTool, r)
	})
	t.Run("brand", func(t *testing.T) {
		newTool := tool
		newTool.Brand = "New Test Brand"
		r, err := ChangeToolValueByField(tool, string(types.ToolFieldBrand), newTool.Brand)
		assert.NoError(t, err)
		assert.Equal(t, newTool, r)
	})
	t.Run("type", func(t *testing.T) {
		newTool := tool
		newTool.ProductType = "New test product type"
		r, err := ChangeToolValueByField(tool, string(types.ToolFieldProductType), newTool.ProductType)
		assert.NoError(t, err)
		assert.Equal(t, newTool, r)
	})
	t.Run("weight", func(t *testing.T) {
		newTool := tool
		newTool.Weight = 999
		r, err := ChangeToolValueByField(tool, string(types.ToolFieldWeight), fmt.Sprintf("%f", newTool.Weight))
		assert.NoError(t, err)
		assert.Equal(t, newTool, r)
	})
	t.Run("stock", func(t *testing.T) {
		newTool := tool
		newTool.Stock = 333
		r, err := ChangeToolValueByField(tool, string(types.ToolFieldStock), strconv.FormatInt(newTool.Stock, 10))
		assert.NoError(t, err)
		assert.Equal(t, newTool, r)
	})
	t.Run("additional info", func(t *testing.T) {
		newTool := tool
		newTool.AdditionalInformation = "new test additional info"
		r, err := ChangeToolValueByField(tool, string(types.ToolFieldAdditionalInfo), newTool.AdditionalInformation)
		assert.NoError(t, err)
		assert.Equal(t, newTool, r)
	})
	t.Run("weight negative", func(t *testing.T) {
		weight := -100
		r, err := ChangeToolValueByField(tool, string(types.ToolFieldWeight), strconv.Itoa(weight))
		assert.Error(t, err)
		assert.Equal(t, tool, r)
	})
	t.Run("weight is not a number", func(t *testing.T) {
		r, err := ChangeToolValueByField(tool, string(types.ToolFieldWeight), "testnotanumber")
		assert.Error(t, err)
		assert.Equal(t, tool, r)
	})
	t.Run("Stock negative", func(t *testing.T) {
		stock := -90
		r, err := ChangeToolValueByField(tool, string(types.ToolFieldStock), strconv.Itoa(stock))
		assert.Error(t, err)
		assert.Equal(t, tool, r)
	})
	t.Run("stock is not a number", func(t *testing.T) {
		r, err := ChangeToolValueByField(tool, string(types.ToolFieldStock), "testnotanumber")
		assert.Error(t, err)
		assert.Equal(t, tool, r)
	})
}
