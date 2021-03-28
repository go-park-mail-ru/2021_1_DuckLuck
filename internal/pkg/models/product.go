package models

import "database/sql"

var (
	ProductsCostSort   = "cost"
	ProductsRatingSort = "rating"

	PaginatorASC  = "ASC"
	PaginatorDESC = "DESC"
)

type Product struct {
	Id          uint64               `json:"id" valid:"type(uint64)"`
	Title       string               `json:"title" valid:"minstringlength(3)"`
	Price       ProductPrice         `json:"price" valid:"notnull, json"`
	Rating      float32              `json:"rating" valid:"float, range(0, 10)"`
	Description sql.NullString       `json:"description" valid:"utfletter"`
	Category    []*CategoriesCatalog `json:"category" valid:"notnull"`
	Images      []string             `json:"images" valid:"notnull"`
}

type ViewProduct struct {
	Id           uint64       `json:"id" valid:"type(uint64)"`
	Title        string       `json:"title" valid:"minstringlength(3)"`
	Price        ProductPrice `json:"price" valid:"notnull, json"`
	Rating       float32      `json:"rating" valid:"float, range(0, 10)"`
	PreviewImage string       `json:"preview_image" valid:"minstringlength(3)"`
}

type ProductPrice struct {
	Discount int `json:"discount" valid:"int"`
	BaseCost int `json:"base_cost" valid:"int"`
}

type RangeProducts struct {
	ListPreviewProducts []*ViewProduct `json:"list_preview_products" valid:"notnull"`
	MaxCountPages       int            `json:"max_count_pages" valid:"int"`
}

type PaginatorProducts struct {
	PageNum       int    `json:"page_num" valid:"int"`
	Count         int    `json:"count" valid:"int"`
	SortKey       string `json:"sort_key" valid:"in(cost|rating)"`
	SortDirection string `json:"sort_direction" valid:"in(ASC|DECS)"`
}
