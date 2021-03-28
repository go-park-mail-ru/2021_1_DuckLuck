package cart

import "net/http"

type Handler interface {
	AddProductInCart(w http.ResponseWriter, r *http.Request)
	DeleteProductInCart(w http.ResponseWriter, r *http.Request)
	ChangeProductInCart(w http.ResponseWriter, r *http.Request)
	GetProductsFromCart(w http.ResponseWriter, r *http.Request)
}
