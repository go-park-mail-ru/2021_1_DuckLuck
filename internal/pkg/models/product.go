package models

type Product struct {
	Id          uint64      `json:"id"`
	Title       string      `json:"title"`
	Cost        ProductCost `json:"cost"`
	Rating      int         `json:"rating"`
	Description string      `json:"description"`
	Category    string      `json:"category"`
	Images      []string    `json:"images"`
}

type ViewProduct struct {
	Id           uint64      `json:"id"`
	Title        string      `json:"title"`
	Cost         ProductCost `json:"cost"`
	Rating       int         `json:"rating"`
	PreviewImage string      `json:"preview_image"`
}

type ProductCost struct {
	Discount float32 `json:"discount"`
	BaseCost float32 `json:"base_cost"`
}

type RangeProducts struct {
	ListPreviewProducts []*ViewProduct `json:"list_preview_products"`
	MaxCountPages       int            `json:"max_count_pages"`
}

type PaginatorProducts struct {
	PageNum       int    `json:"page_num"`
	Count         int    `json:"count"`
	SortKey       string `json:"sort_key"`
	SortDirection string `json:"sort_direction"`
}

var (
	ProductsCostSort   = "cost"
	ProductsRatingSort = "rating"

	PaginatorASC  = "ASC"
	PaginatorDESC = "DESC"
)
