package order

import "net/http"

type Handler interface {
	GetOrderFromCart(w http.ResponseWriter, r *http.Request)
	AddCompletedOrder(w http.ResponseWriter, r *http.Request)
}
