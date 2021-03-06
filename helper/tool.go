package helper

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Jeffail/gabs"
	"github.com/fannyhasbi/lab-tools-lending/types"
)

func GetToolFromChatSessionDetail(manageType types.ManageType, details []types.ChatSessionDetail) types.Tool {
	var tool types.Tool

	for _, detail := range details {
		dataParsed, err := gabs.ParseJSON([]byte(detail.Data))
		if err != nil {
			return tool
		}

		if manageType == types.ManageTypeAdd {
			extractToolAddBasedOnTopic(&tool, detail.Topic, dataParsed)
		} else if manageType == types.ManageTypeEdit {
			extractToolEditBasedOnTopic(&tool, detail.Topic, dataParsed)
		} else if manageType == types.ManageTypeDelete {
			extractToolDeleteBasedOnTopic(&tool, detail.Topic, dataParsed)
		} else if manageType == types.ManageTypePhoto {
			extractToolPhotoBasedOnTopic(&tool, detail.Topic, dataParsed)
		}
	}

	return tool
}

func extractToolAddBasedOnTopic(tool *types.Tool, topic types.TopicType, dataParsed *gabs.Container) {
	switch topic {
	case types.Topic["manage_add_name"]:
		name, _ := dataParsed.Path("name").Data().(string)
		tool.Name = name
	case types.Topic["manage_add_brand"]:
		brand, _ := dataParsed.Path("brand").Data().(string)
		tool.Brand = brand
	case types.Topic["manage_add_type"]:
		product_type, _ := dataParsed.Path("product_type").Data().(string)
		tool.ProductType = product_type
	case types.Topic["manage_add_weight"]:
		weight, _ := dataParsed.Path("weight").Data().(float64)
		w := float32(weight)
		tool.Weight = w
	case types.Topic["manage_add_stock"]:
		stock, _ := dataParsed.Path("stock").Data().(float64)
		tool.Stock = int64(stock)
	case types.Topic["manage_add_info"]:
		info, _ := dataParsed.Path("info").Data().(string)
		tool.AdditionalInformation = info
	}
}

func extractToolEditBasedOnTopic(tool *types.Tool, topic types.TopicType, dataParsed *gabs.Container) {
	switch topic {
	case types.Topic["manage_edit_init"]:
		id, _ := dataParsed.Path("tool_id").Data().(float64)
		tool.ID = int64(id)
	}
}

func extractToolDeleteBasedOnTopic(tool *types.Tool, topic types.TopicType, dataParsed *gabs.Container) {
	switch topic {
	case types.Topic["manage_delete_init"]:
		id, _ := dataParsed.Path("tool_id").Data().(float64)
		tool.ID = int64(id)
	}
}

func extractToolPhotoBasedOnTopic(tool *types.Tool, topic types.TopicType, dataParsed *gabs.Container) {
	switch topic {
	case types.Topic["manage_photo_init"]:
		id, _ := dataParsed.Path("tool_id").Data().(float64)
		tool.ID = int64(id)
	}
}

func GetToolPhotosFromChatSessionDetails(details []types.ChatSessionDetail) []types.TelePhotoSize {
	var photos []types.TelePhotoSize
	for _, detail := range details {
		dataParsed, err := gabs.ParseJSON([]byte(detail.Data))
		if err != nil {
			return photos
		}

		if detail.Topic == types.Topic["manage_add_photo"] || detail.Topic == types.Topic["manage_photo_upload"] {
			fileID, _ := dataParsed.Path("file_id").Data().(string)
			fileUniqueID, _ := dataParsed.Path("file_unique_id").Data().(string)

			photo := types.TelePhotoSize{
				FileID:       fileID,
				FileUniqueID: fileUniqueID,
			}

			photos = append(photos, photo)
		}
	}

	return photos
}

func PickPhoto(photos []types.TelePhotoSize) types.TelePhotoSize {
	highestSizePhoto := photos[0]
	for i := range photos {
		if photos[i].FileSize > highestSizePhoto.FileSize {
			highestSizePhoto = photos[i]
		}
	}
	return highestSizePhoto
}

func toolFields() []types.ToolField {
	return []types.ToolField{
		types.ToolFieldName,
		types.ToolFieldBrand,
		types.ToolFieldProductType,
		types.ToolFieldWeight,
		types.ToolFieldStock,
		types.ToolFieldAdditionalInfo,
		types.ToolFieldPhoto,
	}
}

func IsToolFieldExists(f string) bool {
	fAsserted := types.ToolField(f)
	fields := toolFields()
	for _, field := range fields {
		if fAsserted == field {
			return true
		}
	}
	return false
}

func GetToolValueByField(tool types.Tool, f string) string {
	var result string
	fAsserted := types.ToolField(f)

	switch fAsserted {
	case types.ToolFieldName:
		result = tool.Name
	case types.ToolFieldBrand:
		result = tool.Brand
	case types.ToolFieldProductType:
		result = tool.ProductType
	case types.ToolFieldWeight:
		result = fmt.Sprintf("%.2f", tool.Weight)
	case types.ToolFieldStock:
		result = strconv.FormatInt(tool.Stock, 10)
	case types.ToolFieldAdditionalInfo:
		result = tool.AdditionalInformation
	default:
		result = ""
	}

	return result
}

func ChangeToolValueByField(tool types.Tool, field, newValue string) (types.Tool, error) {
	// make a new copy (clean code)
	updatedTool := tool
	fieldAsserted := types.ToolField(field)

	switch fieldAsserted {
	case types.ToolFieldName:
		updatedTool.Name = newValue
	case types.ToolFieldBrand:
		updatedTool.Brand = newValue
	case types.ToolFieldProductType:
		updatedTool.ProductType = newValue
	case types.ToolFieldWeight:
		i, err := strconv.ParseFloat(newValue, 10)
		if err != nil || i < 0 {
			return updatedTool, errors.New("mohon sebutkan berat dalam angka")
		}
		updatedTool.Weight = float32(i)

	case types.ToolFieldStock:
		i, err := strconv.ParseInt(newValue, 10, 64)
		if err != nil || i < 0 {
			return updatedTool, errors.New("mohon sebutkan jumlah stok dalam angka")
		}
		updatedTool.Stock = i

	case types.ToolFieldAdditionalInfo:
		updatedTool.AdditionalInformation = newValue
	}

	return updatedTool, nil
}
