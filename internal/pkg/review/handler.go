package review

import "net/http"

type Handler interface {
	GetReviewStatistics(w http.ResponseWriter, r *http.Request)
	AddNewReview(w http.ResponseWriter, r *http.Request)
	GetReviewsForProduct(w http.ResponseWriter, r *http.Request)
	CheckReviewRights(w http.ResponseWriter, r *http.Request)
}
