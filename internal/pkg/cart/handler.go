package cart

import "net/http"

//go:generate mockgen -destination=./mock/mock_handler.go -package=mock github.com/go-park-mail-ru/2021_1_DuckLuck/internal/pkg/cart Handler

type Handler interface {
	AddProductInCart(w http.ResponseWriter, r *http.Request)
	DeleteProductInCart(w http.ResponseWriter, r *http.Request)
	ChangeProductInCart(w http.ResponseWriter, r *http.Request)
	GetProductsFromCart(w http.ResponseWriter, r *http.Request)
	DeleteProductsFromCart(w http.ResponseWriter, r *http.Request)
}
