package models

type PreviewCart struct {
	Products []*PreviewCartArticle `json:"products"`
}

type PreviewCartArticle struct {
	Id           uint64       `json:"id"`
	Title        string       `json:"title"`
	Price        ProductPrice `json:"price"`
	PreviewImage string       `json:"preview_image"`
	Count        uint64       `json:"count"`
}

type Cart struct {
	Products map[uint64]*ProductPosition `json:"products"`
}

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
