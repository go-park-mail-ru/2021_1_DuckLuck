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

// View of cart
// This model contains field for preview information
type PreviewCart struct {
	Products []*PreviewCartArticle `json:"products" valid:"notnull"`
}

// View of product in cart
// This model contains field for showing product in user cart
type PreviewCartArticle struct {
	Id           uint64       `json:"id"`
	Title        string       `json:"title" valid:"minstringlength(3)"`
	Price        ProductPrice `json:"price" valid:"notnull"`
	PreviewImage string       `json:"preview_image" valid:"minstringlength(3)"`
	Count        uint64       `json:"count"`
}
