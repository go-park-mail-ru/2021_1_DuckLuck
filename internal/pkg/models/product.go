package models

type Product struct {
	Id          uint64  `json:"-"`
	Title       string  `json:"title"`
	Cost        float32 `json:"cost"`
	Rating      float32 `json:"rating"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Image       string  `json:"image"`
}

type RangeProducts struct {
	ArrayOfProducts []*Product `json:"products"`
	MaxCountPages   int
}

type PaginatorProducts struct {
	PageNum       int `json:"page_num"`
	Count         int `json:"count"`
	SortKey       string `json:"sort_key"`
	SortDirection bool `json:"sort_direction"`
}

var (
	ProductsCostSort   = "cost"
	ProductsRatingSort = "rating"

	PaginatorASC  = true
	PaginatorDESC = false
)
