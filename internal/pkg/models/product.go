package models

import "github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/tools/sanitizer"

var (
	ProductsCostSort      = "cost"
	ProductsRatingSort    = "rating"
	ProductsDateAddedSort = "date"
	ProductsDiscountSort  = "discount"

	PaginatorASC  = "ASC"
	PaginatorDESC = "DESC"
)

// All product information
// This models saved in database
type Product struct {
	Id           uint64               `json:"id"`
	Title        string               `json:"title" valid:"minstringlength(3)"`
	Price        ProductPrice         `json:"price" valid:"notnull, json"`
	Rating       float32              `json:"rating" valid:"float, range(0, 10)"`
	Description  string               `json:"description" valid:"utfletter"`
	Category     uint64               `json:"category"`
	CategoryPath []*CategoriesCatalog `json:"category_path" valid:"notnull"`
	Images       []string             `json:"images" valid:"notnull"`
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
	Discount  int `json:"discount"`
	BaseCost  int `json:"base_cost"`
	TotalCost int `json:"total_cost"`
}

// Set of product with count uniq sets of this size
type RangeProducts struct {
	ListPreviewProducts []*ViewProduct `json:"list_preview_products" valid:"notnull"`
	MaxCountPages       int            `json:"max_count_pages"`
}

type SortOptions struct {
	SortKey       string `json:"sort_key" valid:"in(cost|rating|date|discount)"`
	SortDirection string `json:"sort_direction" valid:"in(ASC|DESC)"`
}

// Paginator for showing page of product
type PaginatorProducts struct {
	PageNum       int            `json:"page_num"`
	Count         int            `json:"count"`
	Category      uint64         `json:"category"`
	Filter        *ProductFilter `json:"filter"`
	SortOptions
}

func (pp *PaginatorProducts) Sanitize() {
	sanitizer := sanitizer.NewSanitizer()
	pp.SortKey = sanitizer.Sanitize(pp.SortKey)
	pp.SortDirection = sanitizer.Sanitize(pp.SortDirection)
}

type ProductFilter struct {
	MinPrice   uint64 `json:"min_price"`
	MaxPrice   uint64 `json:"max_price"`
	IsNew      bool   `json:"is_new"`
	IsRating   bool   `json:"is_rating"`
	IsDiscount bool   `json:"is_discount"`
}

// Search query with options
type SearchQuery struct {
	QueryString string `json:"query_string" valid:"minstringlength(2)"`
	PageNum     int    `json:"page_num"`
	Count       int    `json:"count"`
	SortOptions
	Category uint64 `json:"category"`
}

func (sq *SearchQuery) Sanitize() {
	sanitizer := sanitizer.NewSanitizer()
	sq.QueryString = sanitizer.Sanitize(sq.QueryString)
	sq.SortKey = sanitizer.Sanitize(sq.SortKey)
	sq.SortDirection = sanitizer.Sanitize(sq.SortDirection)
}
