package favorites

import "net/http"

type Handler interface {
	AddProductToFavorites(w http.ResponseWriter, r *http.Request)
	DeleteProductFromFavorites(w http.ResponseWriter, r *http.Request)
	GetListPreviewFavorites(w http.ResponseWriter, r *http.Request)
}
