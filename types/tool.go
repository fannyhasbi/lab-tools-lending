package types

type (
	Tool struct {
		ID                    int64   `json:"id"`
		Name                  string  `json:"name"`
		Brand                 string  `json:"brand"`
		ProductType           string  `json:"product_type"`
		Weight                float32 `json:"weight"`
		Stock                 int64   `json:"stock"`
		AdditionalInformation string  `json:"additional_info"`
		CreatedAt             string  `json:"created_at"`
		UpdatedAt             string  `json:"updated_at"`
	}
)
