package models

type Cart struct {
	Products map[uint64]*ProductPosition `json:"products"`
}

type ProductPosition struct {
	ProductId uint64 `json:"product_id"`
	Count     uint64 `json:"count"`
}
