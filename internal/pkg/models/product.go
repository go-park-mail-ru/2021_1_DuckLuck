package models

import (
	"database/sql"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/tools/sanitizer"
)

var (
	ProductsCostSort   = "cost"
	ProductsRatingSort = "rating"

	PaginatorASC  = "ASC"
	PaginatorDESC = "DESC"
)

// All product information
// This models saved in database
type Product struct {
	Id          uint64               `json:"id"`
	Title       string               `json:"title" valid:"minstringlength(3)"`
	Price       ProductPrice         `json:"price" valid:"notnull, json"`
	Rating      float32              `json:"rating" valid:"float, range(0, 10)"`
	Description sql.NullString       `json:"description" valid:"utfletter"`
	Category    []*CategoriesCatalog `json:"category" valid:"notnull"`
	Images      []string             `json:"images" valid:"notnull"`
}

// View of product
// This model contains field for preview information
type ViewProduct struct {
	Id           uint64       `json:"id"`
	Title        string       `json:"title" valid:"minstringlength(3)"`
	Price        ProductPrice `json:"price" valid:"notnull, json"`
	Rating       float32      `json:"rating" valid:"float, range(0, 10)"`
	PreviewImage string       `json:"preview_image" valid:"minstringlength(3)"`
}

// Price product kept in integer nums
// and contains base price and discount
type ProductPrice struct {
	Discount int `json:"discount"`
	BaseCost int `json:"base_cost"`
}

// Set of product with count uniq sets of this size
type RangeProducts struct {
	ListPreviewProducts []*ViewProduct `json:"list_preview_products" valid:"notnull"`
	MaxCountPages       int            `json:"max_count_pages"`
}

// Paginator for showing page of product
type PaginatorProducts struct {
	PageNum       int    `json:"page_num"`
	Count         int    `json:"count"`
	SortKey       string `json:"sort_key" valid:"in(cost|rating)"`
	SortDirection string `json:"sort_direction" valid:"in(ASC|DECS)"`
}

func (pp *PaginatorProducts) Sanitize() {
	sanitizer := sanitizer.NewSanitizer()
	pp.SortKey = sanitizer.Sanitize(pp.SortKey)
	pp.SortDirection = sanitizer.Sanitize(pp.SortDirection)
}
