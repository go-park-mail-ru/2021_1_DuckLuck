package models

type PromoCodeGroup struct {
	Products []uint64 `json:"products" valid:"notnull"`
	PromoCode string `json:"promo_code" valid:"utfletter, stringlength(1|30)"`
}

type DiscountedPrice struct {
	TotalDiscount int `json:"total_discount"`
	TotalCost     int `json:"total_cost"`
	TotalBaseCost int `json:"total_base_cost"`
}
