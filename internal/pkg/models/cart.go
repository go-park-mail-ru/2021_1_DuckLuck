package models

type PreviewCart struct {
	Products []*PreviewCartArticle `json:"products" valid:"notnull"`
}

type PreviewCartArticle struct {
	Id           uint64       `json:"id"`
	Title        string       `json:"title" valid:"minstringlength(3)"`
	Price        ProductPrice `json:"price" valid:"notnull"`
	PreviewImage string       `json:"preview_image" valid:"minstringlength(3)"`
	Count        uint64       `json:"count"`
}

type Cart struct {
	Products map[uint64]*ProductPosition `json:"products" valid:"notnull"`
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
