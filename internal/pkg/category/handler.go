package category

import "net/http"

type Handler interface {
	GetCatalogCategories(w http.ResponseWriter, r *http.Request)
	GetSubCategories(w http.ResponseWriter, r *http.Request)
}
