package category

import "net/http"

//go:generate mockgen -destination=./mock/mock_handler.go -package=mock github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/category Handler

type Handler interface {
	GetCatalogCategories(w http.ResponseWriter, r *http.Request)
	GetSubCategories(w http.ResponseWriter, r *http.Request)
}
