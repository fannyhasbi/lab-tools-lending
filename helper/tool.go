package helper

import (
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
	case types.Topic["manage_edit_name"]:
		name, _ := dataParsed.Path("name").Data().(string)
		tool.Name = name
	case types.Topic["manage_edit_brand"]:
		brand, _ := dataParsed.Path("brand").Data().(string)
		tool.Brand = brand
	case types.Topic["manage_edit_type"]:
		product_type, _ := dataParsed.Path("product_type").Data().(string)
		tool.ProductType = product_type
	case types.Topic["manage_edit_weight"]:
		weight, _ := dataParsed.Path("weight").Data().(float64)
		w := float32(weight)
		tool.Weight = w
	case types.Topic["manage_edit_stock"]:
		stock, _ := dataParsed.Path("stock").Data().(float64)
		tool.Stock = int64(stock)
	case types.Topic["manage_edit_info"]:
		info, _ := dataParsed.Path("info").Data().(string)
		tool.AdditionalInformation = info
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
