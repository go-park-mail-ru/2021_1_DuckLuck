package models

import (
	"time"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/tools/sanitizer"
)

var (
	ReviewDateAddedSort = "date"

	ReviewPaginatorASC  = "ASC"
	ReviewPaginatorDESC = "DESC"
)

// Paginator for showing page of reviews
type PaginatorReviews struct {
	PageNum int `json:"page_num"`
	Count   int `json:"count"`
	SortReviewsOptions
}

type SortReviewsOptions struct {
	SortKey       string `json:"sort_key" valid:"in(rating|date)"`
	SortDirection string `json:"sort_direction" valid:"in(ASC|DESC)"`
}

func (pr *PaginatorReviews) Sanitize() {
	sanitizer := sanitizer.NewSanitizer()
	pr.SortKey = sanitizer.Sanitize(pr.SortKey)
	pr.SortDirection = sanitizer.Sanitize(pr.SortDirection)
}

type ViewReview struct {
	UserName      string    `json:"user_name" valid:"minstringlength(1)"`
	UserAvatar    string    `json:"user_avatar" valid:"minstringlength(1)"`
	DateAdded     time.Time `json:"date_added" valid:"notnull"`
	Rating        int       `json:"rating" valid:"int"`
	Advantages    string    `json:"advantages"`
	Disadvantages string    `json:"disadvantages"`
	Comment       string    `json:"comment"`
	IsPublic      bool      `json:"-"`
	UserId        int       `json:"-"`
}

type Review struct {
	ProductId     int    `json:"product_id"`
	Rating        int    `json:"rating" valid:"int"`
	Advantages    string `json:"advantages"`
	Disadvantages string `json:"disadvantages"`
	Comment       string `json:"comment"`
	IsPublic      bool   `json:"is_public"`
}

func (r *Review) Sanitize() {
	sanitizer := sanitizer.NewSanitizer()
	r.Advantages = sanitizer.Sanitize(r.Advantages)
	r.Disadvantages = sanitizer.Sanitize(r.Disadvantages)
	r.Comment = sanitizer.Sanitize(r.Comment)
}

type ReviewStatistics struct {
	Stars []int `json:"stars"`
}

// Set of reviews with count uniq sets of this size
type RangeReviews struct {
	ListPreviews  []*ViewReview `json:"list_reviews" valid:"notnull"`
	MaxCountPages int           `json:"max_count_pages"`
}
