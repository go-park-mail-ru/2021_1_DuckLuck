package models

import "database/sql"

var (
	ProductsCostSort   = "cost"
	ProductsRatingSort = "rating"

	PaginatorASC  = "ASC"
	PaginatorDESC = "DESC"
)

type Product struct {
	Id          uint64               `json:"id"`
	Title       string               `json:"title" valid:"minstringlength(3)"`
	Price       ProductPrice         `json:"price" valid:"notnull, json"`
	Rating      float32              `json:"rating" valid:"float, range(0, 10)"`
	Description sql.NullString       `json:"description" valid:"utfletter"`
	Category    []*CategoriesCatalog `json:"category" valid:"notnull"`
	Images      []string             `json:"images" valid:"notnull"`
}

type ViewProduct struct {
	Id           uint64       `json:"id"`
	Title        string       `json:"title" valid:"minstringlength(3)"`
	Price        ProductPrice `json:"price" valid:"notnull, json"`
	Rating       float32      `json:"rating" valid:"float, range(0, 10)"`
	PreviewImage string       `json:"preview_image" valid:"minstringlength(3)"`
}

type ProductPrice struct {
	Discount int `json:"discount"`
	BaseCost int `json:"base_cost"`
}

type RangeProducts struct {
	ListPreviewProducts []*ViewProduct `json:"list_preview_products" valid:"notnull"`
	MaxCountPages       int            `json:"max_count_pages"`
}

type PaginatorProducts struct {
	PageNum       int    `json:"page_num"`
	Count         int    `json:"count"`
	SortKey       string `json:"sort_key" valid:"in(cost|rating)"`
	SortDirection string `json:"sort_direction" valid:"in(ASC|DECS)"`
}
