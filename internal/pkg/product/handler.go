package product

import "net/http"

//go:generate mockgen -destination=./mock/mock_handler.go -package=mock github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/product Handler

type Handler interface {
	GetProduct(w http.ResponseWriter, r *http.Request)
	GetListPreviewProducts(w http.ResponseWriter, r *http.Request)
}
