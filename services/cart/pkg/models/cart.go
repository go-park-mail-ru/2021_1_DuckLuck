package models

// All cart information
// This models saved in database
type Cart struct {
	Products map[uint64]*ProductPosition `json:"products" valid:"notnull"`
}

// One product in cart
type CartArticle struct {
	ProductPosition
	ProductIdentifier
}

type ProductPosition struct {
	Count uint64 `json:"count"`
}

type ProductIdentifier struct {
	ProductId uint64 `json:"product_id"`
}

// Order price kept in integer nums
// and contains base price and discount
type TotalPrice struct {
	TotalDiscount int `json:"total_discount"`
	TotalCost     int `json:"total_cost"`
	TotalBaseCost int `json:"total_base_cost"`
}
