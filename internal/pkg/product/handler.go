package product

import "net/http"

type Handler interface {
	GetProduct(w http.ResponseWriter, r *http.Request)
	GetListPreviewProducts(w http.ResponseWriter, r *http.Request)
	SearchListPreviewProducts(w http.ResponseWriter, r *http.Request)
	GetProductRecommendations(w http.ResponseWriter, r *http.Request)
}
