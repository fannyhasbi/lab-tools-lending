package types

type Tools struct {
	ID                    int64   `json:"id"`
	Name                  string  `json:"name"`
	Brand                 string  `json:"brand"`
	ProductType           string  `json:"product_type"`
	Weight                float32 `json:"weight"`
	AdditionalInformation string  `json:"additional_info"`
}
