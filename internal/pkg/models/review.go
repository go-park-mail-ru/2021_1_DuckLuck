package models

import (
	"github.com/go-park-mail-ru/2021_1_DuckLuck/internal/server/tools/sanitizer"
	"time"
)

// Paginator for showing page of reviews
type PaginatorReviews struct {
	PageNum       int    `json:"page_num"`
	Count         int    `json:"count"`
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
	UserName        string      `json:"user_name" valid:"minstringlength(3)"`
	UserAvatar		string 		`json:"user_avatar" valid:"minstringlength(3)"`
	DateAdded		time.Time	`json:"date_added" valid:"notnull"`
	Images       	[]string    `json:"images"`
	Rating			int			`json:"rating" valid:"int"`

}


