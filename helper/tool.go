package helper

import (
	"github.com/Jeffail/gabs"
	"github.com/fannyhasbi/lab-tools-lending/types"
)

func GetToolFromChatSessionDetail(details []types.ChatSessionDetail) types.Tool {
	var tool types.Tool

	for _, detail := range details {
		dataParsed, err := gabs.ParseJSON([]byte(detail.Data))
		if err != nil {
			return tool
		}

		switch detail.Topic {
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

	return tool
}

func GetToolPhotosFromChatSessionDetails(details []types.ChatSessionDetail) []types.TelePhotoSize {
	var photos []types.TelePhotoSize
	for _, detail := range details {
		dataParsed, err := gabs.ParseJSON([]byte(detail.Data))
		if err != nil {
			return photos
		}

		if detail.Topic == types.Topic["manage_add_photo"] {
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
