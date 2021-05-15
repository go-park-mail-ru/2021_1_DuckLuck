package models

var (
	FavoritesCostSort      = "cost"
	FavoritesRatingSort    = "rating"
	FavoritesDateAddedSort = "date"
	FavoritesDiscountSort  = "discount"

	FavoritesPaginatorASC  = "ASC"
	FavoritesPaginatorDESC = "DESC"
)

// Set of product with count uniq sets of this size
type RangeFavorites struct {
	ListPreviewProducts []*ViewFavorite `json:"list_preview_products" valid:"notnull"`
	MaxCountPages       int            `json:"max_count_pages"`
}

type FavoritesSortOptions struct {
	SortKey       string `json:"sort_key" valid:"in(cost|rating|date|discount)"`
	SortDirection string `json:"sort_direction" valid:"in(ASC|DESC)"`
}

// Paginator for showing page of product
type PaginatorFavorites struct {
	PageNum  int            `json:"page_num"`
	Count    int            `json:"count"`
	SortOptions
}

// View of product
// This model contains field for preview information
type ViewFavorite struct {
	Id           uint64       `json:"id"`
	Title        string       `json:"title" valid:"minstringlength(3)"`
	Price        ProductPrice `json:"price" valid:"notnull, json"`
	Rating       float32      `json:"rating" valid:"float, range(0, 10)"`
	CountReviews uint64       `json:"count_reviews"`
	PreviewImage string       `json:"preview_image" valid:"minstringlength(3)"`
}

